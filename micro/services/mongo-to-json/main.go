package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"

	"github.com/NYTimes/gziphandler"
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

func middleware(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return gziphandler.GzipHandler(
		http.HandlerFunc(f),
	)
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
	r.Handle("/fof/members.json", middleware(userList))
	r.Handle("/destiny/raw/{member}.json", middleware(memberDoc))
	r.Handle("/destiny/raw/{member}/keys.json", middleware(memberSubDocKeys))
	r.Handle("/destiny/raw/{member}/{key}.json", middleware(memberSubDoc))
	r.Handle("/destiny/stats/alltime/keys.json", middleware(allTimeStatKeys))
	r.Handle("/destiny/stats/alltime/{section}/{stat}.json", middleware(allTimeStats))
	r.Handle("/destiny/pva/exotic-kills.json", middleware(exoticStats))
	r.Handle("/destiny/pvp/allTime/aggregate.json", middleware(pvpTotals))
	r.Handle("/destiny/pvp/allTime/aggregate/keys.json", middleware(pvpTotalsKeys))
	r.Handle("/destiny/pvp/allTime/aggregate/{key}.json", middleware(pvpTotal))
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Fatal(http.ListenAndServe(listenOn, r))
}
