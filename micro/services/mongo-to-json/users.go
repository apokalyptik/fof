package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var usernameFilter = regexp.MustCompile("[^0-9a-zA-Z-_ ]")

var userListCache []byte
var userListCacheAge time.Time

func seenList() (map[string]time.Time, error) {
	var recent = map[string]time.Time{}
	if resp, err := http.Get("http://fofgaming.com:8890/seen.json"); err != nil {
		log.Println("error fetching seen.json:", err.Error())
		return recent, err
	} else {
		defer resp.Body.Close()
		var data = map[string]time.Time{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&data); err != nil {
			log.Println("error decoding seen.json:", err.Error())
			return recent, err
		}
		maxAge := time.Now().Add(0 - (30 * 24 * time.Hour))
		for id, t := range data {
			if t.After(maxAge) {
				recent[id] = t
			}
		}
	}
	return recent, nil
}

func getUserlist() []byte {
	if userListCacheAge.After(time.Now().Add(0 - time.Hour)) {
		return userListCache
	}
	resp, err := http.Get("http://127.0.0.1:8879/users.json")
	if err != nil {
		log.Printf("Error fetching user list: %s", err.Error())
		return userListCache
	}
	defer resp.Body.Close()

	recent, recentErr := seenList()

	var details struct {
		Members []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Bot     bool   `json:"is_bot"`
			Deleted bool   `json:"deleted"`
			Profile struct {
				FirstName string `json:"first_name"`
				Avatar    string `json:"image_192"`
			} `json:"profile"`
		} `json:"members"`
	}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&details); err != nil {
		log.Printf("Error unmarshaling user list: %s", err.Error())
		return userListCache
	}

	var destinyLookup = []struct {
		UserID         string `bson:"userid"`
		DestinyAccount string `bson:"account"`
	}{}
	if err := mdb.DB("fof").C("userLookup").Find(nil).Select(bson.M{"_id": -1, "userid": 1, "account": 1}).All(&destinyLookup); err != nil {
		log.Printf("Error fetching destiny lookup docs for: %s", err.Error())
		return userListCache
	}

	rval := []map[string]string{}
	for _, m := range details.Members {
		if m.Bot {
			continue
		}
		if m.Deleted {
			continue
		}
		if recentErr != nil {
			if _, ok := recent[m.ID]; !ok {
				continue
			}
		}
		var destiny = ""
		for _, dUser := range destinyLookup {
			if dUser.UserID == m.ID {
				destiny = dUser.DestinyAccount
			}
		}
		rval = append(rval, map[string]string{
			"gamertag": m.Profile.FirstName,
			"username": m.Name,
			"avatar":   m.Profile.Avatar,
			"destiny":  destiny,
		})
	}
	if encoded, err := json.Marshal(rval); err != nil {
		return userListCache
	} else {
		userListCache = encoded
		userListCacheAge = time.Now()
		return encoded
	}
}

func userList(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	setJSON(w, r)
	w.Write(getUserlist())
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
