package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var errUserNotFound = fmt.Errorf("FOF User not found")

var users = map[string]string{}       // username => gt
var userIDs = map[string]string{}     // gt => username
var gamertags = map[string]string{}   // username => gt and userid => gt
var lastSeen = map[string]time.Time{} // todo: populate
var destinyIDs = map[string]string{}  // todo: populate

var usersLock sync.RWMutex

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
		newUserList[v.UserName] = v.UserID
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

	usersLock.Lock()
	users = newUserList
	gamertags = newGamertags
	userIDs = newUserIDs
	usersLock.Unlock()

	for i := range added {
		pubUserAdded(added[i])
	}
	for i := range deleted {
		pubUserRemoved(deleted[i])
	}
}

func getUserDestinyID(username string) (string, error) {
	usersLock.RLock()
	defer usersLock.RUnlock()
	if ID, ok := users[username]; ok {
		return ID, nil
	}
	return "", errUserNotFound
}

func mindUserList() {
	for {
		updateUserList()
		time.Sleep(5 * time.Minute)
	}
}
