package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
)

type filteredUser struct {
	User privateUser
	Seen time.Time
}

type filteredUsers map[string]filteredUser

func filteredUserList() filteredUsers {
	var u = userList
	var s = seenList
	var rval = filteredUsers{}
	var maxAge = time.Now().Add(0 - (30 * 24 * time.Hour))
	for id, t := range s {
		if id == "USLACKBOT" {
			continue
		}
		if t.Before(maxAge) {
			continue
		}
		rval[id] = filteredUser{
			User: u[id],
			Seen: t,
		}
	}
	return rval
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/users.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(filteredUserList())
	})
	go mindSeenList()
	go mindPrivateUserList()
	n := negroni.New()
	n.Use(cors.New(cors.Options{AllowedOrigins: []string{"*"}}))
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(r)
	n.Run("0.0.0.0:8875")
}
