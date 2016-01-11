package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var q = map[string]*sqlx.Stmt{}
var qs = map[string]string{}

func initSQL() {
	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID TEXT NOT NULL DEFAULT "",
			UserName TEXT NOT NULL DEFAULT "",
			GamerTag TEXT NOT NULL DEFAULT "",
			DestinyID TEXT NOT NULL DEFAULT "",
			Seen DATETIME  NOT NULL DEFAULT "0000-00-00 00:00:00",
			DestinyChecked INTEGER not null default 0
		);
		CREATE UNIQUE INDEX IF NOT EXISTS users_uid ON users (UserID);
		CREATE INDEX IF NOT EXISTS users_did ON users (DestinyID);
		CREATE INDEX IF NOT EXISTS users_seen ON users (Seen);
		CREATE INDEX IF NOT EXISTS users_gt ON users (GamerTag);
	`); err != nil {
		log.Fatalf("error creating users table: %s", err.Error())
	}

	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS userMeta (
			UserID INTEGER NOT NULL DEFAULT 0,
			Key TEXT NOT NULL DEFAULT "",
			Value TEXT NOT NULL DEFAULT ""
		);
		CREATE UNIQUE INDEX IF NOT EXISTS usermeta_lookup ON userMeta (UserID,Key);
		CREATE INDEX IF NOT EXISTS usermeta_userid ON userMeta (UserID);
		CREATE INDEX IF NOT EXISTS usermeta_key ON userMeta (Key);
	`); err != nil {
		log.Fatalf("error creating usermeta table: %s", err.Error())
	}

	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS destinyCharacters (
			UserID INTEGER NOT NULL DEFAULT 0,
			CharacterID STRING NOT NULL DEFAULT "",
			Played DATETIME NOT NULL DEFAULT "0000-00-00 00:00:00"
		);
		CREATE UNIQUE INDEX IF NOT EXISTS destinyCharacters_lookup ON destinyCharacters (UserID);
		CREATE INDEX IF NOT EXISTS destinyCharacters_when ON destinyCharacters (Played);
		CREATE INDEX IF NOT EXISTS destinyCharacters_cid ON destinyCharacters (CharacterID);
	`); err != nil {
		log.Fatalf("error creating destinyCharacters table: %s", err.Error())
	}

	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS destinyAccountSummary (
			UserID INTEGER NOT NULL DEFAULT 0,
			Fetched DATETIME NOT NULL DEFAULT "0000-00-00 00:00:00",
			Raw BLOB NOT NULL DEFAULT ""
		);
		CREATE UNIQUE INDEX IF NOT EXISTS destinyAccountSummary_lookup ON destinyAccountSummary (UserID);
		CREATE INDEX IF NOT EXISTS destinyAccountSummary_when ON destinyAccountSummary (Fetched);
	`); err != nil {
		log.Fatalf("error creating destinyAccountSummary table: %s", err.Error())
	}

	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS destinyAccountValues (
			UserID INTEGER  NOT NULL DEFAULT 0,
			Key TEXT NOT NULL DEFAULT "",
			Value TEXT NOT NULL DEFAULT ""
		);
		CREATE UNIQUE INDEX IF NOT EXISTS destinyAccountValues_lookup ON destinyAccountValues (UserID,Key);
		CREATE INDEX IF NOT EXISTS destinyAccountValues_key ON destinyAccountValues (Key);
	`); err != nil {
		log.Fatalf("error creating destinyAccountValues table: %s", err.Error())
	}

	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS destinyAccountStats (
			UserID INTEGER  NOT NULL DEFAULT 0,
			Key TEXT NOT NULL DEFAULT "",
			Value INTEGER NOT NULL DEFAULT 0
		);
		CREATE UNIQUE INDEX IF NOT EXISTS destinyAccountStats_lookup ON destinyAccountStats (UserID,Key);
		CREATE INDEX IF NOT EXISTS destinyAccountStats_key ON destinyAccountStats (Key);
	`); err != nil {
		log.Fatalf("error creating destinyAccountStats table: %s", err.Error())
	}

	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS destinyAccountDailyStats (
			UserID INTEGER  NOT NULL DEFAULT 0,
			Day DATE NOT NULL DEFAULT "0000-00-00",
			Key TEXT NOT NULL DEFAULT "",
			Value INTEGER NOT NULL DEFAULT 0
		);
		CREATE UNIQUE INDEX IF NOT EXISTS destinyAccountDailyStats_lookup ON destinyAccountDailyStats (UserID,Key,Day);
		CREATE INDEX IF NOT EXISTS destinyAccountDailyStats_key ON destinyAccountDailyStats (Key);
		CREATE INDEX IF NOT EXISTS destinyAccountDailyStats_day ON destinyAccountDailyStats (Day);
	`); err != nil {
		log.Fatalf("error creating destinyAccountDailyStats table: %s", err.Error())
	}

	if _, err := exec(`
		CREATE TABLE IF NOT EXISTS destinyCharacterData (
			UserID INTEGER  NOT NULL DEFAULT 0,
			Fetched DATETIME NOT NULL DEFAULT "0000-00-00 00:00:00",
			CharacterID TEXT NOT NULL DEFAULT "",
			Type TEXT NOT NULL DEFAULT "",
			Raw BLOB NOT NULL DEFAULT ""
		);
		CREATE UNIQUE INDEX IF NOT EXISTS destinyCharacterData_lookup ON destinyCharacterData (UserID,CharacterID,Type);
		CREATE INDEX IF NOT EXISTS destinyCharacterData_char ON destinyCharacterData (CharacterID);
		CREATE INDEX IF NOT EXISTS destinyCharacterData_type ON destinyCharacterData (Type);
		CREATE INDEX IF NOT EXISTS destinyCharacterData_when ON destinyCharacterData (Fetched);
	`); err != nil {
		log.Fatalf("error creating destinyCharacterData table: %s", err.Error())
	}

	for k, v := range qs {
		if qk, err := db.Preparex(v); err != nil {
			log.Fatalf("Error preparing the query for %s: %s", k, err.Error())
		} else {
			q[k] = qk
		}
	}
}
