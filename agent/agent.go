package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chrisurwin/aws-spot-instance-helper/awshelpers"
	"github.com/chrisurwin/aws-spot-instance-helper/healthcheck"
	"github.com/chrisurwin/aws-spot-instance-helper/rancherhelpers"
	"github.com/chrisurwin/aws-spot-instance-helper/slackhelpers"

	client "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

//Agent - Struct for Agent
type Agent struct {
	probePeriod   time.Duration
	httpClient    http.Client
	rancherClient *client.RancherClient
	selfMetadata  *rancherhelpers.SelfMetaData
	slackConfig   *slackhelpers.SlackConfig
}

//NewAgent - Function to expose NewAgent
func NewAgent(probePeriod time.Duration, cattleURL, cattleAccessKey, cattleSecretKey string, slackConfig *slackhelpers.SlackConfig) *Agent {

	var opts = &client.ClientOpts{
		Url:       cattleURL,
		AccessKey: cattleAccessKey,
		SecretKey: cattleSecretKey,
	}

	var rc, _ = client.NewRancherClient(opts)

	var sm = rancherhelpers.DetectSelfMetaData()

	return &Agent{
		probePeriod: probePeriod,
		httpClient: http.Client{
			Timeout: time.Duration(2 * time.Second),
		},
		rancherClient: rc,
		selfMetadata:  sm,
		slackConfig:   slackConfig,
	}
}

//Start - Function to start agent
func (a *Agent) Start() error {
	log.Debug("Agent ", a.selfMetadata.HostName)

	if a.slackConfig.AnnonunceOnInit {
		slackhelpers.SendMessage("Monitoring", "good", a.slackConfig, a.selfMetadata)
	}
	go healthcheck.StartHealthcheck()
	ticker := time.NewTicker(a.probePeriod)
	for tk := range ticker.C {
		log.Debug("Check at ", tk)
		aws, err := awshelpers.GetAWSInfoBool("/latest/meta-data/", 200)
		if aws && err == nil {
			t, err := awshelpers.GetAWSInfoBool("/latest/meta-data/spot/termination-time", 200)
			if t && err != nil {
				log.Info("Instance is marked for termination")
				//Get Host
				hostname, err := rancherhelpers.GetRancherMetadata("/latest/self/host/hostname")
				if hostname != "" && err == nil {
					log.Info("Instance is marked for termination")
					//Get Host
					hostname, err := rancherhelpers.GetRancherMetadata("/latest/self/host/hostname")
					if hostname != "" && err == nil {
						log.Info("Instance is marked for termination")

						// Notify slack channel if configured
						slackhelpers.SendMessage("EVACUATING", "danger", a.slackConfig, a.selfMetadata)

						//Evacuate Host
						_, err := rancherhelpers.EvacuateHost(hostname, a.rancherClient)
						if err != nil {
							log.Error("There was a problem evacuating host...but as its marked for termination everything will get rescheduled anyway!!!")
						} else {
							log.Info("Host has been evacuated")
						}

					} else {
						log.Info("Host has been evacuated")
					}
				} else {
					return err
				}
			}
		} else {
			log.Info("Possibly not an AWS host")
		}
	}
	return fmt.Errorf("Agent returned an error")
}
