package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var usersLock sync.RWMutex
var users = map[string]string{}
var errUserNotFound = fmt.Errorf("FOF User not found")

func updateUserList() {
	rsp, err := http.Get("http://fofgaming.com:8880/fof/members.json")
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	var data = []struct {
		User    string `json:"username"`
		Destiny string `json:"destiny"`
	}{}
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&data); err != nil {
		return
	}
	var newUserList = map[string]string{}
	for _, v := range data {
		newUserList[v.User] = v.Destiny
	}
	usersLock.Lock()
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
	users = newUserList
	usersLock.Unlock()
	for i := range added {
		pubUserAdded(added[i])
	}
	for i := range deleted {
		pubUserRemoved(deleted[i])
	}
}

func getUserListCopy() map[string]string {
	usersLock.RLock()
	defer usersLock.RUnlock()
	var usersCopy = map[string]string{}
	for k, v := range users {
		usersCopy[k] = v
	}
	return usersCopy
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
