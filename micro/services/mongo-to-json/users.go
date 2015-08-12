package main

import (
	"encoding/json"
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
	e := json.NewEncoder(w)
	var result interface{}
	if err := mdb.DB("fof").C("memberDestinyStats").Find(bson.M{"member": member}).One(&result); err != nil {
		log.Printf("Error fetching member doc for %s: %s", member, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e.Encode(result)
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
	var result map[string]map[string]interface{}
	err := mdb.DB("fof").C(
		"memberDestinyStats").Find(bson.M{"member": member}).One(&result)
	if err != nil {
		log.Printf("Error fetching member doc.keys for %s: %s", member, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var rval = []string{}
	for k, _ := range result["data"] {
		rval = append(rval, k)
	}
	sort.Strings(rval)
	e.Encode(rval)
}
