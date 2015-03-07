package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

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
	r.lock.RLock()
	defer r.lock.RUnlock()
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
	if r.data == nil {
		r = &raids{
			data: map[string][]*raid{},
		}
	}
	if err := dec.Decode(&r.data); err != nil {
		return err
	}
	return nil
}
