package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type slackChannelResponse struct {
	Ok       bool   `json:"ok"`
	Error    string `json:"error"`
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

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	var apiResp = &slackChannelResponse{}
	if err := json.NewDecoder(resp.Body).Decode(apiResp); err != nil {
		return err
	}
	if apiResp.Ok != true || apiResp.Error != "" {
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
