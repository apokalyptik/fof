package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type keyedCache struct {
	lock        sync.RWMutex
	cacheKey    string
	cacheExpiry time.Time
	cache       []byte
}

type timedCache struct {
	lock        sync.RWMutex
	cacheExpiry time.Time
	cache       []byte
}

type cacheAllTimeStats struct {
	lock sync.RWMutex
	data map[string]map[string]*cacheAllTimeStat
}

func (c *cacheAllTimeStats) get(section, stat string) ([]byte, error) {
	c.lock.RLock()
	if _, ok := c.data[section]; !ok {
		c.lock.RUnlock()
		c.lock.Lock()
		c.data[section] = map[string]*cacheAllTimeStat{}
		c.data[section][stat] = &cacheAllTimeStat{
			keyedCache: &keyedCache{},
		}
		c.lock.Unlock()
		c.lock.RLock()
	}
	if _, ok := c.data[section][stat]; !ok {
		c.lock.RUnlock()
		c.lock.Lock()
		c.data[section][stat] = &cacheAllTimeStat{
			keyedCache: &keyedCache{},
		}
		c.lock.Unlock()
		c.lock.RLock()
	}
	c.lock.RUnlock()
	return c.data[section][stat].get(section, stat)
}

type cacheAllTimeStat struct {
	*keyedCache
}

func (c *cacheAllTimeStat) get(section, stat string) ([]byte, error) {
	c.lock.RLock()
	if time.Now().Before(c.cacheExpiry) {
		if c.cacheKey != "" {
			if c.cache != nil {
				c.lock.RUnlock()
				return c.cache, nil
			}
		}
	}
	var newCacheKeyID struct {
		ID bson.ObjectId `bson:"_id"`
	}
	var newCacheKey string

	if err := mdb.DB("fof").C("accountStats").Find(nil).Sort("-_id").Limit(1).One(&newCacheKeyID); err != nil {
		c.lock.RUnlock()
		log.Printf("Error getting latest _id from accountStats: %s", err.Error())
		if c.cache != nil {
			return c.cache, nil
		}
		return nil, err
	}

	newCacheKey = newCacheKeyID.ID.String()

	c.lock.RUnlock()
	c.lock.Lock()
	defer c.lock.Unlock()

	if newCacheKey == c.cacheKey {
		c.cacheExpiry = time.Now().Add(5 * time.Minute)
		return c.cache, nil
	}

	query := bson.M{
		"section": section,
		"stat":    stat,
	}
	var docs []struct {
		Member string  `json:"member",bson:"member"`
		Value  float64 `json:"value",bson:"value"`
		PGA    string  `json:"pga",bson:"pgadisplayvalue"`
	}
	if err := mdb.DB("fof").C("accountStats").Find(query).Limit(10000).Sort("-value", "-pgavalue").All(&docs); err != nil {
		log.Printf("Error fetching allTimeStatKeys: %s", err.Error())
		if c.cache != nil {
			return c.cache, nil
		}
		return nil, err
	}
	if cache, err := json.Marshal(docs); err != nil {
		log.Printf("Error marshalling all time stats cache for %s/%s: %s", section, stat, err.Error())
		if c.cache != nil {
			return c.cache, nil
		}
		return nil, err
	} else {
		c.cache = cache
		c.cacheKey = newCacheKey
		c.cacheExpiry = time.Now().Add(5 * time.Minute)
	}
	return c.cache, nil
}

type cacheAllTimeKeys struct {
	*keyedCache
}

func (c *cacheAllTimeKeys) get() ([]byte, error) {
	c.lock.RLock()
	if time.Now().Before(c.cacheExpiry) {
		if c.cacheKey != "" {
			if c.cache != nil {
				c.lock.RUnlock()
				return c.cache, nil
			}
		}
	}
	var newCacheKeyID struct {
		ID bson.ObjectId `bson:"_id"`
	}
	var newCacheKey string

	if err := mdb.DB("fof").C("accountStats").Find(nil).Sort("-_id").Limit(1).One(&newCacheKeyID); err != nil {
		c.lock.RUnlock()
		log.Printf("Error getting latest _id from accountStats: %s", err.Error())
		if c.cache != nil {
			return c.cache, nil
		}
		return nil, err
	}

	newCacheKey = newCacheKeyID.ID.String()

	c.lock.RUnlock()
	c.lock.Lock()
	defer c.lock.Unlock()

	if newCacheKey == c.cacheKey {
		c.cacheExpiry = time.Now().Add(5 * time.Minute)
		return c.cache, nil
	}

	var sections = []string{}
	var rval = map[string][]string{}

	if err := mdb.DB("fof").C("accountStats").Find(nil).Distinct("section", &sections); err != nil {
		log.Printf("Error fetching allTimeStatKeys: %s", err.Error())
		if c.cache != nil {
			return c.cache, nil
		}
		return nil, err
	}
	for _, s := range sections {
		var sectionStats = []string{}
		if err := mdb.DB("fof").C("accountStats").Find(&bson.M{"section": s}).Distinct("stat", &sectionStats); err != nil {
			continue
		}
		rval[s] = sectionStats
	}
	if cache, err := json.Marshal(rval); err != nil {
		log.Printf("Error marshalling keycache: %s", err.Error())
		if c.cache != nil {
			return c.cache, nil
		}
		return nil, err
	} else {
		c.cache = cache
		c.cacheKey = newCacheKey
		c.cacheExpiry = time.Now().Add(5 * time.Minute)
	}
	return c.cache, nil
}

type cacheLeaderboardPVP struct {
	*timedCache
}

func (c *cacheLeaderboardPVP) get() ([]byte, error) {
	c.lock.RLock()
	if time.Now().Before(c.cacheExpiry) {
		if c.cache != nil {
			c.lock.RUnlock()
			return c.cache, nil
		}
	}
	c.lock.RUnlock()

	c.lock.Lock()
	defer c.lock.Unlock()

	var data []struct {
		Member string
		Stat   string
		Value  float64
	}

	err := mdb.DB("fof").C("accountStats").Find(bson.M{
		"section": "pvp",
		"stat": bson.M{
			"$in": []string{
				"activitiesEntered",
				"secondsPlayed",
				"activitiesWon",
				"allParticipantsScore",
				"assists",
				"bestSingleGameKills",
				"bestSingleGameScore",
				"deaths",
				"kills",
				"longest KillSpree",
				"longest SingleLife",
				"orbsDropped",
				"orbsGathered",
				"precision Kills",
				"weaponKillsGrenade",
				"weaponKillsMelee",
				"weaponKillsSuper",
				"zonesCaptured",
			},
		},
	}).Select(bson.M{"member": true, "stat": true, "value": true, "_id": false}).All(&data)

	if err != nil {
		return nil, err
	}

	var rval = map[string]map[string]float64{}
	for _, v := range data {
		if _, ok := rval[v.Stat]; !ok {
			rval[v.Stat] = map[string]float64{}
		}
		rval[v.Stat][v.Member] = v.Value
	}

	if cache, err := json.Marshal(rval); err != nil {
		return nil, err
	} else {
		c.cache = cache
		c.cacheExpiry = time.Now().Add(30 * time.Minute)
		return cache, nil
	}
}

var atsCache = &cacheAllTimeStats{
	data: map[string]map[string]*cacheAllTimeStat{},
}

var atkCache = &cacheAllTimeKeys{
	keyedCache: &keyedCache{},
}

var pvpLeaderboardCache = &cacheLeaderboardPVP{
	timedCache: &timedCache{},
}

func init() {
	go func() {
		wake := time.Tick(10 * time.Minute)
		for {
			<-wake
			atsCache.lock.RLock()
			for _, k1 := range atsCache.data {
				for _, k2 := range k1 {
					if time.Now().After(k2.cacheExpiry) {
						k2.lock.Lock()
						k2.cache = nil
						k2.lock.Unlock()
					}
				}
			}
			atsCache.lock.RUnlock()
		}
	}()
}
