package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var dbPath = "./data.sqlite3"
var sqlMessage = make(chan slackmsg, 1024)

func mindSQL() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS public_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			channel_id TEXT,
			channel_name TEXT,
			user_id TEXT,
			user_name TEXT,
			timestamp TEXT,
			message TEXT
		);
		CREATE INDEX IF NOT EXISTS public_messages_channel_id on public_messages (channel_id);
		CREATE INDEX IF NOT EXISTS public_messages_channel_name on public_messages (channel_name);
		CREATE INDEX IF NOT EXISTS public_messages_user_id on public_messages (user_id);
		CREATE INDEX IF NOT EXISTS public_messages_user_name on public_messages (user_name);
		CREATE INDEX IF NOT EXISTS public_messages_timestamp on public_messages (timestamp);
		`); err != nil {
		log.Fatalf("Error creating `public_messages` table: %s", err.Error())
	}
	insert, err := db.Prepare(
		"INSERT INTO public_messages (channel_id,channel_name,user_id,user_name,timestamp,message) VALUES(?,?,?,?,?,?)",
	)
	if err != nil {
		log.Fatalf("Error preparing insert statement: %s", err.Error())
	}
	go func(insert *sql.Stmt) {
		for {
			select {
			case m := <-sqlMessage:
				if _, err := insert.Exec(m.ChannelID, m.Channel, m.UserID, m.User, m.Timestamp, m.Text); err != nil {
					log.Fatalf("Error inserting sql row: %s", err.Error())
				}
			}
		}
	}(insert)
}
