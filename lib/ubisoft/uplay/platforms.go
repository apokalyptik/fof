package uplay

import "strings"

const (
	PlatformUnknown = 0
	PlatformXBL     = 1
	PlatformPSN     = 2
	PlatformPC      = 3
)

var uplayoverlayPlatforms = map[int]string{
	PlatformXBL: "xbl",
	PlatformPSN: "psn",
	PlatformPC:  "uplay",
}

var uplaywebcenterPlatforms = map[int]string{
	PlatformXBL: "XONE",
	PlatformPSN: "PSN",
	PlatformPC:  "PC",
}

func GuessPlatform(platform string) int {
	platform = strings.ToLower(platform)
	switch string(platform[0]) {
	case "x":
		return PlatformXBL
	case "p":
		switch string(platform[1]) {
		case "s":
			return PlatformPSN
		case "c":
			return PlatformPC
		}
	}
	return PlatformUnknown
}
