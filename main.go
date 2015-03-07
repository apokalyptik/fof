package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/apokalyptik/gopid"
	"github.com/olebedev/config"
)

var configFile string
var dbFile string
var db = &raids{}

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

	if dbFile, err := cfg.String("database"); err != nil {
		log.Fatal(err)
	} else {
		if err := db.load(dbFile); err != nil {
			log.Fatal(err)
		}
	}

	db.save()

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

	if listen, err := cfg.String("listen"); err != nil {
		log.Fatal(err)
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { HTTP404(w) })
		http.HandleFunc("/api", HTTPRouter)
		log.Fatal(http.ListenAndServe(listen, nil))
	}
}
