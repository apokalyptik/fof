package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

var slackMsgQueue = make(chan url.Values, 1024)

func mindSlackMsgQueue() {
	var reqURL = "https://slack.com/api/chat.postMessage"
	for {
		payload := <-slackMsgQueue
		payload.Set("token", slack.apiKey)
		payload.Set("as_user", "true")
		resp, err := http.PostForm(reqURL, payload)
		resp.Body.Close()
		if err != nil {
			log.Println("error sending message via slack:", err.Error())
		}
		time.Sleep(time.Millisecond * 750)
	}
}
