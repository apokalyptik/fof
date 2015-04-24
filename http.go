package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func doHTTPStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

func doHTTP404(w http.ResponseWriter) {
	doHTTPStatus(w, http.StatusNotFound)
}

func doHTTPRouter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Form.Get("token") != slack.key {
		doHTTPStatus(w, http.StatusUnauthorized)
		log.Printf("Unauthorized Request: %#v -- %#v", http.StatusUnauthorized, r.Form)
		return
	}
	command := strings.Split(r.Form.Get("text"), " ")
	subcommand := strings.Join(command[1:], " ")
	channel := r.Form.Get("channel_name")
	username := r.Form.Get("user_name")

	log.Printf("@%s on #%s -- %s %s", username, channel, r.Form.Get("command"), r.Form.Get("text"))

	switch r.Form.Get("command") {
	case "/raid":
		switch channel {
		case "privategroup":
			fmt.Fprint(w, "Please use the /raid command in a public channel")
			return
		case "directmessage":
			fmt.Fprint(w, "Please use the /raid command in a channel")
			return
		}
		if len(command) < 1 || command[0] == "help" || command[0] == "" {
			fmt.Fprint(w, "/raid [command]\n\nWhere command is one of:\n\n")
			fmt.Fprint(w, "\t• list\n")
			fmt.Fprint(w, "\t• host [name of a raid to create]\n")
			fmt.Fprint(w, "\t• join [name of a raid to sign up for]\n")
			fmt.Fprint(w, "\t• leave [name of a raid to remove yourself from]\n")
			fmt.Fprint(w, "\t• finish [name of a raid to remove]\n")
			fmt.Fprint(w, "\t• ping [name of a raid to ping people for]\n\n")
			fmt.Fprint(w, "This will only be for raids in #"+channel+". ")
			fmt.Fprint(w, "To find and use raids in other channels you'll want to ")
			fmt.Fprint(w, "use the /raid command from those channels")
			fmt.Fprint(w, "\n\nFor an introduction, please watch https://www.youtube.com/watch?v=T4g_3Tv5xJU")
			return
		}

		switch strings.ToLower(command[0]) {
		case "list":
			list := raidDb.list(channel)
			if len(list) == 0 {
				fmt.Fprintf(
					w,
					"There are no raids being hosted on #%s yet. Perhaps "+
						"you would like to \"/raid host\" one?",
					channel)
			} else {
				fmt.Fprintf(w, "The following raids are being hosted on #%s:\n", channel)
				for _, v := range raidDb.list(channel) {
					fmt.Fprintf(
						w,
						"• \"%s\" with: _%s_ <http://%s%s?a=rj&u=%s&c=%s&r=%s&h=%s|join>\n",
						v.Name,
						strings.Join(v.Members, "_, _"),
						r.Host,
						r.RequestURI,
						username,
						channel,
						v.UUID,
						v.hmacForUser(username),
					)
				}
			}
		case "host":
			if len(subcommand) < 3 {
				fmt.Fprintf(
					w,
					"Sorry you must give me a little bit more to work with than \"%s\"",
					subcommand)
				return
			}

			if err := raidDb.register(channel, subcommand, username); err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, fmt.Sprintf(
					"OK. \"%s\" has been registered on #%s for you",
					subcommand,
					channel))
				slack.toChannel(channel, fmt.Sprintf(
					"@%s is hosting a new raid: \"%s\".\nUse \"/raid join %s\" to sign up!",
					username, subcommand, subcommand))
			}
		case "join":
			if err := raidDb.join(channel, subcommand, username); err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, fmt.Sprintf(
					"OK. You're signed up for \"%s\" on #%s",
					subcommand,
					channel))
				slack.toChannel(channel, fmt.Sprintf(
					"@%s has signed up for \"%s\".\nUse \"/raid join %s\" to join them!",
					username, subcommand, subcommand))
			}
		case "leave":
			if err := raidDb.leave(channel, subcommand, username); err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, fmt.Sprintf(
					"OK. You're no longer signed up for \"%s\" on #%s",
					subcommand,
					channel))
				slack.toChannel(channel, fmt.Sprintf(
					"@%s is no longer signed up for \"%s\".\nUse \"/raid join %s\" to take their place!",
					username, subcommand, subcommand))
			}
		case "finish":
			if err := raidDb.finish(channel, subcommand, username); err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, fmt.Sprintf(
					"OK. \"%s\" has been removed from the raid list for #%s",
					subcommand,
					channel))
				slack.toChannel(channel, fmt.Sprintf(
					"@%s has closed out \"%s\"",
					username, subcommand))
			}
		case "ping":
			if list, err := raidDb.members(channel, subcommand); err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				slack.toChannel(channel, fmt.Sprintf(
					"pinging @%s about \"%s\" for @%s",
					strings.Join(list, ", @"), subcommand, username))
			}
		default:
			fmt.Fprint(w, "I'm afraid I don't know how to '"+command[0]+"'. ")
			fmt.Fprint(w, "Try '/raid help' to get a list of things I can do for you")
			return
		}
	default:
		doHTTPStatus(w, http.StatusNotImplemented)
		return
	}
}
