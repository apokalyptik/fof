package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type statDesc struct {
	ID       int
	Platform string
	Product  string
	Stat     string
	Sub1     string
	Sub2     string
	Sub3     string
	Info     map[string]interface{}
}

type statList []statDesc

func handleStatsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rval, err := getStatsList()
	if err != nil {
		log.Println("Error getting stats list:", err.Error())
	}
	json.NewEncoder(w).Encode(rval)
}

func getStatsList() (statList, error) {
	var rval statList
	rows, err := db.Query("SELECT * FROM `stats`")
	if err != nil {
		return rval, err
	}
	defer rows.Close()
	for rows.Next() {
		var s statDesc
		var info []byte
		err = rows.Scan(
			&s.ID,
			&s.Platform,
			&s.Product,
			&s.Stat,
			&s.Sub1,
			&s.Sub2,
			&s.Sub3,
			&info,
		)
		if err != nil {
			break
		}
		if err = json.Unmarshal(info, &s.Info); err != nil {
			s.Info = map[string]interface{}{}
			err = nil
		}
		rval = append(rval, s)
	}
	return rval, err
}
