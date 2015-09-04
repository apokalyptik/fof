package main

import (
	"fmt"
)

func getXboxProfileURL(profileName string) (l string) {
	return fmt.Sprintf("<https://account.xbox.com/en-us/profile?gamerTag=%s|%s's Xbox Profile'>", profileName, profileName)
}
