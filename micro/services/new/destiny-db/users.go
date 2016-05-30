package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type bungieSearchResponse []struct {
	MembershipId string `json:"membershipId"`
}

type slackUserListResponse struct {
	OK      bool   `json:"ok"`
	Error   string `json:"error"`
	Members []struct {
		ID              int    `json:"ID"`
		UserID          string `json:"id"`
		UserName        string `json:"name"`
		Bot             bool   `json:"is_bot"`
		Deleted         bool   `json:"deleted"`
		Restricted      bool   `json:"is_restricted"`
		UltraRestricted bool   `json:"is_ultra_restricted"`
		Profile         struct {
			GamerTag string `json:"first_name"`
		} `json:"profile"`
	} `json:"members"`
}

func updateSeenList() {
	rsp, err := http.Get(creds["seenURL"])
	if err != nil {
		log.Println("error fetching seen list:", err.Error())
		return
	}
	defer rsp.Body.Close()
	var newLastSeen = map[string]time.Time{}
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&newLastSeen); err != nil {
		log.Println("error decoding seen data:", err.Error())
		return
	}
	if err := state.Put("last-seen", newLastSeen); err != nil {
		log.Println("error storing last-seen:", err.Error())
	}
}

func updateUserList() {
	var url = fmt.Sprintf("https://slack.com/api/users.list?token=%s", creds["slackAdminToken"])
	rsp, err := http.Get(url)
	if err != nil {
		log.Println("error fetching slack user list:", err.Error())
		return
	}
	defer rsp.Body.Close()
	var slackData slackUserListResponse
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&slackData); err != nil {
		log.Println("error decoding slack user list response:", err.Error())
		return
	}
	if !slackData.OK {
		log.Println("error in slack user list response:", slackData.Error)
		return
	}
	if err := state.Put("slack-members", slackData.Members); err != nil {
		log.Println("error storing slack-members:", err.Error())
	}
}
