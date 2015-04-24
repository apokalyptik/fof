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
		log.Fatalf("Error parsing config file: %s", err.Error())
	}

	if pidFile, err := cfg.String("pidfile"); err != nil {
		log.Fatalf("Error reading pidfile from config: %s", err.Error())
	} else {
		if _, err := pid.Do(pidFile); err != nil {
			log.Fatalf("error creating pidfile (%s): %s", pidFile, err.Error())
		}
	}

	if xlineDbFile, err := cfg.String("database.xline"); err != nil {
		log.Fatal(err)
	} else {
		if err := xlineDB.load(xlineDbFile); err != nil {
			log.Fatal(err)
		}
	}

	if needsDbFile, err := cfg.String("database.needs"); err != nil {
		log.Fatalf("Error reading database.needs from config file: %s", err.Error())
	} else {
		if err := needsDB.load(needsDbFile); err != nil {
			log.Fatalf("Error reading %s: %s", needsDbFile, err.Error())
		}
		go needsDB.mindExpiration()
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
	if slack.xlineKey, err = cfg.String("slack.slashKey.xline"); err != nil {
		log.Fatalf("Error reading slack.slashKey.xline: %s", err.Error())
	}
	if slack.needKey, err = cfg.String("slack.slashKey.needs"); err != nil {
		log.Fatalf("Error reading slack.slashKey.needs: %s", err.Error())
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

	if listen, err := cfg.String("listen"); err != nil {
		log.Fatal(err)
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { doHTTP404(w) })
		http.HandleFunc("/api", doHTTPRouter)
		log.Fatal(http.ListenAndServe(listen, nil))
	}
}
