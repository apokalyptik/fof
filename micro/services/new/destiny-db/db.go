package main

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
)

var dbLock sync.Mutex
var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Open("sqlite3", "./data.sqlite3.db")
	if err != nil {
		log.Fatalf("error opening database: %s", err.Error())
	}
	db.MapperFunc(strings.Title)
}

func exec(q string, args ...interface{}) (sql.Result, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(q, args...)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return res, err
}

func query(q string, args ...interface{}) (*sql.Rows, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	res, err := tx.Query(q, args...)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return res, err
}

func queryx(q string, args ...interface{}) (*sqlx.Rows, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	res, err := tx.Queryx(q, args...)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return res, err
}

func dbSelect(into interface{}, q string, args ...interface{}) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = tx.Select(into, q, args...)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return err
}
