package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/apokalyptik/fof/lib/ubisoft/uplay"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/koding/multiconfig"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
)

type singleUser struct {
	ID       string `json:"id"`
	GamerTag string `json:"gamertag"`
	UserName string `json:"username"`
}

type users []singleUser

// Config holds uplay configuration
type Config struct {
	Username string
	Password string
	Listen   string
}

type userStats struct {
	User  singleUser
	Stats []uplay.DivisionStat
}

var (
	config = &Config{
		Listen:   "0.0.0.0:8885",
		Username: "UplayUsername",
		Password: "UplayPassword",
	}
	failures          = map[string]time.Time{}
	userListAddress   = "http://fofgaming.com:8880/fof/members.json"
	uplayUserID       = map[string]string{}
	skipNoUUID        = 24 * time.Hour
	lastUserlistFetch time.Time
	lastUserlist      users
	stats             = map[string]userStats{}
	statsTime         = map[string]time.Time{}
	statsLock         sync.Mutex
)

func userlist() (users, error) {
	if time.Now().Sub(lastUserlistFetch) < time.Hour {
		if lastUserlist != nil {
			return lastUserlist, nil
		}
	}
	var list users
	rsp, err := http.Get(userListAddress)
	if err != nil {
		if lastUserlist != nil {
			if time.Now().Sub(lastUserlistFetch) < 24*time.Hour {
				return lastUserlist, nil
			} else {
				log.Fatalf("Error fetching user list: %s", err.Error())
			}
		}
		return nil, fmt.Errorf("Error fetching user list: %s", err.Error())
	}
	defer rsp.Body.Close()
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&list); err != nil {
		if lastUserlist != nil {
			if time.Now().Sub(lastUserlistFetch) < 24*time.Hour {
				return lastUserlist, nil
			} else {
				log.Fatalf("Error fetching user list: %s", err.Error())
			}
		}
		return nil, fmt.Errorf("Error decoding user list: %s", err.Error())
	}
	return list, nil
}

func save() {
	statsLock.Lock()
	defer statsLock.Unlock()
	var state = struct {
		Stats             map[string]userStats
		StatsTime         map[string]time.Time
		Failures          map[string]time.Time
		UplayUserID       map[string]string
		LastUserlistFetch time.Time
		LastUserlist      users
	}{
		stats,
		statsTime,
		failures,
		uplayUserID,
		lastUserlistFetch,
		lastUserlist,
	}
	if buf, err := json.Marshal(state); err != nil {
		log.Println("[init] Error marshalling .state.json:", err.Error())
	} else {
		if err := ioutil.WriteFile(".state.json", buf, 0644); err != nil {
			log.Println("[init] Error writing to .state.json:", err.Error())
		} else {
			log.Println("[init] Saved state to .state.json")
		}
	}
}

func load() {
	var state = struct {
		Stats             map[string]userStats
		StatsTime         map[string]time.Time
		Failures          map[string]time.Time
		UplayUserID       map[string]string
		LastUserlistFetch time.Time
		LastUserlist      users
	}{}
	if _, err := os.Stat(".state.json"); err == nil {
		if file, err := os.Open(".state.json"); err != nil {
			log.Println("[init] Error reading .state.json:", err.Error())
			return
		} else {
			defer file.Close()
			dec := json.NewDecoder(file)
			if err := dec.Decode(&state); err != nil {
				log.Println("[init] Error decoding .state.json:", err.Error())
			} else {
				if state.Stats != nil {
					stats = state.Stats
				}
				if state.StatsTime != nil {
					statsTime = state.StatsTime
				}
				if state.Failures != nil {
					failures = state.Failures
				}
				if state.UplayUserID != nil {
					uplayUserID = state.UplayUserID
				}
				lastUserlistFetch = state.LastUserlistFetch
				lastUserlist = state.LastUserlist
				log.Println("[init] Loaded state from .state.json")
			}
		}
	}
}

func init() {
	var m multiconfig.Loader
	var l multiconfig.Loader
	var found = true
	if _, err := os.Stat(".uplay.toml"); err == nil {
		l = &multiconfig.TOMLLoader{Path: ".uplay.toml"}
	} else if _, err := os.Stat(".uplay.json"); err == nil {
		l = &multiconfig.JSONLoader{Path: ".uplay.json"}
	} else if _, err := os.Stat(fmt.Sprintf("%s/.uplay.toml", os.Getenv("HOME"))); err == nil {
		l = &multiconfig.TOMLLoader{Path: fmt.Sprintf("%s/.uplay.toml", os.Getenv("HOME"))}
	} else if _, err := os.Stat(fmt.Sprintf("%s/.uplay.json", os.Getenv("HOME"))); err == nil {
		l = &multiconfig.JSONLoader{Path: fmt.Sprintf("%s/.uplay.json", os.Getenv("HOME"))}
	} else {
		found = false
	}
	if found {
		m = multiconfig.MultiLoader(l, &multiconfig.FlagLoader{}, &multiconfig.EnvironmentLoader{})
	} else {
		m = multiconfig.MultiLoader(&multiconfig.FlagLoader{}, &multiconfig.EnvironmentLoader{})
	}
	m.Load(config)
}

func mindHTTP() {
	r := mux.NewRouter()
	r.HandleFunc("/v1.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")
		statsLock.Lock()
		defer statsLock.Unlock()
		enc := json.NewEncoder(w)
		enc.Encode(stats)
	})
	r.HandleFunc("/v1/gt/{gamertag}.json", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		enc := json.NewEncoder(w)
		statsLock.Lock()
		defer statsLock.Unlock()
		if v, ok := stats[vars["gamertag"]]; ok {
			w.Header().Set("Content-Type", "text/json")
			enc.Encode(v)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	r.HandleFunc("/v1/slackid/{slackuid}.json", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		enc := json.NewEncoder(w)
		statsLock.Lock()
		defer statsLock.Unlock()
		for _, v := range stats {
			if v.User.ID == vars["slackuid"] {
				w.Header().Set("Content-Type", "text/json")
				enc.Encode(v)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})
	r.HandleFunc("/v1/username/{slackusername}.json", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		enc := json.NewEncoder(w)
		statsLock.Lock()
		defer statsLock.Unlock()
		for _, v := range stats {
			if v.User.UserName == vars["slackusername"] {
				w.Header().Set("Content-Type", "text/json")
				enc.Encode(v)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})
	r.HandleFunc("/v1/compact.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		var out = struct {
			Index map[string]interface{}
			Stats []struct {
				User  singleUser
				Stats map[string]string
			}
		}{
			Stats: []struct {
				User  singleUser
				Stats map[string]string
			}{},
		}
		statsLock.Lock()
		for _, v := range stats {
			var entry = struct {
				User  singleUser
				Stats map[string]string
			}{
				User:  v.User,
				Stats: map[string]string{},
			}
			if out.Index == nil {
				out.Index = map[string]interface{}{}
				for _, s := range v.Stats {
					out.Index[s.Name] = s
				}
			}
			for _, s := range v.Stats {
				entry.Stats[s.Name] = s.Value.String()
			}
			out.Stats = append(out.Stats, entry)
		}
		statsLock.Unlock()
		w.Header().Set("Content-Type", "text/json")
		enc.Encode(out)
	})
	n := negroni.New()
	n.Use(cors.New(cors.Options{AllowedOrigins: []string{"*"}}))
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(r)
	n.Run(config.Listen)
}

func main() {
	load()
	go mindHTTP()
	client := uplay.New(config.Username, config.Password)
	t := time.Tick(time.Hour)
	var lastAuth = time.Now().Add(0 - 24*time.Hour)
	for {
		var setupError error
		var list users
		if setupError == nil {
			if time.Now().Sub(lastAuth) > 8*time.Hour {
				if err := client.Authenticate(); err != nil {
					setupError = fmt.Errorf("[uplay] Error authenticating to ubisoft: %s", err.Error())
				} else {
					log.Println("[uplay] Authenticated with ubisoft")
					lastAuth = time.Now()
				}
			} else {
				log.Printf("[uplay] Reusing existing uplay credentials")
			}
		}
		if setupError == nil {
			if l, err := userlist(); err != nil {
				setupError = fmt.Errorf("[users] Error fetching user list: %s", err.Error())
			} else {
				list = l
			}
		}
		if setupError != nil {
			log.Println(setupError.Error())
			<-t
			continue
		}
		for _, user := range list {
			if t, ok := statsTime[user.ID]; ok {
				if diff := time.Now().Sub(t); diff < time.Hour {
					log.Println("[users] Skipping", user.GamerTag, "for recent fetch time of", diff.String())
					continue
				}
			}
			t, ok := failures[user.ID]
			if ok && t.Sub(time.Now()) < skipNoUUID {
				log.Println("[users] Skipping", user.GamerTag, "known not to have a uplay ID")
				continue
			}
			log.Println("[users] Processing", user.GamerTag)
			uuid, ok := uplayUserID[user.ID]
			if !ok {
				if profiles, err := client.UserSearch(uplay.PlatformXBL, user.GamerTag); err != nil {
					log.Println("[uplay] aborting run: error searching for uuid:", err.Error())
					lastAuth = time.Now().Add(0 - 24*time.Hour)
					break
				} else {
					if len(profiles) > 0 {
						if profiles[0].UplayID == "" {
							log.Println("[users] No uuid found for", user.GamerTag, "skipping for", skipNoUUID.String())
							failures[user.ID] = time.Now()
							continue
						}
						uuid = profiles[0].UplayID
						uplayUserID[user.ID] = uuid
						log.Println("[users] Found UUID ", uuid, "for", user.GamerTag)
					} else {
						log.Println("[users] No uuid found for", user.GamerTag, "skipping for", skipNoUUID.String())
						failures[user.ID] = time.Now()
						continue
					}
				}
			} else {
				if uuid == "" {
					delete(uplayUserID, user.ID)
					log.Println("[users] No uuid found for", user.GamerTag, "skipping for", skipNoUUID.String())
					failures[user.ID] = time.Now()
					continue
				}
			}
			if st, err := client.DivisionStats(uplay.PlatformXBL, uuid); err != nil {
				log.Println("[uplay] aborting run: pulling division stats for", user.GamerTag, "/", uuid, "error:", err.Error())
				lastAuth = time.Now().Add(0 - 24*time.Hour)
				break
			} else {
				statsLock.Lock()
				statsTime[user.ID] = time.Now()
				stats[user.GamerTag] = userStats{
					User:  user,
					Stats: st,
				}
				statsLock.Unlock()
				log.Println("[users] updated division stats for", user.GamerTag)
			}
		}
		save()
		<-t
	}
}
