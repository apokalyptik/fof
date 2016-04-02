package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var userList = privateUsers{}

type privateUserList struct {
	Members []struct {
		Deleted         bool   `json:"deleted"`
		Bot             bool   `json:"is_bot"`
		Restricted      bool   `json:"is_restricted"`
		UltraRestricted bool   `json:"is_ultra_restricted"`
		ID              string `json:"id"`
		Name            string `json:"name"`
		Profile         struct {
			GamerTag string `json:"first_name"`
		} `json:"profile"`
	} `json:"members"`
}

type privateUser struct {
	ID       string
	Name     string
	GamerTag string
}

type privateUsers map[string]privateUser

func mindPrivateUserList() {
	if u, err := getPrivateUserList(); err != nil {
		log.Fatal("Error fetching initial user list!", err.Error())
	} else {
		userList = u
	}
	t := time.Tick(10 * time.Minute)
	for {
		select {
		case <-t:
			if u, err := getPrivateUserList(); err != nil {
				log.Println("Error fetching user list:", err.Error())
			} else {
				userList = u
			}
		}
	}
}

func getPrivateUserList() (privateUsers, error) {
	var rval = privateUsers{}
	var raw privateUserList
	rsp, err := http.Get("http://127.0.0.1:8879/users.json")
	if err != nil {
		return rval, err
	}
	defer rsp.Body.Close()
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&raw); err != nil {
		return rval, err
	}
	for _, user := range raw.Members {
		if user.Bot {
			continue
		}
		if user.Deleted {
			continue
		}
		if user.Restricted {
			continue
		}
		if user.UltraRestricted {
			continue
		}
		rval[user.ID] = privateUser{
			ID:       user.ID,
			Name:     user.Name,
			GamerTag: user.Profile.GamerTag,
		}
	}
	return rval, nil
}
