package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

var lfg = &lfgStore{
	Data:  map[string]map[string]*lfgEntry{},
	Users: map[string]*lfgUser{},
}

type lfgEntry struct {
	Expiry   time.Time `json:"expiry"`
	Username string    `json:"username"`
	Gamertag string    `json:"gamertag"`
}

type lfgUser struct {
	events []string
	expiry time.Time
}

type lfgStore struct {
	filename string
	lock     sync.Mutex
	Data     map[string]map[string]*lfgEntry
	Users    map[string]*lfgUser
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
	for event, users := range l.Data {
		for user, entry := range users {
			if runtime.After(entry.Expiry) {
				log.Println("Expired", user, event, entry.Expiry, "<", runtime)
				changed = true
				delete(users, user)
			}
		}
		if len(event) < 1 {
			changed = true
			delete(l.Data, event)
		}
	}
	if changed {
		l.save()
		l.emit()
	}
}

func (l *lfgStore) emit() {
	lfgOutput.set("lfg", lfg.Data)
	xhrOutput.set("lfg", lfg.Data)
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
	if previous, ok := l.Users[username]; ok {
		for _, e := range previous.events {
			delete(l.Data[e], username)
			if len(l.Data[e]) == 0 {
				delete(l.Data, e)
			}
		}
	}
	l.Users[username] = user
	for _, event := range events {
		if _, ok := l.Data[event]; !ok {
			l.Data[event] = map[string]*lfgEntry{}
		}
		l.Data[event][username] = &lfgEntry{
			Expiry:   expireTime,
			Username: username,
			Gamertag: gamertag,
		}
	}
	l.save()
	l.emit()
	return nil
}

func (l *lfgStore) save() error {
	if l.filename == "" {
		return errors.New("cannot persist withuot filename")
	}
	fp, err := os.Create(l.filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if err := enc.Encode(l); err != nil {
		return err
	}
	return nil
}

func (l *lfgStore) load(filename string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			l.filename = filename
			return nil
		}
		return err
	}
	l.filename = filename
	defer fp.Close()
	dec := json.NewDecoder(fp)
	if err := dec.Decode(&l); err != nil {
		return err
	}
	if l == nil {
		l = &lfgStore{
			Data:  map[string]map[string]*lfgEntry{},
			Users: map[string]*lfgUser{},
		}
	}
	l.emit()
	return nil
}
