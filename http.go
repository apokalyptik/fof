package main

import (
	"fmt"
	"net/http"
	"strings"
)

func HTTPStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

func HTTP404(w http.ResponseWriter) {
	HTTPStatus(w, http.StatusNotFound)
}

func HTTPRouter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Form.Get("token") != slack.key {
		HTTPStatus(w, http.StatusUnauthorized)
		return
	}
	switch r.Form.Get("command") {
	case "/raid":
		channel := r.Form.Get("channel_name")
		switch channel {
		case "privategroup":
			fmt.Fprint(w, "Please use the /raid command in a public channel")
			return
		case "directmessage":
			fmt.Fprint(w, "Please use the /raid command in a channel")
			return
		}
		command := strings.Split(r.Form.Get("text"), " ")
		if len(command) < 1 || command[0] == "help" || command[0] == "" {
			fmt.Fprint(w, "/raid [command]\n\nWhere command is one of:\n\n")
			fmt.Fprint(w, "\t• list\n")
			fmt.Fprint(w, "\t• make [name of a raid to create]\n")
			fmt.Fprint(w, "\t• join [name of a raid to sign up for]\n")
			fmt.Fprint(w, "\t• part [name of a raid to remove yourself from]\n")
			fmt.Fprint(w, "\t• info [name of a raid to get info about]\n")
			fmt.Fprint(w, "\t• done [name of a raid to remove]\n")
			fmt.Fprint(w, "\t• ping [name of a raid to ping people for]\n\n")
			fmt.Fprint(w, "This will only be for raids in #"+channel+". ")
			fmt.Fprint(w, "To find and use raids in other channels you'll want to ")
			fmt.Fprint(w, "use the /raid command from those channels")
			return
		}
		//username := r.Form.Get("user_name")
		switch strings.ToLower(command[0]) {
		case "list":
		case "make":
		case "join":
		case "part":
		case "info":
		case "done":
		case "ping":
			fmt.Fprint(w, "ok")
		default:
			fmt.Fprint(w, "I'm afraid I don't know how to '"+command[0]+"'. ")
			fmt.Fprint(w, "Try '/raid help' to get a list of things I can do for you")
			return
		}
	default:
		HTTPStatus(w, http.StatusNotImplemented)
		return
	}
}
