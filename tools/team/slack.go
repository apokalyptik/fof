package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var slack = &Slack{}

type Slack struct {
	apiKey  string
	raidKey string
	name    string
	emoji   string
}

func (s *Slack) msg() *slackMsg {
	return &slackMsg{}
}

func (s *Slack) getUserDataMap() (IDs map[string]string, GTs map[string]string, err error) {
	IDs = map[string]string{}
	GTs = map[string]string{}
	var maxAttempts = 5
	for i := 1; i <= maxAttempts; i++ {
		var response = struct {
			OK      bool   `json:"ok"`
			Error   string `json:"error"`
			Members []struct {
				Name    string `json:"name"`
				ID      string `json:"id"`
				Profile struct {
					FirstName string `json:"first_name"`
				} `json:"profile"`
			} `json:"members"`
		}{}
		resp, err := http.Get(
			fmt.Sprintf("https://slack.com/api/users.list?token=%s", s.apiKey))
		defer resp.Body.Close()
		if err != nil {
			if i != maxAttempts {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			return nil, nil, err
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&response); err != nil {
			if i != maxAttempts {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			return nil, nil, err
		}
		for i := range response.Members {
			IDs[response.Members[i].Name] = response.Members[i].ID
			GTs[response.Members[i].Name] = response.Members[i].Profile.FirstName
		}
		return IDs, GTs, nil
	}
	return nil, nil, errSlackUserIDNotFound
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
