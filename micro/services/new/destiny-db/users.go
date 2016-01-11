package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type metaValue string

func (m metaValue) String() string {
	return string(m)
}

func (m metaValue) Time() (time.Time, error) {
	var v = m.String()
	if v == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", v)
}

type user struct {
	ID             int
	UserID         string
	UserName       string
	GamerTag       string
	DestinyID      string
	Seen           time.Time
	DestinyChecked int
}

func (u *user) getMeta(key string) (metaValue, error) {
	var value string
	res, err := queryx("SELECT Value FROM userMeta Where UserID=? AND Key=? LIMIT 1", u.ID, key)
	if res.Next() {
		if err == nil {
			err = res.Scan(&value)
			if err == sql.ErrNoRows {
				err = nil
			}
		}
	}
	return metaValue(value), err
}

func (u *user) getMetaTime(key string) (time.Time, error) {
	mv, err := u.getMeta(key)
	if err != nil {
		return time.Time{}, err
	}
	return mv.Time()
}

func (u *user) setMeta(key string, value string) error {
	_, err := exec("INSERT OR REPLACE INTO userMeta (UserID,Key,Value) VALUES(?,?,?)", u.ID, key, value)
	return err
}

func getUserByID(ID int) *user {
	var u = user{}
	if err := q["getUserByID"].Get(&u, ID); err != nil {
		log.Printf("error getting by userID: %s", err.Error())
	}
	return &u
}

func getUserByUserID(userID string) *user {
	var u = user{}
	if err := q["getUserByUserID"].Get(&u, userID); err != nil {
		log.Printf("error getting by userID: %s", err.Error())
	}
	return &u
}

var lastSeen = map[string]time.Time{}

type bungieSearchResponse []struct {
	MembershipId string `json:"membershipId"`
}

type slackUserListResponse struct {
	OK      bool   `json:"ok"`
	Error   string `json:"error"`
	Members []struct {
		ID              int    `json:"ID"`
		UserID          string `json:"id"`
		UserName        string `json:"name"`
		Bot             bool   `json:"is_bot"`
		Deleted         bool   `json:"deleted"`
		Restricted      bool   `json:"is_restricted"`
		UltraRestricted bool   `json:"is_ultra_restricted"`
		Profile         struct {
			GamerTag string `json:"first_name"`
		} `json:"profile"`
	} `json:"members"`
}

func init() {
	qs["deleteUser"] = "DELETE FROM users WHERE UserID = ?"
	qs["createUser"] = "INSERT OR IGNORE INTO users (UserID,UserName) VALUES(?,?)"
	qs["updateUserSeen"] = "UPDATE OR IGNORE users SET Seen=datetime(?, 'utc') WHERE userID = ?"
	qs["updateUserName"] = "UPDATE OR IGNORE users SET UserName=? WHERE userID = ?"
	qs["updateGamerTag"] = "UPDATE OR IGNORE users SET GamerTag=? WHERE userID = ?"
	qs["updateDestinyID"] = "UPDATE OR IGNORE users SET DestinyID=? WHERE userID = ?"
	qs["updateDestinyChecked"] = "UPDATE OR IGNORE users SET DestinyChecked=? WHERE userID = ?"
	qs["getDestinyID"] = "SELECT destinyID,destinyChecked FROM users WHERE UserID = ? LIMIT 1"
	qs["getUserByUserID"] = "SELECT * FROM users WHERE UserID = ? LIMIT 1"
	qs["getUserByID"] = "SELECT * FROM users WHERE ID = ? LIMIT 1"
}

func updateSeenList() {
	rsp, err := http.Get(creds["seenURL"])
	if err != nil {
		log.Println("error fetching seen list:", err.Error())
		return
	}
	defer rsp.Body.Close()
	var newLastSeen = map[string]time.Time{}
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&newLastSeen); err != nil {
		log.Println("error decoding seen data:", err.Error())
		return
	}
	lastSeen = newLastSeen
}

func updateUserList() {
	var url = fmt.Sprintf("https://slack.com/api/users.list?token=%s", creds["slackAdminToken"])
	rsp, err := http.Get(url)
	if err != nil {
		log.Println("error fetching slack user list:", err.Error())
		return
	}
	defer rsp.Body.Close()
	var slackData slackUserListResponse
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&slackData); err != nil {
		log.Println("error decoding slack user list response:", err.Error())
		return
	}
	if !slackData.OK {
		log.Println("error in slack user list response:", slackData.Error)
		return
	}
	for _, v := range slackData.Members {
		var deleted = false
		if v.Bot {
			deleted = true
		} else if v.Deleted {
			deleted = true
		} else if v.Restricted {
			deleted = true
		} else if v.UltraRestricted {
			deleted = true
		}
		if deleted {
			if res, err := q["deleteUser"].Exec(v.UserID); err != nil {
				log.Printf("error deleting user %s (%s): %s", v.UserName, v.UserID, err.Error())
			} else {
				if n, _ := res.RowsAffected(); n > 0 {
					pubUser("delete", v.UserID)
				}
			}
			continue
		}
		var created = false
		if res, err := q["createUser"].Exec(v.UserID, v.UserName); err != nil {
			log.Printf("error creating user %s (%s): %s", v.UserName, v.UserID, err.Error())
		} else {
			if n, _ := res.RowsAffected(); n > 0 {
				created = true
			}
		}

		localUser := getUserByUserID(v.UserID)

		// TODO: Move to user struct
		var updated = false
		if seen, ok := lastSeen[v.UserID]; ok {
			if !localUser.Seen.Equal(seen) {
				q["updateUserSeen"].Exec(seen, v.UserID)
				localUser.Seen = seen
				updated = true
			}
		}

		if localUser.UserName != v.UserName {
			q["updateUserName"].Exec(v.UserName, v.UserID)
			localUser.UserName = v.UserName
			updated = true
		}

		if v.Profile.GamerTag != localUser.GamerTag {
			q["updateGamerTag"].Exec(v.Profile.GamerTag, v.UserID)
			q["updateDestinyChecked"].Exec(0, v.UserID)
			localUser.DestinyChecked = 0
			localUser.GamerTag = v.Profile.GamerTag
			updated = true
		}

		if localUser.DestinyID == "" && localUser.DestinyChecked == 0 {
			if req, err := destinyClient.SearchDestinyPlayer(v.Profile.GamerTag); err != nil {
				log.Printf("error generating SearchDestinyPlayer request for %s (%s): %s", v.UserID, v.UserName, err.Error())
			} else {
				var d bungieSearchResponse
				if err := req.Into(&d); err != nil {
					log.Printf("error fetching destiny id for %s (%s): %s", v.UserID, v.UserName, err.Error())
				} else {
					// TODO: move to user struct
					q["updateDestinyChecked"].Exec(1, v.UserID)
					if len(d) > 0 {
						log.Println(v.UserID, v.UserName, d[0].MembershipId)
						q["updateDestinyID"].Exec(d[0].MembershipId, v.UserID)
						updated = true
					} else {
						log.Printf("unable to find %s -- %s -- '%s'", v.UserID, v.UserName, v.Profile.GamerTag)
					}
				}
			}
		}
		if created {
			pubUser("new", v.UserID)
		} else if updated {
			pubUser("update", v.UserID)
		}
	}
}

func mindUserList() {
	for {
		log.Println("mindUserList: Updating last seen list")
		updateSeenList()
		log.Println("mindUserList: Updating user list")
		updateUserList()
		log.Println("mindUserList: Sleeping")
		time.Sleep(5 * time.Minute)
	}
}
