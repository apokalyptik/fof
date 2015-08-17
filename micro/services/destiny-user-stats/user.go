package main

import (
	"encoding/json"
	"errors"

	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

var errUserNotFound = errors.New("User Not Found")

type userDB struct {
	list map[string]*user
	lock sync.RWMutex
}

func (u *userDB) update() {
	log.Println("Refreshing User List")
	u.lock.Lock()
	defer u.lock.Unlock()
	resp, err := http.Get(userListAddress)
	if err != nil {
		log.Printf("Error fetching userlist from %s: %s", userListAddress, err.Error())
		return
	}
	defer resp.Body.Close()
	var data struct {
		Members []struct {
			UserName string `json:"name"`
			UserID   string `json:"id"`
			Deleted  bool   `json:"deleted"`
			Bot      bool   `json:"is_bot"`
			Profile  struct {
				Gamertag string `json:"first_name"`
			} `json:"profile"`
		} `json:"members"`
	}
	dec := json.NewDecoder(resp.Body)
	if dec.Decode(&data); err != nil {
		log.Printf("Error decoding userlist: %s", err.Error())
		return
	}
	found := map[string]bool{}
	re := regexp.MustCompile("[^0-9a-z- ]")
	for _, m := range data.Members {
		if m.Deleted {
			continue
		}
		if m.Bot {
			continue
		}
		name := strings.TrimSpace(m.Profile.Gamertag)
		key := re.ReplaceAllString(strings.ToLower(name), "_")
		found[key] = true
		if _, ok := u.list[key]; !ok {
			u.list[key] = &user{
				name:     name,
				userName: m.UserName,
				userID:   m.UserID,
				data:     map[string]*userBit{},
			}
			log.Printf("+ %s", key)
		}
	}
	for name, _ := range u.list {
		if _, ok := found[name]; !ok {
			log.Printf("- %s", name)
			delete(u.list, name)
		}
	}

}

type userBit struct {
	when time.Time
	data interface{}
}

type user struct {
	name     string
	userName string
	userID   string
	platform int
	account  string
	data     map[string]*userBit
}

func (u *user) cUpsert(collection string, query, data interface{}) error {
	_, err := mgoDB.DB("fof").C(collection).Upsert(query, data)
	return err
}

func (u *user) pull() error {
	var search = []struct {
		MembershipType int    `mapstructure:"membershipType"`
		MembershipId   string `mapstructure:"membershipId"`
	}{}
	// Try XBL First
	info, buf, err := client.get(playerURL(platformXBL, u.name))
	if err != nil {
		return err
	} else {
		if err := mapstructure.Decode(info, &search); err != nil {
			return err
		}
	}
	if len(search) < 1 {
		// Try PSN Second
		info, buf, err = client.get(playerURL(platformPSN, u.name))
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
	u.cUpsert(
		"userData",
		bson.M{
			"account": u.account,
			"name":    u.name,
			"type":    "player",
		},
		map[string]interface{}{
			"username": u.userName,
			"userid":   u.userID,
			"account":  u.account,
			"name":     u.name,
			"type":     "player",
			"data":     info.([]interface{})[0],
			"raw":      buf,
		},
	)
	u.platform = search[0].MembershipType
	u.account = search[0].MembershipId
	u.data["player"] = &userBit{
		when: time.Now(),
		data: info.([]interface{})[0],
	}

	astats, buf, err := client.get(accountStatsURL(u.platform, u.account))
	if err == nil {
		u.cUpsert(
			"userData",
			bson.M{
				"account": u.account,
				"name":    u.name,
				"type":    "accountStats",
			},
			map[string]interface{}{
				"username": u.userName,
				"userid":   u.userID,
				"account":  u.account,
				"name":     u.name,
				"type":     "accountStats",
				"data":     astats,
				"raw":      buf,
			},
		)
		if astatsDocs, err := pullAllTimeStatsDocs(u.name, u.account, astats); err != nil {
			log.Printf("Error pulling all time stats: %s", err.Error())
		} else {
			for _, v := range astatsDocs {
				if err := u.cUpsert("accountStats", bson.M{"account": u.account, "section": v.Section, "stat": v.Stat}, v); err != nil {
					log.Printf("Error inserting into accountStats: %s", err.Error())
				}
			}
		}
	}

	grim, buf, err := client.get(grimoireURL(u.platform, u.account))
	if err == nil {
		u.cUpsert(
			"userData",
			bson.M{
				"account": u.account,
				"name":    u.name,
				"type":    "grimoire",
			},
			map[string]interface{}{
				"username": u.userName,
				"userid":   u.userID,
				"account":  u.account,
				"name":     u.name,
				"type":     "grimoire",
				"data":     grim,
				"raw":      buf,
			},
		)
	}

	t, buf, err := client.get(triumphsURL(u.platform, u.account))
	if err == nil {
		u.cUpsert(
			"userData",
			bson.M{
				"account": u.account,
				"name":    u.name,
				"type":    "triumphs",
			},
			map[string]interface{}{
				"username": u.userName,
				"userid":   u.userID,
				"account":  u.account,
				"name":     u.name,
				"type":     "triumphs",
				"data":     t,
				"raw":      buf,
			},
		)
	}

	account, buf, err := client.get(accountURL(u.platform, u.account))
	if err == nil {
		u.cUpsert(
			"userData",
			bson.M{
				"account": u.account,
				"name":    u.name,
				"type":    "account",
			},
			map[string]interface{}{
				"username": u.userName,
				"userid":   u.userID,
				"account":  u.account,
				"name":     u.name,
				"type":     "account",
				"data":     account,
				"raw":      buf,
			},
		)
	}

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
	var characterDoc map[string]map[string]interface{}
	mapstructure.Decode(u.data["characters"].data, &characterDoc)
	for _, c := range acctResponse.Data.Characters {
		if _, ok := characterDoc[c.CharacterBase.CharacterId]; !ok {
			characterDoc[c.CharacterBase.CharacterId] = map[string]interface{}{}
		}
		cdoc := characterDoc[c.CharacterBase.CharacterId]

		if ch, _, err := client.get(charActivityHistoryURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["activityHistory"] = ch
		}
		if ca, _, err := client.get(charActivitiesURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["activities"] = ca
		}
		if cp, _, err := client.get(charProgressionURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["progression"] = cp
		}
		if cas, _, err := client.get(charActivityStatsURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["activityStats"] = cas
		}
		if cs, _, err := client.get(charStatsURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["stats"] = cs
		}
		if cu, _, err := client.get(charUniqueWeaponsStateURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["uniqueWeapons"] = cu
		}
		if ci, _, err := client.get(charInventoryURL(u.platform, u.account, c.CharacterBase.CharacterId)); err == nil {
			cdoc["inventory"] = ci
			// Enumerate Characters Inventiry -- maybe not important?
			//		InventoryItem
		}
		buf, _ := json.Marshal(cdoc)
		u.cUpsert(
			"userData",
			bson.M{
				"account": u.account,
				"name":    u.name,
				"type":    "character-" + c.CharacterBase.CharacterId,
			},
			map[string]interface{}{
				"username": u.userName,
				"userid":   u.userID,
				"account":  u.account,
				"name":     u.name,
				"type":     "character-" + c.CharacterBase.CharacterId,
				"data":     cdoc,
				"raw":      buf,
			},
		)
		characterDoc[c.CharacterBase.CharacterId] = cdoc
	}

	buf, _ = json.Marshal(characterDoc)
	u.cUpsert(
		"userData",
		bson.M{
			"account": u.account,
			"name":    u.name,
			"type":    "characters",
		},
		map[string]interface{}{
			"username": u.userName,
			"userid":   u.userID,
			"account":  u.account,
			"name":     u.name,
			"type":     "characters",
			"data":     characterDoc,
			"raw":      buf,
		},
	)

	exoticStatsDocs, err := pullExoticStatsDocs(u.name, u.account, characterDoc)
	if err != nil {
		log.Printf("Error pulling exotic stats for %s: %s", u.name, err.Error())
	} else {
		for _, v := range exoticStatsDocs {
			if err := u.cUpsert("accountStats", bson.M{"account": u.account, "section": v.Section, "stat": v.Stat}, v); err != nil {
				log.Printf("Error inserting into accountStats: %s", err.Error())
			}
		}
	}

	// Merge into single document
	return nil
}
