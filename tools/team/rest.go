package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

var errInvalidMethod = errors.New("Invalid Method")
var errInvalidAuth = errors.New("Invalid Authentication")
var hmacKey = make([]byte, 32)
var hmacEmpty = make([]byte, 32)

var store *sessions.CookieStore

func generateAPIKeyForUserTime(username string, age int) string {
	hm := hmac.New(sha1.New, hmacKey)
	fmt.Fprintln(hm, username, int(math.Ceil(float64(time.Now().Unix())/3600))-age)
	return fmt.Sprintf("%x", hm.Sum(nil))
}

func generateKeyForUserTime(username string, age int) string {
	hm := hmac.New(md5.New, hmacKey)
	fmt.Fprintln(hm, username, int(math.Ceil(float64(time.Now().Unix())/60))-age)
	return fmt.Sprintf("%x", hm.Sum(nil))[2:8]
}

func requireMethod(method string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return errInvalidMethod
	}
	return nil
}

func requireAPIKey(session *sessions.Session, w http.ResponseWriter) error {
	var username string
	var apiKey string

	n, ok := session.Values["username"]
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return errInvalidAuth
	}
	username = n.(string)

	a, ok := session.Values["apiKey"]
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return errInvalidAuth
	}
	apiKey = a.(string)

	for i := 0; i < 10; i++ {
		if apiKey == generateAPIKeyForUserTime(username, i) {
			return nil
		}
	}

	w.WriteHeader(http.StatusForbidden)
	return errInvalidAuth
}

func checkAuthorization(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "raidbot")
	session.Options.MaxAge = 604800
	h := r.Header.Get("Authorization")
	if h == "" {
		return session
	}
	p := strings.Split(h, " ")
	if len(p) < 2 {
		return session
	}
	if p[0] != "fof-ut" {
		return session
	}
	t := strings.Split(p[1], ":")
	if len(t) < 3 {
		return session
	}
	mac := hmac.New(sha256.New, hmacKey)
	fmt.Fprintln(mac, t[0], t[1])
	key := fmt.Sprintf("%x", mac.Sum(nil))
	if t[2] == key {
		session.Values["username"] = t[0]
		session.Values["apiKey"] = generateAPIKeyForUserTime(t[0], 0)
	}
	return session
}

func doRESTRouter(w http.ResponseWriter, r *http.Request) {
	session := checkAuthorization(w, r)
	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form values for", r.Method, r.RequestURI, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uri, _ := url.ParseRequestURI(r.RequestURI)
	v, err := url.ParseQuery(uri.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	switch uri.Path {
	case "/rest/login":
		username := v.Get("username")
		t, err := strconv.ParseInt(v.Get("t"), 10, 64)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("//%s/", r.Host), http.StatusFound)
			return
		}
		s := v.Get("signature")
		now := time.Now().Unix()
		if t > now {
			http.Redirect(w, r, fmt.Sprintf("//%s/", r.Host), http.StatusFound)
			return
		}
		if (now - t) > 300 {
			http.Redirect(w, r, fmt.Sprintf("//%s/", r.Host), http.StatusFound)
			return
		}
		mac := hmac.New(sha256.New, hmacKey)
		fmt.Fprintln(mac, username, t)
		if s != fmt.Sprintf("%x", mac.Sum(nil)) {
			http.Redirect(w, r, fmt.Sprintf("//%s/", r.Host), http.StatusFound)
			return
		}
		log.Printf("@%s -- %s", username, "/rest/login")
		session.Values["username"] = username
		session.Values["apiKey"] = generateAPIKeyForUserTime(username, 0)
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("//%s/", r.Host), http.StatusFound)
	case "/rest/login/logout":
		delete(session.Values, "apiKey")
		delete(session.Values, "username")
		session.Save(r, w)
		return
	case "/rest/login/check":
		if name, ok := session.Values["username"]; ok {
			session.Values["apiKey"] = generateAPIKeyForUserTime(name.(string), 0)
			session.Save(r, w)
			data, _ := json.Marshal(map[string]string{"cmd": raidSlashCommand, "username": name.(string)})
			w.Write(data)
			log.Printf("@%s -- /rest/login/check", name.(string))
		} else {
			data, _ := json.Marshal(map[string]string{"cmd": raidSlashCommand})
			w.Write(data)
			log.Println("- -- /rest/login/check")
			return
		}
	case "/rest/lfg":
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			expiry, err := strconv.Atoi(r.Form.Get("time"))
			if err != nil {
				expiry = 120
			} else {
				if expiry < 30 {
					expiry = 30
				} else {
					if expiry > 120 {
						expiry = 120
					}
				}
			}

			if events, ok := r.Form["events[]"]; ok {
				eventList, _ := url.QueryUnescape(strings.Join(events, "', '"))
				lfg.add(username, time.Duration(expiry)*time.Minute, events...)
				log.Printf("@%s -- lfg: %s: '%s'", username, r.Form.Get("time"), eventList)
			} else {
				lfg.add(username, 0)
				log.Printf("@%s -- lfg: -", username)
			}
			return
		}

		var since = v.Get("since")
		if since == "0" || since == "" {
			if err := lfgOutput.send(w); err != nil {
				log.Println("Error sending /rest/lfg:", err.Error())
			}
			return
		}

		closed := w.(http.CloseNotifier).CloseNotify()
		notify := make(chan struct{})
		go func(notify chan struct{}, since string) {
			for lfgOutput.updatedAt == since {
				lfgOutput.cond.Wait()
			}
			notify <- struct{}{}
		}(notify, since)
		select {
		case <-notify:
			if err := lfgOutput.send(w); err != nil {
				log.Println("Error sending /rest/lfg:", err.Error())
			}
		case <-closed:
		}
		return

	case "/rest/get":
		if err := requireAPIKey(session, w); err != nil {
			return
		}

		var since = v.Get("since")
		if since == "0" || since == "" {
			if err := xhrOutput.send(w); err != nil {
				log.Println("Error sending /rest/lfg:", err.Error())
			}
			return
		}

		closed := w.(http.CloseNotifier).CloseNotify()
		notify := make(chan struct{})
		go func(notify chan struct{}, since string) {
			for xhrOutput.updatedAt == since {
				xhrOutput.cond.Wait()
			}
			notify <- struct{}{}
		}(notify, since)
		select {
		case <-notify:
			if err := xhrOutput.send(w); err != nil {
				log.Println("Error sending /rest/get:", err.Error())
			}
		case <-closed:
		}
		return
	case "/rest/raid/join":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		channel := r.Form.Get("channel")
		raid := r.Form.Get("raid")
		log.Printf("@%s on %s -- %s %s", username, channel, "join", raid)
		msgs, err := raidJoin(username, channel, raid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, msgs.stdOut())
			return
		}
		msgs.sendToSlack()
		fmt.Fprint(w, msgs.stdOut())
	case "/rest/raid/leave":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		channel := r.Form.Get("channel")
		raid := r.Form.Get("raid")
		log.Printf("@%s on %s -- %s %s", username, channel, "leave", raid)
		msgs, err := raidLeave(username, channel, raid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, msgs.stdOut())
			return
		}
		msgs.sendToSlack()
		fmt.Fprint(w, msgs.stdOut())
	case "/rest/raid/join-alt":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		channel := r.Form.Get("channel")
		raid := r.Form.Get("raid")
		log.Printf("@%s on %s -- %s %s", username, channel, "join-alt", raid)
		msgs, err := raidAltJoin(username, channel, raid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, msgs.stdOut())
			return
		}
		msgs.sendToSlack()
		fmt.Fprint(w, msgs.stdOut())
	case "/rest/raid/leave-alt":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		channel := r.Form.Get("channel")
		raid := r.Form.Get("raid")
		log.Printf("@%s on %s -- %s %s", username, channel, "leave-alt", raid)
		msgs, err := raidAltLeave(username, channel, raid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, msgs.stdOut())
			return
		}
		msgs.sendToSlack()
		fmt.Fprint(w, msgs.stdOut())
	case "/rest/raid/finish":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		channel := r.Form.Get("channel")
		raid := r.Form.Get("raid")
		log.Printf("@%s on %s -- %s %s", username, channel, "finish", raid)
		msgs, err := raidFinish(username, channel, raid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, msgs.stdOut())
			return
		}
		msgs.sendToSlack()
		fmt.Fprint(w, msgs.stdOut())
	case "/rest/raid/host":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		channel := r.Form.Get("channel")
		raid := r.Form.Get("raid")
		raidTitle := r.Form.Get("raidName")

		//Get Raid time
		raidTime := r.Form.Get("time")
		timeInt, _ := strconv.ParseInt(raidTime, 10, 64)
		timeVal := time.Unix((timeInt / 1000), 0)

		log.Printf("@%s on %s -- %s %s", username, channel, "host", raid)
		msgs, err := raidHost(username, channel, raid, timeVal, raidTitle)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, msgs.stdOut())
			return
		}
		msgs.sendToSlack()
		fmt.Fprint(w, msgs.stdOut())
	case "/rest/raid/ping":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		username, _ := session.Values["username"].(string)
		channel := r.Form.Get("channel")
		raid := r.Form.Get("raid")
		log.Printf("@%s on %s -- %s %s", username, channel, "ping", raid)
		msgs, err := raidPing(username, channel, raid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, msgs.stdOut())
			return
		}
		msgs.sendToSlack()
		fmt.Fprint(w, msgs.stdOut())
	case "/rest/ping":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		to := r.Form.Get("username")
		about := r.Form.Get("about")
		username, _ := session.Values["username"].(string)
		log.Printf("@%s -- ping @%s -- %s", username, to, about)
		slack.msg().to("@" + to).send(fmt.Sprintf(
			"@%s is also looking to play *%s*", username, about,
		))
	case "/rest/report":
		if err := requireMethod("POST", w, r); err != nil {
			return
		}
		if err := requireAPIKey(session, w); err != nil {
			return
		}
		about := r.Form.Get("about")
		username, _ := session.Values["username"].(string)
		message := r.Form.Get("message")
		log.Printf("@%s -- report @%s -- %s", username, about, message)

		slackMessage := fmt.Sprintf(
			"`report-a-member: @%s reported by @%s`\n> %s",
			about,
			username,
			strings.Replace(message, "\n", "\n> ", -1),
		)
		for _, admin := range admins {
			slack.msg().to("@" + admin).send(fmt.Sprintf(slackMessage))
		}
	}
}
