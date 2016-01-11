package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/apokalyptik/fof/lib/destiny"
	_ "github.com/mattn/go-sqlite3"
)

var creds = map[string]string{}
var destinyClient *destiny.Platform

func init() {
	if fp, err := os.Open("creds.json"); err != nil {
		log.Fatalf("Unable to open creds.json: %s", err.Error())
	} else {
		dec := json.NewDecoder(fp)
		if err := dec.Decode(&creds); err != nil {
			log.Fatalf("Unable to decode creds.json: %s", err.Error())
		}
	}
	destinyClient = destiny.New(
		creds["bungieAPI"],
		"github.com/apokalyptik/fof/micro/services/new/destiny-db",
	).Platform(destiny.PlatformXBL)
}

func main() {
	initSQL()
	log.Println("Minding Summaries")
	go mindSummaryUpdates()
	log.Println("Minding Userlists")
	go mindUserList()
	log.Println("Minding CharacterSummaries")
	go mindCharacterSummaryUpdates()
	log.Println("Waiting")
	var wait = make(chan struct{})
	<-wait
}
