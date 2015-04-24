package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

var slack = &Slack{}

type Slack struct {
	xlineKey string
	raidKey  string
	needKey  string
	url      string
	name     string
	emoji    string
}

func (s *Slack) toChannel(which, message string, as ...string) {
	s.sendMessage("#"+which, message, as...)
}

func (s *Slack) toPerson(who, message string, as ...string) {
	s.sendMessage("@"+who, message, as...)
}

func (s *Slack) sendMessage(where, message string, as ...string) {
	var name = s.name
	var emoji = s.emoji
	if len(as) >= 1 {
		switch as[0] {
		case "stickybot":
			name = "stickybot"
			emoji = ":star2:"
		}
	}
	data, err := json.Marshal(map[string]string{
		"text":         message,
		"username":     name,
		"channel":      where,
		"icon_emoji":   emoji,
		"unfurl_links": "true",
	})
	if err != nil {
		log.Println(err.Error)
	}
	if _, err := http.PostForm(s.url, url.Values{"payload": {string(data)}}); err != nil {
		log.Println(err)
	}
}
