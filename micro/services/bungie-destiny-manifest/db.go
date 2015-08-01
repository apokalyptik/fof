package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/daaku/go.httpzip"
	_ "github.com/mattn/go-sqlite3"
)

var tables = []string{
	"DestinyActivityBundleDefinition",
	"DestinyPlaceDefinition",
	"DestinyActivityDefinition",
	"DestinyProgressionDefinition",
	"DestinyActivityTypeDefinition",
	"DestinyRaceDefinition",
	"DestinyClassDefinition",
	"DestinyRewardSourceDefinition",
	"DestinyDestinationDefinition",
	"DestinySandboxPerkDefinition",
	"DestinyDirectorBookDefinition",
	"DestinyScriptedSkullDefinition",
	"DestinyEnemyRaceDefinition",
	"DestinySpecialEventDefinition",
	"DestinyFactionDefinition",
	"DestinyStatDefinition",
	"DestinyGenderDefinition",
	"DestinyStatGroupDefinition",
	"DestinyGrimoireCardDefinition",
	"DestinyTalentGridDefinition",
	"DestinyGrimoireDefinition",
	"DestinyTriumphSetDefinition",
	"DestinyHistoricalStatsDefinition",
	"DestinyUnlockFlagDefinition",
	"DestinyInventoryBucketDefinition",
	"DestinyVendorCategoryDefinition",
	"DestinyInventoryItemDefinition",
	"DestinyVendorDefinition",
	"DestinyItemCategoryDefinition",
}

type db struct {
	path     string
	tempPath string
	handle   *sql.DB
	lock     sync.RWMutex
	version  string
}

var manifest = &db{}

func (d *db) query(q string, i interface{}) error {
	rows, err := d.handle.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(i)
		return err
	}
	return sql.ErrNoRows
}

func (d *db) update() error {
	var manifestListResponse = &struct {
		Status   string `json:"ErrorStatus"`
		Response struct {
			Version                 string            `json:"version"`
			MobileWorldContentPaths map[string]string `json:"mobileWorldContentPaths"`
		}
	}{}
	resp, err := http.Get(findManifestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	listDec := json.NewDecoder(resp.Body)
	if err := listDec.Decode(&manifestListResponse); err != nil {
		return err
	}

	if manifestListResponse.Status != "Success" {
		return errors.New(manifestListResponse.Status)
	}

	if manifestListResponse.Response.Version == d.version {
		return nil
	}

	zr, err := httpzip.ReadURL(fmt.Sprintf(manifestBaseURL, manifestListResponse.Response.MobileWorldContentPaths[lang]))
	if err != nil {
		return err
	}

	tf, err := os.OpenFile(d.tempPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		io.Copy(tf, rc)
		rc.Close()
	}

	tf.Close()

	if err := os.Rename(d.tempPath, d.path); err != nil {
		return err
	}

	dbh, err := sql.Open("sqlite3", d.path)
	if err != nil {
		return err
	}

	d.lock.Lock()
	d.handle = dbh
	d.version = manifestListResponse.Response.Version
	d.lock.Unlock()
	return nil
}
