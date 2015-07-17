package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mgoHost = "127.0.0.1"
var listenOn = "0.0.0.0:8880"

var mdb *mgo.Session

func exoticStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rval interface{}
	_, err := mdb.DB("fof").C("memberDestinyStats").Find(nil).MapReduce(&mgo.MapReduce{
		Map: `function() {
    var rval = { 
      member: this.member,
      uniqueWeapons: {},
    };

    for ( var cidx in this.data.characters ) {
      var char = this.data.characters[cidx];
      for ( var didx in char.uniqueWeapons.definitions.items ) {
        var itemHash = char.uniqueWeapons.definitions.items[didx].itemHash;
        var itemName = char.uniqueWeapons.definitions.items[didx].itemName;
        if ( typeof rval.uniqueWeapons[itemName] == "undefined" ) {
          rval.uniqueWeapons[itemName] = {
            total: 0, 
            normal: 0, 
            precision: 0,
          }
        }
        var item = null;
        for ( var iidx=0; iidx<char.uniqueWeapons.data.weapons.length; iidx++ ) {
          if ( char.uniqueWeapons.data.weapons[iidx].referenceId == itemHash ) {
            item = char.uniqueWeapons.data.weapons[iidx].values;
            break;
          }
        }
        if ( item == null ) {
          continue;
        }
        rval.uniqueWeapons[itemName].total += item.uniqueWeaponKills.basic.value;
        rval.uniqueWeapons[itemName].precision += item.uniqueWeaponKills.basic.value * item.uniqueWeaponKillsPrecisionKills.basic.value;
        rval.uniqueWeapons[itemName].normal = rval.uniqueWeapons[itemName].total - rval.uniqueWeapons[itemName].precision;
      }
    }
    emit(this.member, rval)
  }`,
		Reduce: `function(killstats) {
    rval = [];
    for ( var i=0; i<members.length; i++ ) {
      rval.push({member: members[i], stats: killstats[i]});
    }
    return rval;
  }`,
	}, &rval)
	if err != nil {
		log.Printf("Error querying MongoDB for exoticStats: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e := json.NewEncoder(w)
	if err := e.Encode(rval); err != nil {
		log.Printf("Error marshaling json for exoticStats: %s", err.Error())
	}
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
	filter := regexp.MustCompile("[^0-9a-zA-Z-_]")
	vars := mux.Vars(r)
	member := filter.ReplaceAllString(vars["member"], "")
	part := vars["part"]
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	var result interface{}
	err := mdb.DB("fof").C(
		"memberDestinyStats").Find(
		bson.M{"member": member}).Select(
		bson.M{fmt.Sprintf("data.%s", part): 1}).One(&result)
	if err != nil {
		log.Printf("Error fetching member doc.%s for %s: %s", part, member, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e.Encode(result)
}

func init() {
	flag.StringVar(&mgoHost, "mgo", mgoHost, "MongoDB Address")
	flag.StringVar(&listenOn, "listen", listenOn, "HTTP Server")
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
	r.HandleFunc("/destiny/members/exotic-kill-stats.json", exoticStats)
	r.HandleFunc("/destiny/members/get/{member}.json", memberDoc)
	r.HandleFunc("/destiny/members/get/{member}/{part}.json", memberSubDoc)
	log.Fatal(http.ListenAndServe(listenOn, r))
}
