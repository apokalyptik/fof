package main

import (
	"flag"
	"log"
	"sync"
	"time"
)

var bungieApiKey string

var client = &destinyClient{}

var users = struct {
	list map[string]*user
	lock sync.RWMutex
}{
	list: map[string]*user{
		"demitriousk": &user{
			name: "demitriousk",
			data: map[string]*userBit{},
		},
	},
}

func init() {
	flag.StringVar(&bungieApiKey, "apikey", "...", "Bungie Platform API Key ( https://www.bungie.net/en/User/API )")
}

func mindUsers() {
	for {
		var nextRun = time.Now().Add(4 * time.Hour)

		users.lock.Lock()
		// TODO: Freshen user map, add new users... from slack, remove removed users...
		users.lock.Unlock()

		users.lock.RLock()
		for username, user := range users.list {
			log.Printf("Freshening %s", username)
			log.Printf("%#v", user)
		}
		users.lock.RUnlock()
		if !time.Now().After(nextRun) {
			log.Printf("Sleeping for %s", nextRun.Sub(time.Now()).String())
			time.Sleep(nextRun.Sub(time.Now()))
		}
	}
}

func main() {
	flag.Parse()
	if bungieApiKey == "..." {
		log.Fatal("API Key required")
	}
	client.apiKey = bungieApiKey
	mindUsers()
}
