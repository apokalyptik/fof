package main

import (
	"flag"
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var bungieApiKey = "..."
var mgoServer = "127.0.0.1"
var userListAddress = "http://127.0.0.1:8879/users.json"

var mgoDB *mgo.Session
var client = &destinyClient{}

var users = &userDB{
	list: map[string]*user{},
}

func init() {
	flag.StringVar(&bungieApiKey, "apikey", bungieApiKey, "Bungie Platform API Key ( https://www.bungie.net/en/User/API )")
	flag.StringVar(&mgoServer, "mgo", mgoServer, "MongoDB addresses")
}

func mindUsers() {
	for {
		var nextRun = time.Now().Add(4 * time.Hour)

		users.update()

		users.lock.RLock()
		for username, user := range users.list {
			log.Printf("Freshening %s", username)
			if err := user.pull(); err != nil {
				log.Printf(err.Error())
			}
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
	if session, err := mgo.Dial(mgoServer); err != nil {
		log.Fatal(err)
	} else {
		mgoDB = session
	}
	mindUsers()
}
