package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

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
