package main

import (
	"fmt"
	"strings"
)

func getXboxProfileURL(profileName string) (l string) {
	gamerTag := strings.Replace(profileName, " ", "%20", -1)
	return fmt.Sprintf("<https://account.xbox.com/en-us/profile?gamerTag=%s|%s's Xbox Profile'>", gamerTag, profileName)
}
