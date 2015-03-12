package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var raidDb = &raids{
	data: map[string][]*raid{},
}

type raid struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Members   []string  `json:"members"`
}

type raids struct {
	filename string
	data     map[string][]*raid
	lock     sync.RWMutex
}

func (r *raids) join(channel, name, user string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		return errors.New(fmt.Sprintf("There are no raids for #%s", channel))
	}
	if len(r.data[channel]) == 0 {
		return errors.New(fmt.Sprintf("There are no raids for #%s", channel))
	}
	for _, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		for _, n := range v.Members {
			if n == user {
				return errors.New(fmt.Sprintf(
					"You have already signed up for \"%s\" on #%s",
					name,
					channel))
			}
		}
		slack.toPerson(v.Members[0], fmt.Sprintf(
			"@%s has joined your raid \"%s\" on #%s", user, name, channel))
		v.Members = append(v.Members, user)
		return nil
	}
	return errors.New(fmt.Sprintf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \"/raid host\" one?",
		name,
		channel))
}

func (r *raids) leave(channel, name, user string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		return errors.New(fmt.Sprintf("There are no raids for #%s", channel))
	}
	if len(r.data[channel]) == 0 {
		return errors.New(fmt.Sprintf("There are no raids for #%s", channel))
	}
	for vid, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		for k, n := range v.Members {
			if n != user {
				continue
			}
			v.Members = append(v.Members[:k], v.Members[k+1:]...)
			if len(v.Members) == 0 {
				r.data[channel] = append(r.data[channel][:vid], r.data[channel][vid+1:]...)
				return errors.New(fmt.Sprintf(
					"Since you were the last member of \"%s\" on #%s the raid has been disbanded",
					name,
					channel))
			}
			slack.toPerson(v.Members[0], fmt.Sprintf(
				"@%s has left your raid \"%s\" on #%s", user, name, channel))
			return nil
		}
		return errors.New(fmt.Sprintf(
			"You are not signed up to do \"%s\" on #%s",
			name,
			channel))
	}
	return errors.New(fmt.Sprintf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \"/raid host\" one?",
		name,
		channel))
}

func (r *raids) members(channel, name string) ([]string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if _, ok := r.data[channel]; !ok {
		return nil, errors.New(fmt.Sprintf("There are no raids for #%s", channel))
	}
	if len(r.data[channel]) == 0 {
		return nil, errors.New(fmt.Sprintf("There are no raids for #%s", channel))
	}
	for _, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		rval := make([]string, len(v.Members))
		for k, n := range v.Members {
			rval[k] = n
		}
		return rval, nil
	}
	return nil, errors.New(fmt.Sprintf(
		"I have no \"%s\" registered for #%s",
		name,
		channel))
}

func (r *raids) finish(channel, name, user string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	for k, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		var allowed = false
		if v.Members[0] == user {
			allowed = true
		} else {
			for _, admin := range admins {
				if user == admin {
					allowed = true
					break
				}
			}
		}
		if allowed == false {
			return errors.New(fmt.Sprintf(
				"Only the organizer (_@%s_) can finish a raid",
				v.Members[0]))
		}
		r.data[channel] = append(r.data[channel][:k], r.data[channel][k+1:]...)
		return nil
	}
	return errors.New(fmt.Sprintf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \"/raid host\" one?",
		name,
		channel))
}

func (r *raids) register(channel, name, user string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		r.data[channel] = []*raid{}
	}
	for _, v := range r.data[channel] {
		if v.Name == name {
			return errors.New("A raid by this name is already registered")
		}
	}
	r.data[channel] = append(r.data[channel], &raid{
		Name:      name,
		CreatedAt: time.Now(),
		Members:   []string{user},
	})
	return nil
}

func (r *raids) list(channel string) []raid {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if c, ok := r.data[channel]; ok {
		rval := make([]raid, len(c))
		for k, v := range c {
			newraid := raid{
				Name:      v.Name,
				CreatedAt: v.CreatedAt,
				Members:   make([]string, len(v.Members)),
			}
			for mk := range v.Members {
				newraid.Members[mk] = v.Members[mk]
			}
			rval[k] = newraid
		}
		return rval
	}
	return nil
}

func (r *raids) save() error {
	if r.filename == "" {
		return errors.New("cannot persist withuot filename")
	}
	fp, err := os.Create(r.filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if err := enc.Encode(r.data); err != nil {
		return err
	}
	return nil
}

func (r *raids) mindExpiration(maxAge time.Duration) {
	ticker := time.Tick(10 * time.Minute)
	for {
		<-ticker
		r.lock.RLock()
		for channel, raidlist := range r.data {
			for _, raidentry := range raidlist {
				if time.Now().Add(0 - maxAge).After(raidentry.CreatedAt) {
					go r.finish(channel, raidentry.Name, raidentry.Members[0])
					log.Printf("Expiring %s on #%s", raidentry.Name, channel)
				}
			}
		}
		r.lock.RUnlock()

	}
}

func (r *raids) load(filename string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			r.filename = filename
			return nil
		}
		return err
	}
	r.filename = filename
	defer fp.Close()
	dec := json.NewDecoder(fp)
	if err := dec.Decode(&r.data); err != nil {
		return err
	}
	if r.data == nil {
		r.data = map[string][]*raid{}
	}
	return nil
}
