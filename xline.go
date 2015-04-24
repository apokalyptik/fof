package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

type sticker struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type xline struct {
	lock     sync.RWMutex
	filename string
	DB       []*sticker `json:"db"`
	Version  int        `json:"version"`
}

var xlineDB = &xline{
	DB: []*sticker{},
}

func (x *xline) search(what string) []string {
	x.lock.RLock()
	defer x.lock.RUnlock()
	var rval = []string{}
	for _, v := range x.DB {
		if strings.Contains(v.Name, what) || strings.Contains(v.Description, what) {
			rval = append(rval, fmt.Sprintf("%s \"%s\"", v.Name, v.Description))
		}
	}
	return rval
}

func (x *xline) remove(sticker string) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	defer x.save()
	for i, s := range x.DB {
		if s.Name == sticker {
			x.DB = append(x.DB[:i], x.DB[i+1:]...)
			return nil
		}
	}
	return errors.New("Sticker not found")
}

func (x *xline) add(name, url, description string) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	defer x.save()
	for _, v := range x.DB {
		if v.Name == name {
			v.URL = url
			v.Description = description
			return nil
		}
	}
	x.DB = append(x.DB, &sticker{
		Name:        name,
		URL:         url,
		Description: description,
	})
	return nil
}

func (x *xline) get(what string) (string, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()
	for _, v := range x.DB {
		if v.Name == what {
			return v.URL, nil
		}
	}
	return "", errors.New("sticker not found")
}

func (x *xline) load(filename string) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			x.filename = filename
			return nil
		}
		return err
	}
	x.filename = filename
	defer fp.Close()
	dec := json.NewDecoder(fp)
	if err := dec.Decode(x); err != nil {
		return err
	}
	if x.DB == nil {
		x.DB = []*sticker{}
	}
	return nil
}

func (x *xline) save() error {
	if x.filename == "" {
		return errors.New("cannot persist withuot filename")
	}
	fp, err := os.Create(x.filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if err := enc.Encode(x); err != nil {
		return err
	}
	return nil
}
