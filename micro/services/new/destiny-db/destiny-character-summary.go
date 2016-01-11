package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/apokalyptik/fof/lib/destiny"
)

type charsNeedingSummary []*struct {
	ID          int
	DestinyID   string
	CharacterID string
}

func getCharactersNeedingSummaryUpdates() (charsNeedingSummary, error) {
	var recs charsNeedingSummary
	err := dbSelect(
		&recs,
		`SELECT ID, DestinyID, CharacterID FROM destinyCharacterData JOIN destinyCharacters USING (UserID,CharacterID) JOIN users ON destinyCharacterData.UserID = ID WHERE Type = "Summary" AND Played > Fetched`,
	)
	return recs, err
}

func mindCharacterSummaryUpdates() {
	t := time.Tick(time.Minute)
	for {
		log.Println("mindCharacterSummaryUpdates: waking up")
		exec(`INSERT OR IGNORE INTO destinyCharacterData (UserID,CharacterID,Type) SELECT UserID,CharacterID,"Summary" FROM destinyCharacters`)
		if recs, err := getCharactersNeedingSummaryUpdates(); err != nil {
			log.Println("Error querying for character summaries needing updates:", err.Error())
		} else {
			for _, rec := range recs {
				processCharacterSummary(rec.ID, rec.DestinyID, rec.CharacterID)
			}
		}
		log.Println("mindCharacterSummaryUpdates: sleeping")
		<-t
	}
}

func updateCharacterData(uid int, cid string, kind string, data []byte) error {
	_, err := exec(
		"INSERT OR REPLACE INTO destinyCharacterData (UserID,CharacterID,Type,Raw,Fetched) VALUES(?,?,?,?,datetime(?, 'utc'))",
		uid,
		cid,
		kind,
		data,
		time.Now(),
	)
	return err
}

func handleCharacterSummaryRequest(uid int, cid string, kind string, req *destiny.Request) error {
	var i json.RawMessage
	if err := req.DebugInto(&i); err != nil {
		return err
	}
	if err := updateCharacterData(uid, cid, kind, []byte(i)); err != nil {
		return err
	}
	return nil
}

func processCharacterSummary(uid int, did, cid string) {
	var character = destinyClient.Account(did).Character(cid)

	if req, err := character.CharacterSummary(); err != nil {
		log.Println("error creating CharacterSummary Request:", err.Error())
	} else {
		if err := handleCharacterSummaryRequest(uid, cid, "Summary", req); err != nil {
			log.Println("error handling CharacterSummary request:", err.Error())
		}
	}

	if req, err := character.CharacterInventory(); err != nil {
		log.Println("error creating CharacterInventory Request:", err.Error())
	} else {
		if err := handleCharacterSummaryRequest(uid, cid, "Inventory", req); err != nil {
			log.Println("error handling CharacterInventory request:", err.Error())
		}
	}

	if req, err := character.CharacterActivities(); err != nil {
		log.Println("error creating CharacterActivities Request:", err.Error())
	} else {
		if err := handleCharacterSummaryRequest(uid, cid, "Activities", req); err != nil {
			log.Println("error handling CharacterActivities request:", err.Error())
		}
	}

	if req, err := character.CharacterProgression(); err != nil {
		log.Println("error creating CharacterProgression Request:", err.Error())
	} else {
		if err := handleCharacterSummaryRequest(uid, cid, "Progression", req); err != nil {
			log.Println("error handling CharacterProgression request:", err.Error())
		}
	}

	if req, err := character.AggregateActivityStats(); err != nil {
		log.Println("error creating AggregateActivityStats Request:", err.Error())
	} else {
		if err := handleCharacterSummaryRequest(uid, cid, "AggregateActivityStats", req); err != nil {
			log.Println("error handling AggregateActivityStats request:", err.Error())
		}
	}

	if req, err := character.UniqueWeapons(); err != nil {
		log.Println("error creating UniqueWeapons Request:", err.Error())
	} else {
		if err := handleCharacterSummaryRequest(uid, cid, "UniqueWeapons", req); err != nil {
			log.Println("error handling UniqueWeapons request:", err.Error())
		}
	}

	log.Printf("processCharacterSummary: %s on %s for %d updated", cid, did, uid)
}
