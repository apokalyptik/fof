package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type usersHourlyStats map[string]userHourlyStatsList
type userHourlyStatsList map[string]userHourlyStats
type userHourlyStats map[string]int

func getUserHourlyStats(user, stat string, last int) (userHourlyStats, error) {
	var rval = userHourlyStats{}
	var err error
	var rows *sql.Rows
	var stat_id int
	stat_id, err = strconv.Atoi(stat)
	if err != nil {
		return rval, err
	}
	if last >= 0 {
		rows, err = db.Query(
			"SELECT `when`,`value` FROM stats_hourly WHERE `member`=? AND `stat_id`=? AND `when` >= DATE_SUB(NOW(), INTERVAL ? HOUR)",
			user,
			stat_id,
			last)
	} else {
		rows, err = db.Query(
			"SELECT `when`,`value` FROM stats_hourly WHERE `member`=? AND `stat_id`=?",
			user,
			stat_id)
	}
	if err != nil {
		return rval, err
	}
	defer rows.Close()
	for rows.Next() {
		var when string
		var value int
		err = rows.Scan(&when, &value)
		if err != nil {
			break
		}
		rval[when] = value
	}
	return rval, err
}

func handleHourlyJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users := strings.Split(r.URL.Query().Get("users"), ",")
	if len(users) < 1 {
		return
	}
	stats := strings.Split(r.URL.Query().Get("stats"), ",")
	if len(stats) < 1 {
		return
	}
	last, err := strconv.Atoi(r.URL.Query().Get("last"))
	if err != nil {
		last = -1
	}
	switch vars["type"] {
	case "json":
		var rval = usersHourlyStats{}
		for _, u := range users {
			rval[u] = userHourlyStatsList{}
			for _, s := range stats {
				l, err := getUserHourlyStats(u, s, last)
				if err != nil {
					log.Println("Error in getUserHourlyStats:", err.Error())
				}
				rval[u][s] = l
			}
		}
		json.NewEncoder(w).Encode(rval)
	case "csv":
		enc := csv.NewWriter(w)
		defer enc.Flush()
		enc.Write([]string{"user_id", "stat_id", "hour", "value"})
		for _, u := range users {
			for _, s := range stats {
				l, err := getUserHourlyStats(u, s, last)
				if err != nil {
					log.Println("Error in getUserHourlyStats:", err.Error())
				}
				for k, v := range l {
					if err := enc.Write([]string{u, s, k, strconv.Itoa(v)}); err != nil {
						log.Println("error writing csv:", err.Error())
					}
				}
			}
		}
	}
}
