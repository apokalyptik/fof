package main

import (
	"crypto/rand"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/apokalyptik/gopid"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/sessions"
	"github.com/olebedev/config"
)

var configFile string
var admins []string

func init() {
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to YAML configuration")
	rand.Read(hmacKey)
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

	if slack.raidKey, err = cfg.String("slack.slashKey.raids"); err != nil {
		log.Fatalf("Error reading slack.slashKey.raids: %s", err.Error())
	}
	if slack.name, err = cfg.String("slack.webhooks.name"); err != nil {
		log.Fatalf("Error reading slack.webhooks.name: %s", err.Error())
	}
	if slack.url, err = cfg.String("slack.webhooks.url"); err != nil {
		log.Fatalf("Error reading slack.webhooks.url: %s", err.Error())
	}
	if slack.emoji, err = cfg.String("slack.webhooks.emoji"); err != nil {
		log.Fatalf("Error reading slack.webhooks.emoji: %s", err.Error())
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
	}

	if apiKey, err := cfg.String("slack.apiKey"); err != nil {
		log.Fatalf("Error reading slack.apiKey from config: %s", err.Error())
	} else {
		slack.apiKey = apiKey
	}

	if auth, err := cfg.String("cookie.auth"); err == nil {
		if key, err := cfg.String("cookie.key"); err == nil {
			store = sessions.NewCookieStore([]byte(auth), []byte(key))
		} else {
			log.Fatalf("error reading cookie.key from config: %s", err.Error())
		}
	} else {
		log.Fatalf("error reading cookie.auth from config: %s", err.Error())
	}

	if cmd, err := cfg.String("slack.slashCommand.raids"); err == nil {
		raidSlashCommand = cmd
	}
	go mindSlackMsgQueue()
	if listen, err := cfg.String("listen"); err != nil {
		log.Fatal(err)
	} else {
		var devmode = false
		if _, err := os.Stat("www/index.html"); err == nil {
			devmode = true
		}
		if devmode == false {
			http.Handle("/",
				http.FileServer(
					&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "www"}))
		} else {
			http.Handle("/", http.FileServer(http.Dir("www/")))
		}
		http.HandleFunc("/api", doHTTPPost)
		http.HandleFunc("/rest/", doRESTRouter)

		go mindChannelList()

		log.Fatal(http.ListenAndServe(listen, nil))
	}
}
