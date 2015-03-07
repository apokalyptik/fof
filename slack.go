package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

var slack = &Slack{}

type Slack struct {
	key   string
	url   string
	name  string
	emoji string
}

func (s *Slack) toChannel(which, message string) {
	s.sendMessage("#"+which, message)
}

func (s *Slack) toPerson(who, message string) {
	s.sendMessage("@"+who, message)
}

func (s *Slack) sendMessage(where, message string) {
	data, err := json.Marshal(map[string]string{
		"text":         message,
		"username":     s.name,
		"channel":      where,
		"icon_emoji":   s.emoji,
		"unfurl_links": "true",
	})
	if err != nil {
		log.Println(err.Error)
	}
	if _, err := http.PostForm(s.url, url.Values{"payload": {string(data)}}); err != nil {
		log.Println(err)
	}
}
