package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/apokalyptik/gopid"
	"github.com/olebedev/config"
)

var configFile string
var admins []string

func init() {
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to YAML configuration")
}

func main() {
	flag.Parse()
	cfg, err := config.ParseYamlFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if pidFile, err := cfg.String("pidfile"); err != nil {
		log.Fatal(err)
	} else {
		if _, err := pid.Do(pidFile); err != nil {
			log.Fatalf("error creating pidfile (%s): %s", pidFile, err.Error())
			log.Fatal(err)
		}
	}

	if needsDbFile, err := cfg.String("database.needs"); err != nil {
		log.Fatal(err)
	} else {
		if err := needsDB.load(needsDbFile); err != nil {
			log.Fatal(err)
		}
		go needsDB.mindExpiration()
	}

	if raidsDbFile, err := cfg.String("database.raids"); err != nil {
		log.Fatal(err)
	} else {
		if err := raidDb.load(raidsDbFile); err != nil {
			log.Fatal(err)
		}
		if dur, err := cfg.String("maxAge"); err != nil {
			log.Fatal(err)
		} else {
			if maxAge, err := time.ParseDuration(dur); err != nil {
				log.Fatal(err)
			} else {
				go raidDb.mindExpiration(maxAge)
			}
		}

	}

	if slack.key, err = cfg.String("slack.slashKey"); err != nil {
		log.Fatal(err)
	}
	if slack.name, err = cfg.String("slack.webhooks.name"); err != nil {
		log.Fatal(err)
	}
	if slack.url, err = cfg.String("slack.webhooks.url"); err != nil {
		log.Fatal(err)
	}
	if slack.emoji, err = cfg.String("slack.webhooks.emoji"); err != nil {
		log.Fatal(err)
	}

	if adminsvar, err := cfg.List("slack.admins"); err != nil {
		log.Fatal(err)
	} else {
		for _, v := range adminsvar {
			if admin, ok := v.(string); ok {
				admins = append(admins, admin)
			} else {
				log.Fatalf("%#v from %#v is not a string", v, adminsvar)
			}
		}
	}

	if listen, err := cfg.String("listen"); err != nil {
		log.Fatal(err)
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { doHTTP404(w) })
		http.HandleFunc("/api", doHTTPRouter)
		log.Fatal(http.ListenAndServe(listen, nil))
	}
}
