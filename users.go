package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var udb = &userDB{}

type userDB struct {
	filename    string
	lock        sync.RWMutex
	LastUpdated time.Time
	LookupID    map[string]string `json:"LookupUserToID"`
	IMs         map[string]string `json:"LookupUserToIM"`
}

func (u *userDB) save() error {
	if u.filename == "" {
		return errors.New("cannot persist withuot filename")
	}
	fp, err := os.Create(u.filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if err := enc.Encode(u); err != nil {
		return err
	}
	return nil
}

func (u *userDB) load(filename string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			u.filename = filename
			return nil
		}
		return err
	}
	u.filename = filename
	defer fp.Close()
	dec := json.NewDecoder(fp)
	if err := dec.Decode(&u); err != nil {
		return err
	}
	if u == nil {
		u = &userDB{
			LookupID: map[string]string{},
			IMs:      map[string]string{},
		}
	}
	return nil
}
