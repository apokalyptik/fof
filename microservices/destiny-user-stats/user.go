package main

import "time"

type userBit struct {
	url  string
	when time.Time
	data interface{}
}

type user struct {
	name     string
	platform int
	account  string
	data     map[string]*userBit
	document []byte
}

func (u *user) pull() error {
	//if r, e := client.get(playerURL(platformXBL, "demitriousk")); e != nil {
	//	log.Fatal(e)
	//} else {
	//	log.Printf("%#v", r)
	//}
	if u.platform == 0 {
		// Fill from Player XBL
		// else Fill from Player PSN
	}
	// u.data["grimoire"].url = ...
	// u.data["grimoire"].response = ...
	// u.data["grimoire"].when = time.Now()
	//
	// and so on for:
	//
	// StatsAccount
	// PostGameCarnageReport
	// Triumphs
	// Account
	// Enumerate Characters
	//		ActivityHistory, Activities, Progression, ActivityStats, Stats, StatsUniqueWeapons, Inventory
	// Enumerate Characters Inventiry
	//		InventoryItem
	// Merge into single document
	return nil
}
