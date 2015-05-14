package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"

	"net/http"
	"net/url"
	"time"
)

func doHTTPStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

func doHTTP404(w http.ResponseWriter) {
	doHTTPStatus(w, http.StatusNotFound)
}

func doHTTPPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("user_name")
	mac := hmac.New(sha256.New, hmacKey)
	t := time.Now().Unix()
	fmt.Fprintln(mac, username, t)
	slack.msg().to("@" + username).send(fmt.Sprintf(
		"<http://%s/rest/login?username=%s&t=%d&signature=%s|Click here> to log into and use the team tool",
		r.Host,
		url.QueryEscape(username),
		t,
		fmt.Sprintf("%x", mac.Sum(nil)),
	))
	fmt.Fprint(w, "You have requested access to the FoF Team site. You will receive a direct message from SLACKBOT with a link to the site.")
	log.Printf(
		"@%s on %s -- %s %s",
		username,
		r.Form.Get("channel_name"),
		r.Form.Get("command"),
		r.Form.Get("text"))
}
