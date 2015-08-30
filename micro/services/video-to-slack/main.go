package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/nlopes/slack"
)

func main() {
	log.Println("Starting up...")
	if buf, err := ioutil.ReadFile(slackConfigFile); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(buf, &slackCFG); err != nil {
			log.Fatal(err)
		}
	}
	slackClient = slack.New(slackCFG.Key)
	flag.Parse()
	mindTwitch()
}
