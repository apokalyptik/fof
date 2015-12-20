package main

import "sync"

var subs = map[string][]chan string{
	"userAdded":   []chan string{},
	"userRemoved": []chan string{},
}

var subsLock sync.RWMutex

func pubUserAdded(username string) {
	subsLock.RLock()
	defer subsLock.RUnlock()
	for _, v := range subs["userAdded"] {
		go func(v chan string, u string) {
			v <- u
		}(v, username)
	}
}

func subUserAdded() <-chan string {
	subsLock.Lock()
	defer subsLock.Unlock()
	c := make(chan string)
	subs["userAdded"] = append(subs["userAdded"], c)
	return c
}

func pubUserRemoved(username string) {
	subsLock.RLock()
	defer subsLock.RUnlock()
	for _, v := range subs["userRemoved"] {
		go func(v chan string, u string) {
			v <- u
		}(v, username)
	}
}

func subUserRemoved() <-chan string {
	subsLock.Lock()
	defer subsLock.Unlock()
	c := make(chan string)
	subs["userRemoved"] = append(subs["userRemoved"], c)
	return c
}
