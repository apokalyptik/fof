package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var slack = &Slack{}

type slackChannelResponse struct {
	Ok       bool `json:"ok"`
	Channels []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		General  string `json:"is_general"`
		Archives string `json:"is_archived"`
		Channel  string `json:"is_channel"`
	} `json:"channels"`
}

func mindChannelList() {
	for {
		resp, err := http.Get(fmt.Sprintf(
			"https://slack.com/api/channels.list?token=%s&exclude_archived=1",
			url.QueryEscape(slack.apiKey)))

		if err != nil {
			resp.Body.Close()
			log.Println(err.Error())
			time.Sleep(time.Second * 5)
			continue
		}

		var apiResp = &slackChannelResponse{}
		json.NewDecoder(resp.Body).Decode(apiResp)
		resp.Body.Close()
		if apiResp.Ok != true {
			log.Printf("%#v", apiResp)
			time.Sleep(time.Second * 5)
			continue
		}
		var newChannelList = []string{}
		for _, c := range apiResp.Channels {
			newChannelList = append(newChannelList, c.Name)
		}
		if strings.Join(newChannelList, ",") != strings.Join(xhrOutput.data["channels"].([]string), ",") {
			xhrOutput.set("channels", newChannelList)
		}
	}
}

type slackGetCache struct {
	when time.Time
	lock sync.RWMutex
	url  string
	data []byte
}

type Slack struct {
	apiKey  string
	raidKey string
	url     string
	name    string
	emoji   string
}

func (s *Slack) msg() *slackMsg {
	return &slackMsg{}
}

type slackMsg struct {
	Where string `json:"channel"`
	From  string `json:"username,omitempty"`
	Icon  string `json:"icon_url,omitempty"`
	Emoji string `json:"icon_emoji,omitempty"`
	Text  string `json:"text"`
}

func (s *slackMsg) to(where string) *slackMsg {
	s.Where = where
	return s
}

func (s *slackMsg) from(who string) *slackMsg {
	s.From = who
	return s
}

func (s *slackMsg) icon(which string) *slackMsg {
	if len(which) < 7 {
		s.Emoji = which
		return s
	}
	var bit = strings.ToLower(which[:8])
	if bit[0:0] != "h" {
		s.Emoji = which
		return s
	}
	if bit[:4] != "http:" || bit[:5] != "https:" {
		s.Emoji = which
		return s
	}
	s.Icon = which
	return s
}

func (s *slackMsg) send(text string) error {
	if s.Icon == "" && s.Emoji == "" {
		s.Emoji = slack.emoji
	}
	if s.From == "" {
		s.From = slack.name
	}
	s.Text = text
	data, err := json.Marshal(s)
	if err != nil {
		log.Println("error marshing slack message:", err.Error())
		return err
	}
	slackMsgQueue <- url.Values{"payload": {string(data)}}
	return nil
}

var slackMsgQueue = make(chan url.Values, 1024)

func mindSlackMsgQueue() {
	for {
		payload := <-slackMsgQueue
		if _, err := http.PostForm(slack.url, payload); err != nil {
			log.Println("error sending message via slack:", err.Error())
		}
		time.Sleep(time.Millisecond * 750)
	}
}
