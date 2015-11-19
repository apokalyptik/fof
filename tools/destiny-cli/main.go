package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/apokalyptik/fof/lib/destiny"
)

var apiKey = ""

func init() {
	flag.StringVar(&apiKey, "apikey", apiKey, "Your API key (see: https://www.bungie.net/en/User/API)")
}

func main() {
	flag.Parse()
	if apiKey == "" {
		fmt.Println("Please specify a valid API key")
		os.Exit(1)
	}
	client := destiny.New(apiKey, "github.com/apokalyptik/fof/tools/destiny-cli")

	var rsp struct {
		Data struct {
			Characters []struct {
				CharacterBase struct {
					CurrentActivityHash int64     `json:"currentActivityHash"`
					DateLastPlayed      time.Time `json:"dateLastPlayed"`
				} `json:"characterBase"`
			} `json:"characters"`
		} `json:"data"`
	}
	if err := client.AccountSummary(destiny.PlatformXBL, "4611686018437911483", &rsp); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(rsp)
}
