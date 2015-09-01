package main

import (
	"flag"
	"log"
	"runtime"
	"sync"
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
	var maxWorkers = 10
	var concurrency = make(chan struct{}, maxWorkers)
	go func() {
		for i := 0; i < maxWorkers; i++ {
			concurrency <- struct{}{}
		}
	}()
	for {
		startAll := time.Now()
		var nextRun = time.Now().Add(12 * time.Hour)
		log.Printf("Updating Users")
		users.update()
		log.Printf("Users Updated")
		users.lock.RLock()
		n := len(users.list)
		i := 0
		var w sync.WaitGroup
		log.Printf("Pulling user data")
		for username, u := range users.list {
			<-concurrency
			i++
			w.Add(1)
			go func(u *user, username string, i int) {
				start := time.Now()

				if err := u.pull(); err != nil {
					log.Printf("Error %s on %s %d/%d in %s", err.Error(), username, i, n, time.Now().Sub(start))
				}
				concurrency <- struct{}{}
				log.Printf("Freshened %s %d/%d in %s", username, i, n, time.Now().Sub(start))
				w.Done()
			}(u, username, i)
		}
		w.Wait()
		log.Println("finished in", time.Now().Sub(startAll))
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	mindUsers()
}
