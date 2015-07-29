package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	platformNone = iota // 0
	platformXBL         // 1
	platformPSN         // 2
)

type destinyClientResponse interface{}

type destinyClient struct {
	apiKey string
}

func (d *destinyClient) get(URL string) (interface{}, error) {
	var rval struct {
		ErrorCode       int
		ThrottleSeconds int
		ErrorStatus     string
		Message         string
		MessageData     interface{}
		Response        interface{}
	}
	if req, err := http.NewRequest("GET", URL, nil); err != nil {
		return nil, err
	} else {
		var c = &http.Client{}
		req.Header.Add("X-API-Key", d.apiKey)
		req.Header.Add("User-Agent", "FederationOfFathers-StatsBot/0.1 (http://federationoffathers.com/)")
		if resp, err := c.Do(req); err != nil {
			return nil, err
		} else {
			defer resp.Body.Close()
			decoder := json.NewDecoder(resp.Body)
			if err := decoder.Decode(&rval); err != nil {
				return nil, err
			} else {
				if rval.ThrottleSeconds != 0 {
					// Super Naive
					log.Println("Asked to throttle for", rval.ThrottleSeconds, "seconds")
					time.Sleep(time.Duration(rval.ThrottleSeconds) * time.Second)
				}
				if rval.Message != "Ok" {
					return nil, fmt.Errorf("%s (%d) -- %s", rval.Message, rval.ErrorCode, rval.ErrorStatus)
				}
				return rval.Response, nil
			}
		}
	}
}

func platformURL(suffix string) string {
	var url = fmt.Sprintf("http://www.bungie.net/Platform/Destiny/%s?definitions=true", suffix)
	// var url = fmt.Sprintf("http://www.bungie.net/Platform/Destiny/%s", suffix)
	log.Println(url)
	return url
}

func playerURL(platform int, username string) string {
	return platformURL(fmt.Sprintf("SearchDestinyPlayer/%d/%s/", platform, username))
}

func grimoireURL(platform int, memberID string) string {
	return platformURL(fmt.Sprintf("Vanguard/Grimoire/%d/%s/", platform, memberID))
}

func accountStatsURL(platform int, memberID string) string {
	return platformURL(fmt.Sprintf("Stats/Account/%d/%s/", platform, memberID))
}

func triumphsURL(platform int, memberID string) string {
	return platformURL(fmt.Sprintf("%d/Account/%s/Triumphs/", platform, memberID))
}

func accountURL(platform int, memberID string) string {
	return platformURL(fmt.Sprintf("%d/Account/%s/", platform, memberID))
}

func charURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("%d/Account/%s/Character/%s/", platform, memberID, charID))
}

func charInventoryURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("%d/Account/%s/Character/%s/Inventory/", platform, memberID, charID))
}

func charInventoryItemURL(platform int, memberID, charID, itemID string) string {
	return platformURL(fmt.Sprintf("%d/Account/%s/Character/%s/Inventory/%s", platform, memberID, charID, itemID))
}

func charActivityHistoryURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("Stats/ActivityHistory/%d/%s/%s/", platform, memberID, charID))
}

func charActivitiesURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("Stats/ActivityHistory/%d/%s/%s/Activities/", platform, memberID, charID))
}

func charProgressionURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("Stats/ActivityHistory/%d/%s/%s/Progression/", platform, memberID, charID))
}

func charActivityStatsURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("Stats/AggregateActivityStats/%d/%s/%s/", platform, memberID, charID))
}

func charUniqueWeaponsStateURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("Stats/UniqueWeapons/%d/%s/%s/", platform, memberID, charID))
}

func charStatsURL(platform int, memberID, charID string) string {
	return platformURL(fmt.Sprintf("Stats/Stats/%d/%s/%s/", platform, memberID, charID))
}
