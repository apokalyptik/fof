package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"

	"net/http"
	"net/url"
	"strings"
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
	text := r.Form.Get("text")

	if text != "" {
		firstSpace := strings.Index(text, " ")
		cmd := text
		if firstSpace > 0 {
			cmd = text[:firstSpace]
		}
		switch cmd {
		case "xbox":
			if len(text) > len("xbox") {
				fmt.Fprint(w, getXboxProfileURL(text[firstSpace+1:]))
			} else {
				fmt.Fprint(w, "You need to specify a gamer tag")
			}
		}
	} else {
		mac := hmac.New(sha256.New, hmacKey)
		t := time.Now().Unix()
		fmt.Fprintln(mac, username, t)
		go func() {
			slack.msg().to("@" + username).send(fmt.Sprintf(
				"<http://%s/rest/login?username=%s&t=%d&signature=%s|Click here> to log into and use the team tool. The link is only valid for the next 5 or so minutes. You can request a new one at any time with /team. You will be logged out of the team tool after about a week, and will need to log in again when that happens.",
				r.Host,
				url.QueryEscape(username),
				t,
				fmt.Sprintf("%x", mac.Sum(nil)),
			))
		}()
		fmt.Fprint(w, "You have requested access to the FoF Team site. You will receive a direct message from FOFBOT with a link to the site.")
		log.Printf(
			"@%s on %s -- %s %s",
			username,
			r.Form.Get("channel_name"),
			r.Form.Get("command"),
			r.Form.Get("text"))
	}

}
