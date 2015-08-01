package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var findManifestURL = "http://www.bungie.net/Platform/Destiny/Manifest/"
var manifestBaseURL = "http://www.bungie.net%s"
var dataPath = "/tmp/bungie-desiny-manifest-db"
var dataTempPath = "/tmp/bungie-desiny-manifest-db.tmp"
var lang = "en"

var listen = "0.0.0.0:8883"

func init() {
	flag.StringVar(&listen, "listen", listen, "address and port number to listen on")
	flag.StringVar(&dataPath, "data", dataPath, "path to store the manifest DB in")
	flag.StringVar(&dataTempPath, "tmp", dataTempPath, "path to temporarily store the new database when updating")
	flag.StringVar(&lang, "lang", lang, "language")
}

func main() {
	flag.Parse()
	manifest.path = dataPath
	manifest.tempPath = dataTempPath
	if err := manifest.update(); err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			if err := manifest.update(); err != nil {
				log.Printf("Error updating manifest: %s", err.Error())
			}
			time.Sleep(time.Hour)
		}
	}()
	r := mux.NewRouter()

	r.HandleFunc("/destiny/types.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var rval = struct {
			OK       bool        `json:"ok"`
			Error    string      `json:"error"`
			Response interface{} `json:"response"`
		}{
			OK:       true,
			Response: tables,
		}
		e := json.NewEncoder(w)
		e.Encode(rval)
	})

	r.HandleFunc("/destiny/{type}/{hash}.json", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		table := vars["type"]
		h := vars["hash"]

		hash, hashErr := overflowString(h)

		var rval = struct {
			OK       bool        `json:"ok"`
			Error    string      `json:"error"`
			Response interface{} `json:"response"`
		}{
			OK: true,
		}

		e := json.NewEncoder(w)

		if hashErr != nil {
			rval.Error = hashErr.Error()
			rval.OK = false
			e.Encode(rval)
			return
		}

		var foundTable = false
		for _, v := range tables {
			if v == table {
				foundTable = true
				break
			}
			if v == fmt.Sprintf("Destiny%sDefinition", table) {
				table = v
				foundTable = true
				break
			}
		}

		if !foundTable {
			rval.Error = "Manifest type not found"
			rval.OK = false
			e.Encode(rval)
			return
		}

		manifest.lock.RLock()
		defer manifest.lock.RUnlock()
		var buf []byte

		var query = fmt.Sprintf("SELECT json FROM %s WHERE id = %d", table, hash)

		if err := manifest.query(query, &buf); err != nil {
			if err != nil {
				rval.Error = err.Error()
				rval.OK = false
			}
		} else {
			if err := json.Unmarshal(buf, &rval.Response); err != nil {
				rval.Error = err.Error()
				rval.OK = false
			}
		}

		e.Encode(rval)
	})
	http.ListenAndServe(listen, r)
}
