package slackhelpers

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/chrisurwin/aws-spot-instance-helper/rancherhelpers"
	log "github.com/sirupsen/logrus"
)

//SlackConfig - struct for keeping slack configuration
type SlackConfig struct {
	WebhookURL      string
	AnnonunceOnInit bool
	MessageSuffix   string
}

//SendMessage - method with error handling to send slack messages
func SendMessage(message, color string, slackConfig *SlackConfig, envMetaData *rancherhelpers.SelfMetaData) {

	currentEnvName := envMetaData.EnvName
	hostName := envMetaData.HostName
	messageSuffix := slackConfig.MessageSuffix
	hostLabels := ""

	for _, label := range envMetaData.HostLabels {
		//excluding standard rancher headers
		if !strings.Contains(label, "io.rancher.host") {
			hostLabels += "," + label
		}
	}
	hostLabels = strings.Replace(hostLabels, ",", "", 1)

	footer := "Spot Host Evacuator _(" + hostLabels + ")_"
	//building the message
	message = message + ": *" + currentEnvName + " / " + hostName + ":*  " + messageSuffix

	if slackConfig.WebhookURL != "" {
		log.Debug("Slack sending: " + message)

		attachment1 := slack.Attachment{
			Color:    &color,
			Text:     &message,
			Footer:   &footer,
			Fallback: &message,
		}

		payload := slack.Payload{
			Attachments: []slack.Attachment{attachment1},
			LinkNames:   "true",
		}

		err := slack.Send(slackConfig.WebhookURL, "", payload)
		if len(err) > 0 {
			fmt.Println(fmt.Sprintf("There was a problem sending Slack notification: %s\n", err))
		}

	} else {
		log.Debug("Slack not configured, not sending: " + message)
	}

}
