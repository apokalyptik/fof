package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var slackToken string
var address string

type userDB struct {
	OK        bool        `json:"ok"`
	Error     string      `json:"error"`
	Members   interface{} `json:"members"`
	Timestamp time.Time   `json:"timestamp"`
	jsonOut   []byte
}

func (u *userDB) fetch() (*userDB, error) {
	rsp, err := http.Get(fmt.Sprintf("https://slack.com/api/users.list?token=%s", slackToken))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var newudb = &userDB{}
	var decoder = json.NewDecoder(rsp.Body)
	if err := decoder.Decode(newudb); err != nil {
		return nil, err
	}
	newudb.Timestamp = time.Now()
	if newudb.jsonOut, err = json.Marshal(newudb); err != nil {
		return newudb, err
	}
	if !newudb.OK {
		return newudb, errors.New(newudb.Error)
	}
	return newudb, nil
}

var udb = &userDB{
	OK:        false,
	Error:     "just initializing",
	Timestamp: time.Now(),
}

func init() {
	flag.StringVar(&slackToken, "st", "", "Slack API Token")
	flag.StringVar(&address, "listen", "127.0.0.1:8879", "listen on")
}

func main() {
	flag.Parse()
	if newudb, err := udb.fetch(); err != nil {
		log.Fatal(err.Error())
	} else {
		udb = newudb
	}
	go func() {
		t := time.Tick(15 * time.Minute)
		for {
			select {
			case <-t:
				if newudb, err := udb.fetch(); err != nil {
					log.Println(err.Error())
				} else {
					udb = newudb
				}
			}
		}
	}()
	http.HandleFunc("/users.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(udb.jsonOut)
	})
	http.ListenAndServe(address, nil)
}
