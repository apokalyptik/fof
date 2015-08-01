package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2"
)

func pvpTotal(w http.ResponseWriter, r *http.Request) {
	cleaner := regexp.MustCompile("[^0-9a-zA-Z_-]")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	stat := cleaner.ReplaceAllString(vars["key"], "")

	var rval interface{}
	var cleaned = []struct {
		Member string `mapstructure:"_id"`
		Value  map[string]struct {
			PGA   float64 `mapstructure:"pga"`
			Total int64   `mapstructure:"total"`
		} `mapstructure:"value"`
	}{}
	_, err := mdb.DB("fof").C("memberDestinyStats").Find(nil).MapReduce(&mgo.MapReduce{
		Map: fmt.Sprintf(`function() {
				if ( typeof this.data.accountStats != "object" ) {
					return
				}
				var rval = {};
				for( k in this.data.accountStats.mergedAllCharacters.results.allPvP.allTime ) {
					if ( k != "%s" ) {
						continue;
					}
					name = k;
					data = this.data.accountStats.mergedAllCharacters.results.allPvP.allTime[k];
					pga = null;
					if ( typeof data.pga != "undefined" ) {
						pga = data.pga.value;
					}
					rval[k] = { total: data.basic.value, pga: pga } 
				}
				emit( this.member, rval ); 
			}`, stat),
		Reduce: `function(stats) {
			rval = [];
			for ( var i=0; i<stats.length; i++ ) {
			  rval.push(stats[i]);
			}
			return rval;
		  }`,
	}, &rval)
	if err != nil {
		log.Printf("Error querying MongoDB for pvpTotal: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := mapstructure.Decode(rval, &cleaned); err != nil {
		log.Printf("Error decoding rval for pvpTotal: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var final = map[string]interface{}{}
	for _, v := range cleaned {
		if _, ok := v.Value[stat]; ok {
			s := v.Value[stat]
			final[v.Member] = struct {
				Total int64   `json:"total"`
				PGA   float64 `json:"pga"`
			}{
				Total: s.Total,
				PGA:   s.PGA,
			}
		}
	}
	e := json.NewEncoder(w)
	if err := e.Encode(final); err != nil {
		log.Printf("Error marshaling json for pvpTotals: %s", err.Error())
	}
}

func pvpTotalsKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rval interface{}
	var converted = []struct {
		ID   string   `mapstructure:"_id"`
		Keys []string `mapstructure:"value"`
	}{}
	_, err := mdb.DB("fof").C("memberDestinyStats").Find(nil).MapReduce(&mgo.MapReduce{
		Map: `function() {
				if ( typeof this.data.accountStats != "object" ) {
					return
				}
				var rval = [];
				for( k in this.data.accountStats.mergedAllCharacters.results.allPvP.allTime ) {
					rval.push(k)
				}
				emit( this.member, rval ); 
			}`,
		Reduce: `function(stats) {
			rval = [];
			for ( var i=0; i<stats.length; i++ ) {
			  var found = false;
			  for ( var n=0; n<rval.length; n++ ) {
				  if ( rval[n] == stats[i] ) {
					  found = true;
					  break;
				  }
			  }
			  if ( found == false ) {
				  rval.push(typeof stats);
			  }
			}
			return rval;
		  }`,
	}, &rval)
	if err != nil {
		log.Printf("Error querying MongoDB for pvpKeys: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := mapstructure.Decode(rval, &converted); err != nil {
		log.Printf("Error converting query result for pvpKeys: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var aggregate = map[string]struct{}{}
	var final = []string{}
	for _, m := range converted {
		for _, v := range m.Keys {
			if _, ok := aggregate[v]; !ok {
				final = append(final, v)
				aggregate[v] = struct{}{}
			}
		}
	}
	sort.Strings(final)
	e := json.NewEncoder(w)
	if err := e.Encode(final); err != nil {
		log.Printf("Error marshaling json for pvpKeys: %s", err.Error())
	}
}

func pvpTotals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rval interface{}
	_, err := mdb.DB("fof").C("memberDestinyStats").Find(nil).MapReduce(&mgo.MapReduce{
		Map: `function() {
				if ( typeof this.data.accountStats != "object" ) {
					return
				}
				var rval = {};
				for( k in this.data.accountStats.mergedAllCharacters.results.allPvP.allTime ) {
					name = k;
					data = this.data.accountStats.mergedAllCharacters.results.allPvP.allTime[k];
					pga = null;
					if ( typeof data.pga != "undefined" ) {
						pga = data.pga.value;
					}
					rval[k] = { value: data.basic.value, pga: pga } 
				}
				emit( this.member, rval ); 
			}`,
		Reduce: `function(stats) {
			rval = [];
			for ( var i=0; i<stats.length; i++ ) {
			  rval.push(stats[i].value);
			}
			return rval;
		  }`,
	}, &rval)
	if err != nil {
		log.Printf("Error querying MongoDB for pvpTotals: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e := json.NewEncoder(w)
	if err := e.Encode(rval); err != nil {
		log.Printf("Error marshaling json for pvpTotals: %s", err.Error())
	}
}
