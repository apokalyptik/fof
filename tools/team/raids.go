package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"crypto/hmac"
	"crypto/sha256"

	"code.google.com/p/go-uuid/uuid"
)

var errDisbanded = errors.New("Since you were the last member of \"%s\" on #%s the raid has been disbanded")

type raidCommandResponses map[string]string

func (r raidCommandResponses) sendToSlack() {
	for k, v := range r {
		if k == "-" {
			continue
		}
		slack.msg().to(k).send(v)
	}
}

func (r raidCommandResponses) stdOut() string {
	if out, ok := r["-"]; ok {
		return out
	}
	return ""
}

var raidSlashCommand = "/raid"

var raidDb = &raids{
	data: map[string][]*raid{},
}

type raid struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Members   []string  `json:"members"`
	Alts      []string  `json:"alts"`
	UUID      string    `json:"uuid"`
	Secret    string    `json:"secret"`
	Type      string    `json:"type"`
}

func (r *raid) hmacForUser(username string) string {
	mac := hmac.New(sha256.New, []byte(r.Secret))
	fmt.Fprintf(mac, "@%s:%s:%s", username, r.Name, r.UUID)
	expectedMAC := mac.Sum(nil)
	return hex.EncodeToString(expectedMAC[8:18])
}

func (r *raid) validateHmacForUser(username, hm string) error {
	want := r.hmacForUser(username)
	if want == hm {
		return nil
	}
	log.Printf("'%s' != '%s'", hm, want)
	return errors.New("Invalid HMAC")
}

type raids struct {
	filename string
	data     map[string][]*raid
	lock     sync.RWMutex
}

func (r *raids) joinAlt(channel, name, user string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	if len(r.data[channel]) == 0 {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	for _, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		v.Alts = append(v.Alts, user)
		return v.Members[0], nil
	}
	return "", fmt.Errorf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \""+raidSlashCommand+" host\" one?",
		name,
		channel)
}

func (r *raids) join(channel, name, user string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	if len(r.data[channel]) == 0 {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	for _, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		v.Members = append(v.Members, user)
		return v.Members[0], nil
	}
	return "", fmt.Errorf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \""+raidSlashCommand+" host\" one?",
		name,
		channel)
}

func (r *raids) leaveAlt(channel, name, user string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	if len(r.data[channel]) == 0 {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	for _, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		toRemove := -1
		for i := range v.Alts {
			checkOn := (len(v.Alts) - i) - 1
			if v.Alts[checkOn] == user {
				toRemove = checkOn
				break
			}
		}
		if toRemove < 0 {
			return "", fmt.Errorf(
				"You are not signed up to do \"%s\" on #%s",
				name,
				channel)
		}
		var newAlts = []string{}
		for i, n := range v.Alts {
			if i == toRemove {
				continue
			}
			newAlts = append(newAlts, n)
		}
		v.Alts = newAlts
		return v.Members[0], nil
	}
	return "", fmt.Errorf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \""+raidSlashCommand+" host\" one?",
		name,
		channel)
}

func (r *raids) leave(channel, name, user string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	if len(r.data[channel]) == 0 {
		return "", fmt.Errorf("There are no raids for #%s", channel)
	}
	for vid, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		toRemove := -1
		for i := range v.Members {
			checkOn := (len(v.Members) - i) - 1
			if v.Members[checkOn] == user {
				toRemove = checkOn
				break
			}
		}
		if toRemove < 0 {
			return "", fmt.Errorf(
				"You are not signed up to do \"%s\" on #%s",
				name,
				channel)
		}
		var newMembers = []string{}
		for i, n := range v.Members {
			if i == toRemove {
				continue
			}
			newMembers = append(newMembers, n)
		}
		v.Members = newMembers
		if len(v.Members) == 0 {
			r.data[channel] = append(r.data[channel][:vid], r.data[channel][vid+1:]...)
			if len(r.data[channel]) == 0 {
				delete(r.data, channel)
			}
			return "", errDisbanded
		}
		return v.Members[0], nil
	}
	return "", fmt.Errorf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \""+raidSlashCommand+" host\" one?",
		name,
		channel)
}

func (r *raids) members(channel, name string) ([]string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if _, ok := r.data[channel]; !ok {
		return nil, fmt.Errorf("There are no raids for #%s", channel)
	}
	if len(r.data[channel]) == 0 {
		return nil, fmt.Errorf("There are no raids for #%s", channel)
	}
	for _, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		rval := make([]string, len(v.Members))
		for k, n := range v.Members {
			rval[k] = n
		}
		return rval, nil
	}
	return nil, fmt.Errorf(
		"I have no \"%s\" registered for #%s",
		name,
		channel)
}

func (r *raids) finish(channel, name, user string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	for k, v := range r.data[channel] {
		if v.Name != name {
			continue
		}
		var allowed = false
		if v.Members[0] == user {
			allowed = true
		} else {
			for _, admin := range admins {
				if user == admin {
					allowed = true
					break
				}
			}
		}
		if allowed == false {
			return fmt.Errorf(
				"Only the organizer (_@%s_) can finish a raid",
				v.Members[0])
		}
		r.data[channel] = append(r.data[channel][:k], r.data[channel][k+1:]...)
		return nil
	}
	return fmt.Errorf(
		"I have no \"%s\" registered for #%s. Perhaps you would like to \"/"+raidSlashCommand+" host\" one?",
		name,
		channel)
}

func (r *raids) register(channel, name, user string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.save()
	if _, ok := r.data[channel]; !ok {
		r.data[channel] = []*raid{}
	}
	for _, v := range r.data[channel] {
		if v.Name == name {
			return errors.New("A raid by this name is already registered")
		}
	}
	r.data[channel] = append(r.data[channel], &raid{
		Name:      name,
		CreatedAt: time.Now(),
		Members:   []string{user},
		UUID:      uuid.New(),
		Secret:    uuid.New(),
	})
	return nil
}

func (r *raids) list(channel string) []*raid {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if c, ok := r.data[channel]; ok {
		rval := make([]*raid, len(c))
		for k, v := range c {
			rval[k] = v
		}
		return rval
	}
	return nil
}

func (r *raids) cache() {
	var data = map[string]map[string]map[string]interface{}{}
	for channel, channelRaids := range raidDb.data {
		data[channel] = map[string]map[string]interface{}{}
		for _, raid := range channelRaids {
			data[channel][raid.UUID] = map[string]interface{}{
				"uuid":       raid.UUID,
				"name":       raid.Name,
				"members":    raid.Members,
				"alts":       raid.Alts,
				"created_at": raid.CreatedAt,
			}
		}
	}
	xhrOutput.set("raids", data)
}

func (r *raids) save() error {
	if r.filename == "" {
		return errors.New("cannot persist withuot filename")
	}
	fp, err := os.Create(r.filename)
	if err != nil {
		return err
	}
	defer r.cache()
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if err := enc.Encode(r.data); err != nil {
		return err
	}
	return nil
}

func (r *raids) mindExpiration(maxAge time.Duration) {
	ticker := time.Tick(10 * time.Minute)
	for {
		<-ticker
		r.lock.RLock()
		for channel, raidlist := range r.data {
			for _, raidentry := range raidlist {
				switch raidentry.Type {
				case "event":
					if time.Now().Add(0 - maxAge).After(raidentry.CreatedAt) {
						go r.finish(channel, raidentry.Name, raidentry.Members[0])
						log.Printf("Expiring %s on #%s", raidentry.Name, channel)
					}
				}
			}
		}
		r.lock.RUnlock()
	}
}

func (r *raids) load(filename string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	defer r.cache()
	fp, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			r.filename = filename
			return nil
		}
		return err
	}
	r.filename = filename
	defer fp.Close()
	dec := json.NewDecoder(fp)
	if err := dec.Decode(&r.data); err != nil {
		return err
	}
	if r.data == nil {
		r.data = map[string][]*raid{}
	}
	for c := range r.data {
		for i := range r.data[c] {
			if r.data[c][i].UUID == "" {
				r.data[c][i].UUID = uuid.New()
				r.data[c][i].Secret = uuid.New()
			}
			if r.data[c][i].Type == "" {
				r.data[c][i].Type = "event"
			}
		}
	}
	return nil
}

func raidHost(username, channel, raid string) (raidCommandResponses, error) {
	if err := raidDb.register(channel, raid, username); err != nil {
		return raidCommandResponses{
			"@" + username: err.Error(),
			"-":            err.Error(),
		}, err
	}
	return raidCommandResponses{
		"@" + username: fmt.Sprintf(
			"*%s* has been hosted on *#%s*. You will be notified when members join you",
			raid, channel),
		"-": fmt.Sprintf(
			"*%s* has been hosted on *#%s*. You will be notified when members join you",
			raid, channel),
		"#" + channel: fmt.Sprintf(
			"*@%s* is hosting a new event *%s*. Use the <http://fofgaming.com/team|team tool> or type ‘/team’ to sign up",
			username, raid),
	}, nil
}

func raidPing(username, channel, raid string) (raidCommandResponses, error) {
	list, err := raidDb.members(channel, raid)
	if err != nil {
		return raidCommandResponses{
			"@" + username: err.Error(),
			"-":            err.Error(),
		}, err
	}
	return raidCommandResponses{
		"#" + channel: fmt.Sprintf(
			"pinging @%s about *%s*",
			strings.Join(list, ", @"), raid),
	}, nil
}

func raidFinish(username, channel, raid string) (raidCommandResponses, error) {
	if err := raidDb.finish(channel, raid, username); err != nil {
		return raidCommandResponses{
			"@" + username: err.Error(),
			"-":            err.Error(),
		}, err
	}
	return raidCommandResponses{
		"@" + username: fmt.Sprintf(
			"*%s* has been removed from #%s",
			raid, channel),
		"#" + channel: fmt.Sprintf(
			"*@%s* has removed *%s*. Use the <http://fofgaming.com/team|team tool> or type ‘/team’ to host/join a new event",
			username, raid),
	}, nil
}

func raidLeave(username, channel, raid string) (raidCommandResponses, error) {
	owner, err := raidDb.leave(channel, raid, username)
	if err != nil {
		if err == errDisbanded {
			return raidCommandResponses{
				"@" + username: fmt.Sprintf(err.Error(), raid, channel),
			}, nil
		}
		return raidCommandResponses{
			"@" + username: err.Error(),
			"-":            err.Error(),
		}, err
	}
	return raidCommandResponses{
		"@" + username: fmt.Sprintf(
			"You have removed yourself from *%s* on #%s",
			raid, channel),
		"#" + channel: fmt.Sprintf(
			"*@%s* has removed themselves from *%s*. Use the <http://fofgaming.com/team/|team tool> or type ‘/team’ to take their place",
			username, raid),
		"@" + owner: fmt.Sprintf(
			"@%s has removed themselves from *%s* on #%s",
			username, raid, channel),
	}, nil
}

func raidAltLeave(username, channel, raid string) (raidCommandResponses, error) {
	owner, err := raidDb.leaveAlt(channel, raid, username)
	if err != nil {
		return raidCommandResponses{
			"@" + username: err.Error(),
			"-":            err.Error(),
		}, err
	}
	return raidCommandResponses{
		"@" + username: fmt.Sprintf(
			"You have removed yourself from *%s* on #%s",
			raid, channel),
		"#" + channel: fmt.Sprintf(
			"*@%s* has removed themselves as an alternate from *%s*. Use the <http://fofgaming.com/team/|team tool> or type ‘/team’ to take their place",
			username, raid),
		"@" + owner: fmt.Sprintf(
			"*@%s* is no longer signed up as an alternate for *%s* on #%s",
			username, raid, channel),
	}, nil
}

func raidJoin(username, channel, raid string) (raidCommandResponses, error) {
	owner, err := raidDb.join(channel, raid, username)
	if err != nil {
		return raidCommandResponses{
			"@" + username: err.Error(),
			"-":            err.Error(),
		}, err
	}
	return raidCommandResponses{
		"@" + username: fmt.Sprintf(
			"You have signed up for *%s* on #%s",
			raid, channel),
		"#" + channel: fmt.Sprintf(
			"*@%s* has signed for *%s* Use the <http://fofgaming.com/team/|team tool> or type ‘/team’ to join them",
			username, raid),
		"@" + owner: fmt.Sprintf(
			"*@%s* has signed up to join you in *%s* on #%s",
			username, raid, channel),
	}, nil
}

func raidAltJoin(username, channel, raid string) (raidCommandResponses, error) {
	owner, err := raidDb.joinAlt(channel, raid, username)
	if err != nil {
		return raidCommandResponses{
			"@" + username: err.Error(),
			"-":            err.Error(),
		}, err
	}
	return raidCommandResponses{
		"@" + username: fmt.Sprintf(
			"You have signed up to be an alternate in *%s* on #%s",
			raid, channel),
		"#" + channel: fmt.Sprintf(
			"*@%s* has signed for *%s* as an alternate. Use the <http://fofgaming.com/team|team tool> or type ‘/team’ to join them",
			username, raid),
		"@" + owner: fmt.Sprintf(
			"*@%s* has signed up to be an alternate in *%s* on #%s",
			username, raid, channel),
	}, nil
}
