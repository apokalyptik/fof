package main

import "github.com/apokalyptik/cfg"

var bungieApiKey string

func main() {
	var c *cfg.Options = cfg.New("bungie")
	c.StringVar(&bungieApiKey, "api", "", "see: https://www.bungie.net/en-US/User/API")
	cfg.Parse()
}
