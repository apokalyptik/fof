package main

import (
	"log"
	"sync"
	"time"
)

var lfg *lfgStore

type lfgEntry struct {
	Expiry   time.Time `json:"expiry"`
	Username string    `json:"username"`
	Gamertag string    `json:"gamertag"`
}

func init() {
	lfg = &lfgStore{
		data:  map[string]map[string]*lfgEntry{},
		users: map[string]*lfgUser{},
	}
	lfg.emit()
	go lfg.mindExpiration()
}

type lfgUser struct {
	events []string
	expiry time.Time
}

type lfgStore struct {
	filename string
	lock     sync.Mutex
	data     map[string]map[string]*lfgEntry
	users    map[string]*lfgUser
}

func (l *lfgStore) mindExpiration() {
	ticker := time.Tick(time.Minute)
	for {
		select {
		case <-ticker:
			l.prune()
		}
	}
}

func (l *lfgStore) prune() {
	var changed = false
	var runtime = time.Now()
	l.lock.Lock()
	defer l.lock.Unlock()
	for event, users := range l.data {
		for user, entry := range users {
			if runtime.After(entry.Expiry) {
				log.Println("Expired", user, event, entry.Expiry, "<", runtime)
				changed = true
				delete(users, user)
			}
		}
		if len(event) < 1 {
			changed = true
			delete(l.data, event)
		}
	}
	if changed {
		l.emit()
	}
}

func (l *lfgStore) emit() {
	lfgOutput.set("lfg", lfg.data)
}

func (l *lfgStore) add(username string, expiry time.Duration, events ...string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	var expireTime = time.Now().Add(expiry)
	var user = &lfgUser{
		expiry: expireTime,
		events: events,
	}
	var gamertag string
	if gt, err := udb.getGamertagForUser(username); err != nil {
		return err
	} else {
		gamertag = gt
	}
	if previous, ok := l.users[username]; ok {
		for _, e := range previous.events {
			delete(l.data[e], username)
			if len(l.data[e]) == 0 {
				delete(l.data, e)
			}
		}
	}
	l.users[username] = user
	for _, event := range events {
		if _, ok := l.data[event]; !ok {
			l.data[event] = map[string]*lfgEntry{}
		}
		l.data[event][username] = &lfgEntry{
			Expiry:   expireTime,
			Username: username,
			Gamertag: gamertag,
		}
	}
	l.emit()
	return nil
}
