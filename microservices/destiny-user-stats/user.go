package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
)

var errUserNotFound = errors.New("User Not Found")

type userBit struct {
	when time.Time
	data interface{}
}

type user struct {
	name     string
	platform int
	account  string
	data     map[string]*userBit
	document []byte
}

func (u *user) set(what string, info interface{}, err error) {
	if _, ok := u.data[what]; !ok {
		u.data[what] = &userBit{}
	}
	if err != nil {
		log.Printf("Error fetching %s for %s: %s", what, u.name, err.Error())
		return
	}
	u.data[what].data = info
	u.data[what].when = time.Now()
}

func (u *user) pull() error {
	defer func() {
		var doc = map[string]interface{}{}
		for k, v := range u.data {
			doc[k] = v.data
		}
		if document, err := json.Marshal(doc); err != nil {
			log.Printf("Error marshalling user data for %s: %s", u.name, err.Error())
		} else {
			u.document = document
			ioutil.WriteFile(fmt.Sprintf("data/%d.%s.json", u.platform, u.account), u.document, 0644)
		}
	}()

	if u.platform == 0 {
		var search = []struct {
			MembershipType int    `mapstructure:"membershipType"`
			MembershipId   string `mapstructure:"membershipId"`
		}{}
		// Try XBL First
		info, err := client.get(playerURL(platformXBL, u.name))
		if err != nil {
			return err
		} else {
			if err := mapstructure.Decode(info, &search); err != nil {
				return err
			}
		}
		if len(search) < 1 {
			// Try PSN Second
			info, err = client.get(playerURL(platformPSN, u.name))
			if err != nil {
				return err
			} else {
				if err := mapstructure.Decode(info, &search); err != nil {
					return err
				}
			}
		}
		if len(search) < 1 {
			return errUserNotFound
		}
		u.platform = search[0].MembershipType
		u.account = search[0].MembershipId
		u.data["player"] = &userBit{
			when: time.Now(),
			data: info.([]interface{})[0],
		}
	}

	grim, err := client.get(grimoireURL(u.platform, u.account))
	u.set("grimoire", grim, err)

	astats, err := client.get(accountStatsURL(u.platform, u.account))
	u.set("accountStats", astats, err)

	t, err := client.get(triumphsURL(u.platform, u.account))
	u.set("triumphs", t, err)

	account, err := client.get(accountURL(u.platform, u.account))
	u.set("account", account, err)

	if err != nil {
		// todo... decide how to handle error handing in this function...
		return nil
	}

	var acctResponse struct {
		Data struct {
			Characters []struct {
				CharacterBase struct {
					CharacterId string `mapstructure:"characterId"`
				} `mapstructure:"characterBase"`
			} `mapstructure:"characters"`
		} `mapstructure:"data"`
	}

	if err := mapstructure.Decode(account, &acctResponse); err != nil {
		return err
	}

	if _, ok := u.data["characters"]; !ok {
		u.data["characters"] = &userBit{
			data: map[string]map[string]interface{}{},
		}
	}
	characterDoc := u.data["characters"].data.(map[string]map[string]interface{})
	for _, c := range acctResponse.Data.Characters {
		if _, ok := characterDoc[c.CharacterBase.CharacterId]; !ok {
			characterDoc[c.CharacterBase.CharacterId] = map[string]interface{}{}
		}
		cdoc := characterDoc[c.CharacterBase.CharacterId]

		if ch, err := client.get(charActivityHistoryURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["activityHistory"] = ch
		}
		if ca, err := client.get(charActivitiesURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["activities"] = ca
		}
		if cp, err := client.get(charProgressionURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["progression"] = cp
		}
		if cas, err := client.get(charActivityStatsURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["activityStats"] = cas
		}
		if cs, err := client.get(charStatsURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["stats"] = cs
		}
		if cu, err := client.get(charUniqueWeaponsStateURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["uniqueWeapons"] = cu
		}
		if ci, err := client.get(charInventoryURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["inventory"] = ci
			// Enumerate Characters Inventiry -- maybe not important?
			//		InventoryItem
		}
		characterDoc[c.CharacterBase.CharacterId] = cdoc
	}
	u.set("characters", characterDoc, nil)

	// Merge into single document
	return nil
}
