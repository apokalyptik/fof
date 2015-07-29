package main

import "github.com/mitchellh/mapstructure"

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
