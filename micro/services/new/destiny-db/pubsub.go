package main

import "sync"

var userSubs = map[string][]chan string{}
var userSubsLock sync.RWMutex

func subUser(kind string) <-chan string {
	userSubsLock.Lock()
	defer userSubsLock.Unlock()
	if _, ok := userSubs[kind]; !ok {
		userSubs[kind] = []chan string{}
	}
	c := make(chan string)
	userSubs[kind] = append(userSubs[kind], c)
	return c
}

func pubUser(kind, userID string) {
	userSubsLock.RLock()
	defer userSubsLock.RUnlock()
	if _, ok := userSubs[kind]; !ok {
		return
	}
	for _, v := range userSubs[kind] {
		go func(v chan string, u string) {
			v <- u
		}(v, userID)
	}
}
