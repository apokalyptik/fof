package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var xhrOutput *Json
var lfgOutput *Json

func init() {
	xhrOutput = &Json{
		data: map[string]interface{}{
			"channels": []string{},
		},
		cond: &sync.WaitGroup{},
	}
	xhrOutput.cond.Add(1)
	xhrOutput.set("updated_at", time.Now().Unix())

	lfgOutput = &Json{
		data: map[string]interface{}{
			"lfg": []string{},
		},
		cond: &sync.WaitGroup{},
	}
	lfgOutput.cond.Add(1)
	lfgOutput.set("updated_at", time.Now().Unix())
}

type Json struct {
	data      map[string]interface{}
	updatedAt string
	cache     []byte
	cond      *sync.WaitGroup
	lock      sync.RWMutex
}

func (j *Json) send(w http.ResponseWriter) error {
	j.lock.RLock()
	_, err := w.Write(j.cache)
	j.lock.RUnlock()
	return err
}

func (j *Json) set(key string, value interface{}) error {
	j.lock.Lock()
	j.updatedAt = fmt.Sprintf("%d", time.Now().Unix())
	j.data["updated_at"] = j.updatedAt
	j.data[key] = value
	cache, err := json.Marshal(j.data)
	if err != nil {
		log.Println(err.Error())
		j.lock.Unlock()
		return err
	}
	j.cache = cache
	newCond := &sync.WaitGroup{}
	newCond.Add(1)
	j.cond.Done()
	j.cond = newCond
	j.lock.Unlock()
	return nil
}
