package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var usernameFilter = regexp.MustCompile("[^0-9a-zA-Z-_ ]")

func userList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := http.Get("http://127.0.0.1:8879/users.json")
	if err != nil {
		log.Printf("Error fetching user list: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var details struct {
		Members []struct {
			Bot     bool `json:"is_bot"`
			Deleted bool `json:"deleted"`
			Profile struct {
				FirstName string `json:"first_name"`
			} `json:"profile"`
		} `json:"members"`
	}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&details); err != nil {
		log.Printf("Error unmarshaling user list: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rval := []string{}
	for _, m := range details.Members {
		if m.Bot {
			continue
		}
		if m.Deleted {
			continue
		}
		rval = append(rval, m.Profile.FirstName)
	}
	e := json.NewEncoder(w)
	e.Encode(rval)
}

func memberDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	member := vars["member"]
	w.Header().Set("Content-Type", "application/json")
	var result []struct {
		Type string `bson:"type"`
		Raw  string `bson:"raw"`
	}
	if err := mdb.DB("fof").C("userData").Find(bson.M{"username": member}).Select(bson.M{"_id": -1, "type": 1, "raw": 1}).All(&result); err != nil {
		log.Printf("Error fetching member doc for %s: %s", member, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var found = 0
	for _, v := range result {
		if len(v.Type) > 11 && v.Type == "character-" {
			continue
		}
		if found == 0 {
			fmt.Fprintf(w, "{\"%s\":%s", v.Type, v.Raw)
			found++
		} else {
			fmt.Fprintf(w, ",\"%s\":%s", v.Type, v.Raw)
		}
	}
	fmt.Fprint(w, "}")
}

func memberSubDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	member := usernameFilter.ReplaceAllString(vars["member"], "")
	part := vars["key"]
	w.Header().Set("Content-Type", "application/json")
	var result struct {
		Raw []byte `bson:"raw"`
	}
	err := mdb.DB("fof").C("userData").Find(bson.M{
		"username": member,
		"type":     part,
	}).Select(bson.M{"_id": -1, "raw": 1}).One(&result)
	if err != nil {
		log.Printf("Error fetching member doc.%s for %s: %s", part, member, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result.Raw)
}

func memberSubDocKeys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	member := usernameFilter.ReplaceAllString(vars["member"], "")
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	var result []struct {
		Type string `bson:"type"`
	}
	err := mdb.DB("fof").C(
		"userData").Find(bson.M{"username": member}).Select(bson.M{"type": 1}).All(&result)
	if err != nil {
		log.Printf("Error fetching member doc.keys for %s: %s", member, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var rval = []string{}
	for _, v := range result {
		rval = append(rval, v.Type)
	}
	sort.Strings(rval)
	e.Encode(rval)
}
