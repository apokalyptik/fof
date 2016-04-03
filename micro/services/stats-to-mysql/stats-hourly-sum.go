package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type usersHourlySumStats map[string]userHourlySumStatsList
type userHourlySumStatsList map[string]int

func getUserHourlySumStats(user, stat string, last int) (int, error) {
	var rval *int
	var err error
	var row *sql.Row
	var stat_id int
	stat_id, err = strconv.Atoi(stat)
	if err != nil {
		return *rval, err
	}
	if last >= 0 {
		row = db.QueryRow(
			"SELECT SUM(`value`) FROM stats_hourly WHERE `member`=? AND `stat_id`=? AND `when` >= DATE_SUB(NOW(), INTERVAL ? HOUR)",
			user,
			stat_id,
			last)
	} else {
		row = db.QueryRow(
			"SELECT SUM(`value`) FROM stats_hourly WHERE `member`=? AND `stat_id`=?",
			user,
			stat_id)
	}
	err = row.Scan(&rval)
	if rval != nil {
		return *rval, err
	}
	return 0, err
}

func handleHourlySumJson(w http.ResponseWriter, r *http.Request) {
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
	var rval = usersHourlySumStats{}
	for _, u := range users {
		rval[u] = userHourlySumStatsList{}
		for _, s := range stats {
			l, err := getUserHourlySumStats(u, s, last)
			if err != nil {
				log.Println("Error in getUserHourlySumStats:", err.Error())
			}
			rval[u][s] = l
		}
	}
	json.NewEncoder(w).Encode(rval)
}
