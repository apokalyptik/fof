package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type need struct {
	Until time.Time `json:"until"`
	User  string    `json:"user"`
}

type needs struct {
	lock     sync.RWMutex
	filename string
	DB       map[string]map[string][]need `json:"db"`
	Version  int                          `json:"version"`
}

var needsDB = &needs{
	DB: map[string]map[string][]need{},
}

func (n *needs) load(filename string) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			n.filename = filename
			return nil
		}
		return err
	}
	n.filename = filename
	defer fp.Close()
	dec := json.NewDecoder(fp)
	if err := dec.Decode(n); err != nil {
		return err
	}
	if n.DB == nil {
		n.DB = map[string]map[string][]need{}
	}
	return nil
}

func (n *needs) save() error {
	if n.filename == "" {
		return errors.New("cannot persist withuot filename")
	}
	fp, err := os.Create(n.filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if err := enc.Encode(n); err != nil {
		return err
	}
	return nil
}

func (n *needs) mindExpiration() {
}

func (n *needs) add(channel, user, wants string, until time.Duration) error {
	return nil
}

func (n *needs) del(channl, user, wants string) error {
	return nil
}

func (n *needs) list(channel string) error {
	return nil
}
