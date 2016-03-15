package uplay

import "encoding/json"

type DivisionStat struct {
	ID           int         `json:"id"`
	GameCode     string      `json:"gameCode"`
	GameModeID   int         `json:"gameModeId"`
	GameModeName string      `json:"gameModeName"`
	GameName     string      `json:"gameName"`
	IconURL      string      `json:"iconUrl"`
	IsHighScore  int         `json:"isHighScore"`
	Name         string      `json:"name"`
	UnitLabel    string      `json:"unitLabel"`
	Value        json.Number `json:"value"`
}
