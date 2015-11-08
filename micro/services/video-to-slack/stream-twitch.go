package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/levi/twch"
	"github.com/nlopes/slack"
)

var twitchConfigFile = "./twitch.json"
var twitchCFG = &twitchConfig{
	Channels: []string{},
}

type twitchConfig struct {
	Key      string   `json:"oauth_key"`
	Channels []string `json:"channels"`
}

func init() {
	flag.StringVar(&twitchConfigFile, "twitch", twitchConfigFile, "path to the twitch config file")
}

func mindTwitch() {
	log.Println("[tw] Starting up")
	if buf, err := ioutil.ReadFile(twitchConfigFile); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(buf, &twitchCFG); err != nil {
			log.Fatal(err)
		}
	}
	client, err := twch.NewClient(twitchCFG.Key, nil)
	if err != nil {
		log.Fatal(err)
	}

	var active = map[string]int{}
	var lastActive = map[string]time.Time{}

	var twitchTicker = time.Tick(time.Second)
	for {
		streamers := map[string]bool{}
		nextRun := time.Now().Add(1 * time.Minute)
		for _, c := range twitchCFG.Channels {
			<-twitchTicker
			found := false
			if s, _, err := client.Streams.GetStream(c); err == nil {
				if s.ID != nil {
					details, _ := json.MarshalIndent(s, "\t", "\t")
					detailsString := string(details)
					streamers[c] = true
					sid := *s.ID
					if oid, ok := active[c]; ok {
						if oid != sid {
							found = true
							log.Printf(
								"[tw] Found changed stream for %s id: %d -> %d\n%s",
								c, oid, *s.ID, detailsString,
							)
						} else {
							lastActive[c] = time.Now()
							log.Printf(
								"[tw] Continued twitch stream for %s id: %d\n",
								c, *s.ID,
							)
						}
					} else {
						found = true
						log.Printf(
							"[tw] Found twitch stream for %s id: %d\n%s",
							c, *s.ID, detailsString,
						)
					}
					if found {
						active[c] = sid
						lastSeen := time.Now().Add(0 - (24 * time.Hour))
						if ls, ok := lastActive[c]; ok {
							lastSeen = ls
						}
						lastActive[c] = time.Now()
						if time.Now().Sub(lastSeen) < (15 * time.Minute) {
							log.Printf(
								"[tw] Skipping twitch %s because of recent activity %s ago",
								c,
								time.Now().Sub(lastSeen).String(),
							)
							continue
						}
						displayName := c
						if s.Channel.DisplayName != nil {
							displayName = *s.Channel.DisplayName
						}
						game := "something"
						if s.Game != nil {
							game = *s.Game
						}
						messageParams := slack.NewPostMessageParameters()
						messageParams.AsUser = true
						messageParams.Parse = "full"
						messageParams.LinkNames = 1
						messageParams.UnfurlMedia = true
						messageParams.UnfurlLinks = true
						messageParams.EscapeText = false
						messageParams.Attachments = append(messageParams.Attachments, slack.Attachment{
							Title:     fmt.Sprintf("Watch %s play %s", displayName, game),
							TitleLink: *s.Channel.URL,
							ThumbURL:  *s.Preview.Small,
						})
						_, _, err := slackClient.PostMessage(
							slackCFG.Channel,
							fmt.Sprintf("*%s* has begun streaming *%s* at %s", *s.Channel.Name, game, *s.Channel.URL),
							messageParams,
						)
						if err != nil {
							log.Printf("[tw] Error sending slack message: %s", err.Error())
						}
					}
				} else {
					log.Printf("[tw] Found no twitch stream for %s", c)
				}
			} else {
				log.Printf("[tw] Error checking twitch stream: %s", err.Error())
			}
		}
		for k := range active {
			if _, ok := streamers[k]; !ok {
				log.Printf("[tw] Removing twitch stream: %s", k)
				delete(active, k)
			}
		}
		for k := range lastActive {
			if time.Now().Sub(lastActive[k]) > (30 * time.Minute) {
				log.Printf("[tw] Removing last activity marker for twitch stream: %s", k)
				delete(lastActive, k)
			}
		}
		sleepFor := nextRun.Sub(time.Now())
		if sleepFor < (15 * time.Second) {
			sleepFor = 15 * time.Second
		}
		time.Sleep(sleepFor)
	}
}
