package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var mgoHost = "127.0.0.1"
var listenOn = "0.0.0.0:8880"

var mdb *mgo.Session

func init() {
	flag.StringVar(&mgoHost, "mgo", mgoHost, "MongoDB Address")
	flag.StringVar(&listenOn, "listen", listenOn, "HTTP Server")
}

func setJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func setCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	flag.Parse()
	if session, err := mgo.Dial(mgoHost); err != nil {
		log.Fatalf("Error dialing mongodb %s: %s", mgoHost, err.Error())
	} else {
		mdb = session
	}
	log.Println("Starting up...")
	r := mux.NewRouter()
	r.HandleFunc("/fof/members.json", userList)
	r.HandleFunc("/destiny/raw/{member}.json", memberDoc)
	r.HandleFunc("/destiny/raw/{member}/keys.json", memberSubDocKeys)
	r.HandleFunc("/destiny/raw/{member}/{key}.json", memberSubDoc)
	r.HandleFunc("/destiny/stats/alltime/keys.json", allTimeStatKeys)
	r.HandleFunc("/destiny/stats/alltime/{section}/{stat}.json", allTimeStats)
	r.HandleFunc("/destiny/pva/exotic-kills.json", exoticStats)
	r.HandleFunc("/destiny/pvp/allTime/aggregate.json", pvpTotals)
	r.HandleFunc("/destiny/pvp/allTime/aggregate/keys.json", pvpTotalsKeys)
	r.HandleFunc("/destiny/pvp/allTime/aggregate/{key}.json", pvpTotal)
	log.Fatal(http.ListenAndServe(listenOn, r))
}
