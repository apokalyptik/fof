package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
			Deleted bool `json:"deleted"`
			Bot     bool `json:"is_bot"`
			Profile struct {
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
	re := regexp.MustCompile("[^0-9a-z-]")
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
				name: name,
				data: map[string]*userBit{},
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
	platform int
	account  string
	data     map[string]*userBit
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

func (u *user) cUpsert(collection string, query, data interface{}) error {
	_, err := mgoDB.DB("fof").C(collection).Upsert(query, data)
	return err
}

func (u *user) pull() error {
	defer func() {
		var doc = map[string]interface{}{}
		for k, v := range u.data {
			doc[k] = v.data
			delete(u.data, k)
		}
		if bytes, err := json.Marshal(doc); err != nil {
			log.Printf("Error marshalling data for %s: %s", u.name, err.Error())
		} else {
			err = ioutil.WriteFile(fmt.Sprintf("data/%d.%s.json", u.platform, u.account), bytes, 0644)
			if err != nil {
				log.Printf("Error writing datafile for user data/%d.%s.json: %s", u.platform, u.account, err.Error())
			}
		}
		c := mgoDB.DB("fof").C("memberDestinyStats")
		_, err := c.Upsert(bson.M{"member": u.name}, &struct {
			Member string
			Data   map[string]interface{}
		}{
			Member: u.name,
			Data:   doc,
		})
		if err != nil {
			log.Printf("Error inserting document into MongoDB: %s", err.Error())
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

	if bytes, err := ioutil.ReadFile(fmt.Sprintf("data/%d.%s.json", u.platform, u.account)); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error reading %s: %s", fmt.Sprintf("data/%d.%s.json", u.platform, u.account), err.Error())
		}
	} else {
		var doc = map[string]interface{}{}
		if err := json.Unmarshal(bytes, &doc); err != nil {
			log.Printf("Error unmarshalling data for %d/%s (%s): %s", u.platform, u.account, u.name, err.Error())
		} else {
			for k, v := range doc {
				u.set(k, v, nil)
			}
		}
	}

	astats, err := client.get(accountStatsURL(u.platform, u.account))
	u.set("accountStats", astats, err)
	if err == nil {
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

	grim, err := client.get(grimoireURL(u.platform, u.account))
	u.set("grimoire", grim, err)

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
	var characterDoc map[string]map[string]interface{}
	mapstructure.Decode(u.data["characters"].data, &characterDoc)
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
	u.set("characters", characterDoc, nil)

	// Merge into single document
	return nil
}
