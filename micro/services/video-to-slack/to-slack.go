package main

import (
	"flag"

	"github.com/nlopes/slack"
)

var slackConfigFile = "./slack.json"
var slackCFG = &slackConfig{}
var slackClient *slack.Client

type slackConfig struct {
	Key     string `json:"api_key"`
	Channel string `json:"channel"`
}

func init() {
	flag.StringVar(&slackConfigFile, "slack", slackConfigFile, "path to the slack config file")
}
