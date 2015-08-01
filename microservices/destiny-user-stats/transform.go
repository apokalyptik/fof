package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type accountStatDetail struct {
	DisplayValue string  `mapstructure:"displayValue"`
	Value        float64 `mapstructure:"value"`
}

type accountStat struct {
	Basic  accountStatDetail `mapstructure:"basic"`
	PGA    accountStatDetail `mapstructure:"pga"`
	StatID string            `mapstructure:"statId"`
}

type accountStatsDecoded struct {
	All struct {
		Merged struct {
			AllTime map[string]accountStat `mapstructure:"allTime"`
		} `mapstructure:"merged"`
		Results struct {
			PVP struct {
				AllTime map[string]accountStat `mapstructure:"allTime"`
			} `mapstructure:"allPvP"`
			PVE struct {
				AllTime map[string]accountStat `mapstructure:"allTime"`
			} `mapstructure:"allPvE"`
		} `mapstructure:"results"`
	} `mapstructure:"mergedAllCharacters"`
}

type accountStatsDoc struct {
	Member          string
	Account         string
	Stat            string
	Section         string
	DisplayValue    string
	Value           float64
	PGADisplayValue string
	PGAValue        float64
}

func makeAllTimeStatDoc(member, account, section string, v accountStat) accountStatsDoc {
	return accountStatsDoc{
		Member:          member,
		Account:         account,
		Stat:            v.StatID,
		Section:         section,
		DisplayValue:    v.Basic.DisplayValue,
		Value:           v.Basic.Value,
		PGADisplayValue: v.PGA.DisplayValue,
		PGAValue:        v.PGA.Value,
	}
}

func pullAllTimeStatsDocs(member string, account string, cdoc interface{}) ([]accountStatsDoc, error) {
	var decoded accountStatsDecoded
	if err := mapstructure.Decode(cdoc, &decoded); err != nil {
		return nil, err
	}
	var rval = []accountStatsDoc{}
	for _, v := range decoded.All.Merged.AllTime {
		rval = append(rval, makeAllTimeStatDoc(member, account, "merged", v))
	}
	for _, v := range decoded.All.Results.PVP.AllTime {
		rval = append(rval, makeAllTimeStatDoc(member, account, "pvp", v))
	}
	for _, v := range decoded.All.Results.PVE.AllTime {
		rval = append(rval, makeAllTimeStatDoc(member, account, "pve", v))
	}
	return rval, nil
}

type charactersUniqueWeapons map[string]struct {
	Weapons struct {
		Data struct {
			Weapons []struct {
				ReferenceID float64 `mapstructure:"referenceId"`
				Name        string
				Values      map[string]struct {
					StatID string `mapstructure:"statId"`
					Basic  struct {
						DisplayValue string  `mapstructure:"displayValue"`
						Value        float64 `mapstructure:"value"`
					} `mapstructure:"basic"`
				} `mapstructure:"values"`
			} `mapstructure:"weapons"`
		} `mapstructure:"data"`
	} `mapstructure:"uniqueWeapons"`
}

var exoticHashLookup = map[float64]string{}

type exoticLookupResponse struct {
	OK       bool   `json:"ok"`
	Error    string `json:"error"`
	Response struct {
		Name string `json:"itemName"`
	} `json:"response"`
}

func lookupExoticHash(hash float64) (string, error) {
	if name, ok := exoticHashLookup[hash]; ok {
		return name, nil
	}
	resp, err := http.Get(fmt.Sprintf(
		"http://localhost:8883/destiny/DestinyInventoryItemDefinition/%0.f.json",
		hash,
	))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var unwrapped exoticLookupResponse
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&unwrapped); err != nil {
		return "", err
	}
	if !unwrapped.OK {
		return "", errors.New(unwrapped.Error)
	}
	exoticHashLookup[hash] = unwrapped.Response.Name
	return unwrapped.Response.Name, nil

}

func pullExoticStatsDocs(member string, account string, chars interface{}) ([]accountStatsDoc, error) {
	var decoded charactersUniqueWeapons
	var rval = map[float64]map[string]*accountStatsDoc{}
	mapstructure.Decode(chars, &decoded)
	for _, data := range decoded {
		for _, weapon := range data.Weapons.Data.Weapons {
			if name, err := lookupExoticHash(weapon.ReferenceID); err != nil {
				log.Printf(
					"Error looking up exotic weapon hash: %0.f: %s",
					weapon.ReferenceID,
					err.Error(),
				)
				continue
			} else {
				weapon.Name = name
			}
			if _, ok := rval[weapon.ReferenceID]; !ok {
				rval[weapon.ReferenceID] = map[string]*accountStatsDoc{}
			}
			for name, value := range weapon.Values {
				if _, ok := rval[weapon.ReferenceID][name]; !ok {
					rval[weapon.ReferenceID][name] = &accountStatsDoc{}
				}
				rval[weapon.ReferenceID][name].Member = member
				rval[weapon.ReferenceID][name].Account = account
				rval[weapon.ReferenceID][name].Section = "exoticWeapons"
				rval[weapon.ReferenceID][name].Value += value.Basic.Value
				rval[weapon.ReferenceID][name].Stat = strings.Replace(
					strings.Replace(value.StatID, "KillsPrecisionKills", "KillsToPrecisionKills", -1),
					"uniqueWeapon",
					strings.Replace(weapon.Name, " ", "", -1),
					1,
				)
				rval[weapon.ReferenceID][name].DisplayValue = fmt.Sprintf(
					"%0.f",
					rval[weapon.ReferenceID][name].Value,
				)
			}
		}
	}
	var returnValue = []accountStatsDoc{}
	for _, w := range rval {
		for _, s := range w {
			returnValue = append(returnValue, *s)
		}
	}
	return returnValue, nil
}
