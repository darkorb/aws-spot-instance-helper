package rancherhelpers

import (
	"io/ioutil"
	"net/http"
	"strings"

	client "github.com/rancher/go-rancher/v2"
	log "github.com/sirupsen/logrus"
)

const (
	rancherMetaData = "http://169.254.169.250"
)

//SelfMetaData - struct to keep rancher metadata to be passed around
type SelfMetaData struct {
	EnvName    string
	HostName   string
	HostLabels []string
}

//DetectSelfMetaData - detects basic info from Rancher metadata about itself
func DetectSelfMetaData() *SelfMetaData {

	envName, _ := GetRancherMetadata("/latest/self/stack/environment_name")
	log.Debug("Environment: " + envName)
	hostName, _ := GetRancherMetadata("/latest/self/host/name")
	log.Debug("Host: " + hostName)
	allHostLabels, _ := GetRancherMetadata("/latest/self/host/labels")

	cattleLabels := []string{}
	for _, label := range strings.Fields(allHostLabels) {
		value, _ := GetRancherMetadata("/latest/self/host/labels/" + label)
		log.Debug("Host Label: " + label + "=" + value)
		cattleLabels = append(cattleLabels, label+"="+value)
	}

	return &SelfMetaData{
		EnvName:    envName,
		HostName:   hostName,
		HostLabels: cattleLabels,
	}

}

//GetRancherMetadata - Function to query local rancher metadata
func GetRancherMetadata(path string) (string, error) {

	resp, err := http.Get(rancherMetaData + path)
	if err != nil {
		log.Warn("invalid path passed " + path)
		if resp != nil {
			defer resp.Body.Close()
		}
		return "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn("GetRancherMetadata: Received invalid response " + path)
		return "", err
	}
	return string(body), err
}

//EvacuateHost - Function to evacuate a Rancher host
func EvacuateHost(hostName string, c *client.RancherClient) (bool, error) {

	//Get a list of Hosts
	hosts, err := c.Host.List(nil)
	if err != nil {
		log.Error("EvacuateHost: Error getting host list")
		return false, err
	}
	for _, h := range hosts.Data {
		if h.Hostname == hostName {
			_, err := c.Host.ActionEvacuate(&h)
			if err != nil {
				return false, err
			}
		}
	}
	return true, err
}
