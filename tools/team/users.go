package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var errSlackUserIDNotFound = errors.New("User not found on slack")
var errSlackIMNotCreated = errors.New("Could not create slack IM")
var errSlackGTNotFound = errors.New("GT not found on slack")

var udb = &userDB{}

type userDB struct {
	filename     string
	lock         sync.Mutex
	LastUpdated  time.Time
	APITokenHash []byte
	LookupID     map[string]string `json:"LookupUserToID"`
	IMs          map[string]string `json:"LookupUserToIM"`
	LookupGT     map[string]string `json:"LookupGT"`
}

func (u *userDB) getGamertagForUser(username string) (string, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if gt, ok := u.LookupGT[username]; ok {
		return gt, nil
	}

	if LookupID, LookupGT, err := slack.getUserDataMap(); err != nil {
		return "", err
	} else {
		u.LookupID = LookupID
		u.LookupGT = LookupGT
	}

	if gt, ok := u.LookupGT[username]; ok {
		return gt, nil
	}

	return "", errSlackGTNotFound
}

func (u *userDB) getChannelForIM(username string) (string, error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	channel, ok := u.IMs[username]
	if ok && channel != "" {
		return channel, nil
	}

	userid, ok := u.LookupID[username]
	if !ok {
		if LookupID, LookupGT, err := slack.getUserDataMap(); err != nil {
			return "", err
		} else {
			u.LookupID = LookupID
			u.LookupGT = LookupGT
			userid, ok = u.LookupID[username]
			if !ok {
				return "", errSlackUserIDNotFound
			}
		}
	}

	channel, err := slack.doOpenIM(userid)
	if err != nil {
		log.Println("Got error opening IM to", username, err.Error())
		return "", err
	}
	u.IMs[username] = channel

	if err := u.save(); err != nil {
		log.Println("Error saving users db:", err.Error())
	}

	return channel, nil
}

func (u *userDB) save() error {
	if u.filename == "" {
		return errors.New("cannot persist withuot filename")
	}
	fp, err := os.Create(u.filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if err := enc.Encode(u); err != nil {
		return err
	}
	return nil
}

func (u *userDB) load(filename string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			u.filename = filename
			return nil
		}
		return err
	}
	u.filename = filename
	defer fp.Close()
	dec := json.NewDecoder(fp)
	if err := dec.Decode(&u); err != nil {
		return err
	}
	if u == nil {
		u = &userDB{
			LookupID: map[string]string{},
			IMs:      map[string]string{},
		}
	}
	h := md5.New()
	io.WriteString(h, slack.apiKey)
	hash := h.Sum(nil)
	if bytes.Equal(hash, u.APITokenHash) {
		return nil
	}
	u.APITokenHash = hash
	u.LookupID = map[string]string{}
	u.LookupGT = map[string]string{}
	u.IMs = map[string]string{}
	return nil
}
