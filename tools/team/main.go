package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/apokalyptik/gopid"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/sessions"
	"github.com/olebedev/config"
)

var configFile string
var admins []string
var cleanFromFilename = regexp.MustCompile(`[^0-9a-zA-Z _-]$`)

func init() {
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to YAML configuration")
	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)
}

func main() {
	flag.Parse()

	cfg, err := config.ParseYamlFile(configFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err.Error())
	}

	if pidFile, err := cfg.String("pidfile"); err != nil {
		log.Fatalf("Error reading pidfile from config: %s", err.Error())
	} else {
		if _, err := pid.Do(pidFile); err != nil {
			log.Fatalf("error creating pidfile (%s): %s", pidFile, err.Error())
		}
	}

	if raidsDbFile, err := cfg.String("database.raids"); err != nil {
		log.Fatalf("Error reading database.raids from config file: %s", err.Error())
	} else {
		if err := raidDb.load(raidsDbFile); err != nil {
			log.Fatalf("Error reading %s: %s", raidsDbFile, err.Error())
		}
		if dur, err := cfg.String("maxAge"); err != nil {
			log.Fatalf("Error reading maxAge from config file: %s", err.Error())
		} else {
			if maxAge, err := time.ParseDuration(dur); err != nil {
				log.Fatalf("Error parsing maxAge as a time.Duration: %s", err.Error())
			} else {
				go raidDb.mindExpiration(maxAge)
			}
		}
	}

	if userDbFile, err := cfg.String("database.users"); err != nil {
		log.Fatalf("Error reading database.users from config file: %s", err.Error())
	} else {
		if err := udb.load(userDbFile); err != nil {
			log.Fatalf("Error reading %s: %s", userDbFile, err.Error())
		}
	}

	if lfgDbFile, err := cfg.String("database.lfg"); err != nil {
		log.Fatalf("Error reading database.lfg from config file: %s", err.Error())
	} else {
		if err := lfg.load(lfgDbFile); err != nil {
			log.Fatalf("Error reading %s: %s", lfgDbFile, err.Error())
		}
		go lfg.mindExpiration()
	}

	if slack.raidKey, err = cfg.String("slack.slashKey.raids"); err != nil {
		log.Fatalf("Error reading slack.slashKey.raids: %s", err.Error())
	}

	if adminsvar, err := cfg.List("slack.admins"); err != nil {
		log.Fatalf("Error reading slack.admins: %s", err.Error())
	} else {
		for _, v := range adminsvar {
			if admin, ok := v.(string); ok {
				admins = append(admins, admin)
			} else {
				log.Fatalf("%#v from %#v is not a string", v, adminsvar)
			}
		}
		xhrOutput.set("admins", admins)
	}

	if apiKey, err := cfg.String("slack.apiKey"); err != nil {
		log.Fatalf("Error reading slack.apiKey from config: %s", err.Error())
	} else {
		slack.apiKey = apiKey
	}

	if hmac, err := cfg.String("security.hmac_key"); err != nil {
		rand.Read(hmacKey)
	} else {
		hmacKey = []byte(hmac)
	}

	if auth, err := cfg.String("security.cookie.auth"); err == nil {
		if key, err := cfg.String("security.cookie.key"); err == nil {
			store = sessions.NewCookieStore([]byte(auth), []byte(key))
		} else {
			log.Fatalf("error reading security.cookie.key from config: %s", err.Error())
		}
	} else {
		log.Fatalf("error reading security.cookie.auth from config: %s", err.Error())
	}

	if cmd, err := cfg.String("slack.slashCommand.raids"); err == nil {
		raidSlashCommand = cmd
		xhrOutput.set("command", cmd)
	}

	go mindSlackMsgQueue()
	if listen, err := cfg.String("listen"); err != nil {
		log.Fatal(err)
	} else {
		mux := http.NewServeMux()
		var devmode = false
		if _, err := os.Stat("www/index.html"); err == nil {
			devmode = true
		}
		if devmode == false {
			mux.Handle("/",
				http.FileServer(
					&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "www"}))
		} else {
			mux.Handle("/", http.FileServer(http.Dir("www/")))
		}
		mux.HandleFunc("/api", doHTTPPost)
		mux.HandleFunc("/rest/", doRESTRouter)
		mux.HandleFunc("/ics", func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "raidbot")
			session.Options.MaxAge = 604800
			uri, _ := url.ParseRequestURI(r.RequestURI)
			if v, err := url.ParseQuery(uri.RawQuery); err != nil {
				log.Println(err.Error())
				doHTTPStatus(w, http.StatusInternalServerError)
				return
			} else {
				if err := requireAPIKey(session, w); err != nil {
					log.Println(err.Error())
					doHTTPStatus(w, http.StatusInternalServerError)
					return
				}
				if name, ok := session.Values["username"]; ok {
					w.Header().Set("Content-Type", "text/Calendar")
					w.Header().Set("Content-Disposition", "inline; filename="+cleanFromFilename.ReplaceAllString(v.Get("title"), "-")+".ics")
					fmt.Fprintf(w, v.Get("data"))
					log.Printf("@%s -- ics: %s", name.(string), v.Get("title"))
				} else {
					doHTTPStatus(w, http.StatusForbidden)
				}
			}
		})

		go mindPvtGroupList()
		go mindChannelList()

		log.Println("Starting Up!")
		log.Fatal(http.ListenAndServe(listen, ch.Handler(mux)))
	}
}
