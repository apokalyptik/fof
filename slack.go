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
		General  bool   `json:"is_general"`
		Archived bool   `json:"is_archived"`
		Channel  bool   `json:"is_channel"`
	} `json:"channels"`
}

func updateChannelList() error {
	resp, err := http.Get(fmt.Sprintf(
		"https://slack.com/api/channels.list?token=%s&exclude_archived=1",
		url.QueryEscape(slack.apiKey)))

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	var apiResp = &slackChannelResponse{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		return err
	}
	if apiResp.Ok != true {
		return fmt.Errorf("Error decoding channel list: %#v", apiResp)
	}
	var newChannelList = []string{}
	for _, c := range apiResp.Channels {
		newChannelList = append(newChannelList, c.Name)
	}
	if strings.Join(newChannelList, ",") != strings.Join(xhrOutput.data["channels"].([]string), ",") {
		xhrOutput.set("channels", newChannelList)
	}
	return nil
}

func mindChannelList() {
	for {
		if err := updateChannelList(); err != nil {
			log.Println("error getting channel list from slack:", err.Error())
			time.Sleep(time.Second * 30)
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

func (s *Slack) getUserListToIDs() (map[string]string, error) {
	var response = struct {
		OK      bool   `json:"ok"`
		Error   string `json:"error"`
		Members []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"members"`
	}{}
	resp, err := http.Get(
		fmt.Sprintf("https://slack.com/api/users.list?token=%s", s.apiKey))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&response); err != nil {
		return nil, err
	}
	var rval = map[string]string{}
	for i := range response.Members {
		rval[response.Members[i].Name] = response.Members[i].ID
	}
	return rval, nil
}

func (s *Slack) doOpenIM(username string) (string, error) {
	var response = struct {
		OK      bool   `json:"ok"`
		Error   string `json:"error"`
		Channel struct {
			ID string `json:"id"`
		} `json:"channel"`
	}{}
	reqURL := fmt.Sprintf("https://slack.com/api/im.open?token=%s&user=%s",
		url.QueryEscape(s.apiKey),
		url.QueryEscape(username),
	)
	resp, err := http.Get(reqURL)

	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&response); err != nil {
		return "", err
	}
	return response.Channel.ID, nil
}

type slackMsg struct {
	Where string
}

func (s *slackMsg) to(where string) *slackMsg {
	s.Where = where
	return s
}

func (s *slackMsg) send(text string) error {
	if s.Where[0] == '@' {
		if channel, err := udb.getChannelForIM(s.Where[1:]); err != nil {
			return err
		} else {
			s.Where = channel
		}
	} else {
		s.Where = "#slack-tools-testing"
	}

	slackMsgQueue <- url.Values{
		"channel": []string{s.Where},
		"text":    []string{text},
	}
	return nil
}

var slackMsgQueue = make(chan url.Values, 1024)

func mindSlackMsgQueue() {
	var reqURL = "https://slack.com/api/chat.postMessage"
	for {
		log.Println("neg1")
		payload := <-slackMsgQueue
		log.Println("zero")
		log.Println(payload)
		payload.Set("token", slack.apiKey)
		log.Println("one")
		resp, err := http.PostForm(reqURL, payload)
		log.Println("two")
		resp.Body.Close()
		log.Println("three")
		if err != nil {
			log.Println("error sending message via slack:", err.Error())
		}
		time.Sleep(time.Millisecond * 750)
		log.Println("four")
	}
}
