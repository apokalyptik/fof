package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var slack = &Slack{}

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
