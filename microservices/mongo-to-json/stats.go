package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func allTimeStats(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	setJSON(w, r)
	vars := mux.Vars(r)
	query := bson.M{
		"section": vars["section"],
		"stat":    vars["stat"],
	}
	var docs []struct {
		Member string  `json:"member",bson:"member"`
		Value  float64 `json:"value",bson:"value"`
		PGA    string  `json:"pga",bson:"pgadisplayvalue"`
	}
	if err := mdb.DB("fof").C("accountStats").Find(query).Limit(10000).Sort("-value", "-pgavalue").All(&docs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error fetching allTimeStatKeys: %s", err.Error())
		return
	}
	e := json.NewEncoder(w)
	e.Encode(docs)
}

func allTimeStatKeys(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	setJSON(w, r)

	var rval = map[string][]string{}
	var sections = []string{}

	if err := mdb.DB("fof").C("accountStats").Find(nil).Distinct("section", &sections); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error fetching allTimeStatKeys: %s", err.Error())
		return
	}
	for _, s := range sections {
		var sectionStats = []string{}
		if err := mdb.DB("fof").C("accountStats").Find(&bson.M{"section": s}).Distinct("stat", &sectionStats); err != nil {
			continue
		}
		rval[s] = sectionStats
	}
	e := json.NewEncoder(w)
	e.Encode(rval)
	return
	query := bson.M{
		"member": "Adm Wright Meow",
	}
	var docs []struct {
		Section string `bson:"section"`
		Stat    string `bson:"stat"`
	}
	if err := mdb.DB("fof").C("accountStats").Find(query).Limit(1000).All(&docs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error fetching allTimeStatKeys: %s", err.Error())
		return
	}
	for _, doc := range docs {
		if _, ok := rval[doc.Section]; !ok {
			rval[doc.Section] = []string{}
		}
		rval[doc.Section] = append(rval[doc.Section], doc.Stat)
	}
}
