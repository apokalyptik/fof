package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/apokalyptik/fof/lib/destiny"
)

var creds = map[string]string{}

var client *destiny.Client

func init() {
	if fp, err := os.Open("creds.json"); err != nil {
		log.Fatalf("Unable to open creds.json: %s", err.Error())
	} else {
		dec := json.NewDecoder(fp)
		if err := dec.Decode(&creds); err != nil {
			log.Fatalf("Unable to decode creds.json: %s", err.Error())
		}
	}
	client = destiny.New(creds["bungieAPI"], "github.com/apokalyptik/fof/micro/services/new/destiny-db")
}

func main() {
	go func() {
		c := subUserAdded()
		for {
			u := <-c
			log.Println("added user", u)
		}
	}()
	go func() {
		c := subUserRemoved()
		for {
			u := <-c
			log.Println("deleted user", u)
		}
	}()
	mindUserList()
}
