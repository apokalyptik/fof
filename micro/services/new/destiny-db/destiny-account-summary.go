package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/apokalyptik/fof/lib/destiny"
)

func mindSummaryUpdates() {
	t := time.Tick(time.Minute)
	for {
		exec("INSERT OR IGNORE INTO destinyAccountSummary (UserID) SELECT ID FROM users")
		users := getUsersNeedingSummaryUpdates()
		for _, uid := range users {
			user := getUserByID(uid)
			if skipUntil, err := user.getMetaTime("skipDestinyAccountSummary"); err == nil {
				if time.Now().Before(skipUntil) {
					continue
				}
			} else {
				log.Printf("Error seeking skipDestinyAccountSummary for %d:%s: %s", user.ID, user.GamerTag, err.Error())
			}
			if rsp, err := destinyClient.AccountSummary(user.DestinyID); err != nil {
				log.Printf("error getting account summary: %s", err.Error())
			} else {
				var i json.RawMessage
				var d struct {
					Chars []struct {
						Base struct {
							ID   string    `json:"characterId"`
							Seen time.Time `json:"dateLastPlayed"`
						} `json:"characterBase"`
					} `json:"characters"`
				}
				if err := rsp.Into(&i); err != nil {
					if err == destiny.ErrDestinyAccountNotFound {
						log.Printf("No destiny account for user: %s: Skipping for 24 hours", user.GamerTag)
						if err := user.setMeta("skipDestinyAccountSummary", time.Now().Add(24*time.Hour).String()); err != nil {
							log.Printf("Error setting skipDestinyAccountSummary for user: %s", err.Error())
						}
						continue
					}
					log.Printf("Error decoding account summary response: %s for user %#v", err.Error(), uid)
					continue
				}
				_, e := exec(
					"INSERT OR REPLACE INTO destinyAccountSummary (UserID,Fetched,Raw) VALUES(?,datetime(?, 'utc'),?)",
					uid,
					time.Now(),
					[]byte(i),
				)
				if e != nil {
					log.Fatalf("Error updating destinyAccountSummary: %s", e.Error())
				}
				if err := json.Unmarshal(i, &d); err == nil {
					for _, ch := range d.Chars {
						_, e := exec(
							"INSERT OR REPLACE INTO destinyCharacters (UserID,CharacterID,Played) VALUES(?,?,datetime(?, 'utc'))",
							uid,
							ch.Base.ID,
							ch.Base.Seen,
						)
						if e != nil {
							log.Fatal(e)
						}
					}
				}
			}
		}
		<-t
	}
}

func getUsersNeedingSummaryUpdates() []int {
	var users = []int{}
	err := q["getUsersNeedingSummary"].Select(&users)
	if err != nil {
		log.Printf("error executing getUsersNeedingSummary query: %s", err.Error())
	}
	return users
}
