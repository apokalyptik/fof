package main

import (
	"encoding/json"
	"strings"

	"github.com/apokalyptik/fof/lib/ubisoft/uplay"
)

type statDetail struct {
	info string
	name string
}

var statDetailMap = map[string]*statDetail{}

func getStatDetail(s *uplay.DivisionStat) (*statDetail, int) {
	val, err := s.Value.Float64()
	if sd, ok := statDetailMap[s.Name]; ok {
		return sd, int(val * 100)
	}
	buf, err := json.Marshal(map[string]interface{}{
		"id":           s.ID,
		"gameCode":     s.GameCode,
		"gameModeId":   s.GameModeID,
		"gameModeName": s.GameModeName,
		"gameName":     s.GameName,
		"isHighScore":  s.IsHighScore,
		"name":         s.Name,
		"unitLabel":    s.UnitLabel,
		"iconUrl":      s.IconURL,
		"multipliedBy": 100,
	})
	if err != nil {
		panic(err)
	}
	statDetailMap[s.Name] = &statDetail{
		info: string(buf),
		name: strings.Replace(strings.ToLower(s.Name), " ", "-", -1),
	}
	if err != nil {
		panic(err)
	}
	return statDetailMap[s.Name], int(val * 100)
}
