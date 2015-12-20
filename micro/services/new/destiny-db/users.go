package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var errUserNotFound = fmt.Errorf("FOF User not found")

var users = map[string]string{}       // username => gt. Only people seen < 30 days ago
var userIDs = map[string]string{}     // gt => username. All members. Regardless of last seen time
var gamertags = map[string]string{}   // username => gt and userid => gt. All members regardless of last seen time
var lastSeen = map[string]time.Time{} // userID => last seen time, used for deciding whether a user gets into the users list
var destinyIDs = map[string]string{}  // todo: populate, persist

type slackUserListResponse struct {
	OK      bool   `json:"ok"`
	Error   string `json:"error"`
	Members []struct {
		UserID          string `json:"id"`
		UserName        string `json:"name"`
		Bot             bool   `json:"is_bot"`
		Deleted         bool   `json:"deleted"`
		Restricted      bool   `json:"is_restricted"`
		UltraRestricted bool   `json:"is_ultra_restricted"`
		Profile         struct {
			GamerTag string `first_name`
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
	lastSeen = newLastSeen
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
	var newUserList = map[string]string{}
	var newUserIDs = map[string]string{}
	var newGamertags = map[string]string{}
	var seenCutoff = time.Now().Add(0 - (24 * 30 * time.Hour))
	for _, v := range slackData.Members {
		if v.Bot {
			continue
		}
		if v.Deleted {
			continue
		}
		if v.Restricted {
			continue
		}
		if v.UltraRestricted {
			continue
		}
		if seen, ok := lastSeen[v.UserID]; ok {
			if seen.After(seenCutoff) {
				newUserList[v.UserName] = v.UserID
			}
		}
		newUserIDs[v.UserID] = v.UserName
		newGamertags[v.UserName] = v.Profile.GamerTag
		newGamertags[v.UserID] = v.Profile.GamerTag
	}
	var added = []string{}
	var deleted = []string{}
	for k := range newUserList {
		if _, ok := users[k]; !ok {
			added = append(added, k)
		}
	}
	for k := range users {
		if _, ok := newUserList[k]; !ok {
			deleted = append(deleted, k)
		}
	}

	log.Println("going from", len(users), "to", len(newUserList), "users")
	users = newUserList
	gamertags = newGamertags
	userIDs = newUserIDs

	for i := range added {
		pubUserAdded(added[i])
	}
	for i := range deleted {
		pubUserRemoved(deleted[i])
	}
}

func getUserDestinyID(username string) (string, error) {
	if ID, ok := users[username]; ok {
		return ID, nil
	}
	return "", errUserNotFound
}

func mindUserList() {
	for {
		updateSeenList()
		updateUserList()
		time.Sleep(5 * time.Minute)
	}
}
