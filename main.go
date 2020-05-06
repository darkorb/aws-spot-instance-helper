package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chrisurwin/aws-spot-instance-helper/agent"
	"github.com/chrisurwin/aws-spot-instance-helper/slackhelpers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

//VERSION of the program
var VERSION = "v0.1.1-aznamier-dev"

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "aws-spot-instance-helper"
	app.Version = VERSION
	app.Usage = "Evacuates an AWS Spot Instance host when marked for termination"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug,d",
			Usage:  "Debug logging",
			EnvVar: "DEBUG",
		},
		cli.DurationFlag{
			Name:   "poll-interval,i",
			Value:  5 * time.Second,
			Usage:  "Polling interval for checks",
			EnvVar: "POLL_INTERVAL",
		},
		cli.StringFlag{
			Name:   "cattleURL,u",
			Usage:  "Cattle URL",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "cattleAccessKey,ck",
			Usage:  "Cattle Access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "cattleSecretKey,cs",
			Usage:  "Cattle Secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "slackWebhookUrl,s",
			Usage:  "Slack Webhook URL",
			EnvVar: "SLACK_WEBHOOK",
		},
		cli.StringFlag{
			Name:   "slackMessageSuffix,ss",
			Usage:  "Appears at the end of slack message - eg: @bob.jon or Tracking: 77283",
			EnvVar: "SLACK_MESSAGE_SUFFIX",
		},
		cli.BoolFlag{
			Name:   "slackInitAnnouncement,si",
			Usage:  "Initial announcement will be send to slack on a startup",
			EnvVar: "SLACK_INIT_ANNOUNCEMENT",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	if c.String("cattleURL") == "" {
		return fmt.Errorf("Cattle URL required")
	}
	if c.String("cattleAccessKey") == "" {
		return fmt.Errorf("Cattle Access Key required")
	}
	if c.String("cattleSecretKey") == "" {
		return fmt.Errorf("Cattle Secret Key required")
	}
	slc := &slackhelpers.SlackConfig{
		AnnonunceOnInit: c.Bool("slackInitAnnouncement"),
		WebhookURL:      c.String("slackWebhookUrl"),
		MessageSuffix:   c.String("slackMessageSuffix"),
	}

	a := agent.NewAgent(c.Duration("poll-interval"), c.String("cattleURL"), c.String("cattleAccessKey"), c.String("cattleSecretKey"), slc)

	return a.Start()
}
