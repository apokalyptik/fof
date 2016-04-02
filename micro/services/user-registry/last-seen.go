package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var seenList = map[string]time.Time{}

func mindSeenList() {
	if s, err := getSeenList(); err != nil {
		log.Println("Error fetching initial seen list!", err.Error())
	} else {
		seenList = s
	}
	t := time.Tick(10 * time.Minute)
	for {
		select {
		case <-t:
			if s, err := getSeenList(); err != nil {
				log.Println("Error fetching seen list:", err.Error())
			} else {
				seenList = s
			}
		}
	}
}

func getSeenList() (map[string]time.Time, error) {
	var rval = map[string]time.Time{}
	rsp, err := http.Get("http://127.0.0.1:8890/seen.json")
	if err != nil {
		return rval, err
	}
	defer rsp.Body.Close()
	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&rval)
	return rval, err
}
