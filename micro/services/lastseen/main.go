package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nlopes/slack"
)

type seenMap map[string]time.Time

func (s seenMap) saw(userID string) {
	now := time.Now()
	if last, ok := s[userID]; ok {
		if now.Sub(last) < time.Minute {
			return
		}
	}
	if _, err := updateQuery.Exec(userID, now.String()); err != nil {
		log.Fatal(err)
	}
	s[userID] = now
	log.Println("saw", userID)
}

var myDB *sql.DB
var seen = seenMap{}
var updateQuery *sql.Stmt

var dbPath = "./db.sqlite3"
var listenOn = ":8890"
var slackAPIKey = "..."

func init() {
	flag.StringVar(&listenOn, "listen", listenOn, "ip:port to listen on for requests")
	flag.StringVar(&dbPath, "db", dbPath, "path to the database")
	flag.StringVar(&slackAPIKey, "slack", slackAPIKey, "slack bot users api key")
}

func main() {
	flag.Parse()

	if db, err := sql.Open("sqlite3", dbPath); err != nil {
		log.Fatal(err)
	} else {
		myDB = db
	}

	if _, err := myDB.Exec("CREATE TABLE IF NOT EXISTS `seen` (`id` STRING PRIMARY KEY,`when` TEXT);"); err != nil {
		log.Fatal(err)
	}

	if stmt, err := myDB.Prepare("INSERT OR REPLACE INTO `seen` (`id`,`when`) VALUES(?,?)"); err != nil {
		log.Fatal(err)
	} else {
		updateQuery = stmt
	}

	if rows, err := myDB.Query("SELECT * FROM `seen`"); err != nil {
		log.Fatal(err)
	} else {
		for rows.Next() {
			var id string
			var when string
			rows.Scan(&id, &when)
			seen[id], err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", when)
			if err != nil {
				log.Println("failed parsing time", when, "for", id)
			}
		}
		rows.Close()
	}

	for u, t := range seen {
		log.Println("Loaded", u, "last seen time", t)
	}

	http.HandleFunc("/seen.json", func(w http.ResponseWriter, r *http.Request) {
		e := json.NewEncoder(w)
		e.Encode(seen)
	})
	go http.ListenAndServe(listenOn, nil)

	api := slack.New(slackAPIKey)
	api.SetDebug(false)
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	for {
		msg := <-rtm.IncomingEvents
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			seen.saw(ev.User)
		case *slack.PresenceChangeEvent:
			seen.saw(ev.User)
		case *slack.UserTypingEvent:
			seen.saw(ev.User)
		case *slack.UserChangeEvent:
			seen.saw(ev.User.ID)
		case *slack.TeamJoinEvent:
			seen.saw(ev.User.ID)
		case *slack.RTMError:
			log.Printf(ev.Error())
		case *slack.InvalidAuthEvent:
			log.Fatal("Invalid credentials")
		}
	}
}
