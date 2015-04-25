package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	switch r.Method {
	case "POST":
		doHTTPPost(w, r)
	default:
		doHTTPGet(w, r)
	}
}

func raidLeave(w http.ResponseWriter, channel, raid, username, cmd string) {
	if err := raidDb.leave(channel, raid, username); err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprint(w, fmt.Sprintf(
			"OK. You're no longer signed up for \"%s\" on #%s",
			raid,
			channel))
		slack.toChannel(channel, fmt.Sprintf(
			"@%s is no longer signed up for \"%s\".\nUse \"%s join %s\" to take their place!",
			username, raid, cmd, raid))
	}
}

func raidJoin(w http.ResponseWriter, user, channel, raid, cmd string) error {
	if err := raidDb.join(channel, raid, user); err != nil {
		fmt.Fprint(w, err.Error())
		return err
	}
	fmt.Fprint(w, fmt.Sprintf("OK. You're signed up for \"%s\" on #%s", raid, channel))
	slack.toChannel(channel, fmt.Sprintf(
		"@%s has signed up for \"%s\".\nUse \"%s join %s\" to join them!",
		user, raid, cmd, raid))
	return nil
}

func doHTTPGet(w http.ResponseWriter, r *http.Request) {
	var q = r.URL.Query()
	switch q.Get("a") {
	case "rj": // Raid Join
		channel := q.Get("c")
		uuid := q.Get("r")
		user := q.Get("u")
		sig := q.Get("h")
		if channel == "" || user == "" || uuid == "" || sig == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, r := range raidDb.list(channel) {
			if r.UUID == uuid {
				if r.validateHmacForUser(user, sig) == nil {
					raidJoin(w, user, channel, r.Name, q.Get("cmd"))
					return
				}
			}
		}
	case "rp": // Raid Part
		channel := q.Get("c")
		uuid := q.Get("r")
		user := q.Get("u")
		sig := q.Get("h")
		if channel == "" || user == "" || uuid == "" || sig == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, r := range raidDb.list(channel) {
			if r.UUID == uuid {
				if r.validateHmacForUser(user, sig) == nil {
					raidLeave(w, channel, r.Name, user, q.Get("cmd"))
					return
				}
			}
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func doHTTPPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := strings.Split(r.Form.Get("text"), " ")
	subcommand := strings.Join(command[1:], " ")
	channel := r.Form.Get("channel_name")
	username := r.Form.Get("user_name")

	log.Printf("@%s on #%s -- %s %s", username, channel, r.Form.Get("command"), r.Form.Get("text"))

	switch r.Form.Get("command") {
	case "/xline":
		if r.Form.Get("token") != slack.xlineKey {
			doHTTPStatus(w, http.StatusUnauthorized)
			log.Printf("Unauthorized Request: %#v -- %#v", http.StatusUnauthorized, r.Form)
			return
		}
		switch channel {
		case "privategroup":
			fmt.Fprint(w, "Please use the "+r.Form.Get("command")+" command in a public channel")
			return
		case "directmessage":
			fmt.Fprint(w, "Please use the "+r.Form.Get("command")+" command in a channel")
			return
		}
		if len(command) < 1 || command[0] == "help" || command[0] == "" {
			for _, v := range admins {
				if v != username {
					continue
				}
				fmt.Fprint(w, r.Form.Get("command")+" add <sticker> <url> <description>\n")
				fmt.Fprint(w, r.Form.Get("command")+" remove <sticker>\n")
				break
			}
			fmt.Fprint(w, r.Form.Get("command")+" search <term>\n")
			fmt.Fprint(w, r.Form.Get("command")+" <sticker>\n")
			return
		}
		switch strings.ToLower(command[0]) {
		case "add":
			var isAdmin = false
			for _, v := range admins {
				if v != username {
					continue
				}
				isAdmin = true
			}
			if !isAdmin {
				fmt.Fprint(w, "Only an admin may use this feature")
				return
			}
			var err error
			switch len(command) {
			case 0, 1, 2:
				err = errors.New("You must provide a name and a url at least")
			case 3:
				err = xlineDB.add(strings.ToLower(command[1]), command[2], "")
			default:
				err = xlineDB.add(strings.ToLower(command[1]), command[2], strings.ToLower(strings.Join(command[3:], " ")))
			}
			if err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, "added")
			}
		case "remove":
			var isAdmin = false
			for _, v := range admins {
				if v != username {
					continue
				}
				isAdmin = true
			}
			if !isAdmin {
				fmt.Fprint(w, "Only an admin may use this feature")
				return
			}
			fmt.Fprint(w, "remove")
			if err := xlineDB.remove(strings.ToLower(command[1])); err != nil {
				fmt.Fprint(w, err.Error())
			} else {
				fmt.Fprint(w, "removed \""+command[1]+"\"")
			}
		case "search":
			found := xlineDB.search(strings.ToLower(subcommand))
			if len(found) > 0 {
				fmt.Fprintf(w, "found %d stickers:\n\t%s", len(found), strings.Join(found, "\n\t"))
			} else {
				fmt.Fprintf(w, "found 0 stickers")
			}
		default:
			if out, err := xlineDB.get(strings.ToLower(command[0])); err != nil {
				found := xlineDB.search(strings.ToLower(strings.Join(command, " ")))
				switch len(found) {
				case 0:
					fmt.Fprint(w, err.Error())
				case 1:
					out, _ = xlineDB.get(strings.Split(found[0], " ")[0])
					slack.toChannel(
						channel,
						fmt.Sprintf("%s: %s %s\n%s", username, r.Form.Get("command"), r.Form.Get("text"), out),
						"stickybot")
				default:
					fmt.Fprintf(w, "did you mean one of these?\n\t%s", strings.Join(found, "\n\t"))
				}
			} else {
				slack.toChannel(
					channel,
					fmt.Sprintf("%s: %s %s\n%s", username, r.Form.Get("command"), r.Form.Get("text"), out),
					"stickybot")
			}
		}
	case "/needs":
		if r.Form.Get("token") != slack.needKey {
			doHTTPStatus(w, http.StatusUnauthorized)
			log.Printf("Unauthorized Request: %#v -- %#v", http.StatusUnauthorized, r.Form)
			return
		}
	case "/raid", "/raidtest":
		if r.Form.Get("token") != slack.raidKey {
			doHTTPStatus(w, http.StatusUnauthorized)
			log.Printf("Unauthorized Request: %#v -- %#v", http.StatusUnauthorized, r.Form)
			return
		}
		switch channel {
		case "privategroup":
			fmt.Fprint(w, "Please use the "+r.Form.Get("command")+" command in a public channel")
			return
		case "directmessage":
			fmt.Fprint(w, "Please use the "+r.Form.Get("command")+" command in a channel")
			return
		}
		if len(command) < 1 || command[0] == "help" || command[0] == "" {
			fmt.Fprint(w, r.Form.Get("command")+" [command]\n\nWhere command is one of:\n\n")
			fmt.Fprint(w, "\t• list\n")
			fmt.Fprint(w, "\t• host [name of a raid to create]\n")
			fmt.Fprint(w, "\t• join [name of a raid to sign up for]\n")
			fmt.Fprint(w, "\t• leave [name of a raid to remove yourself from]\n")
			fmt.Fprint(w, "\t• finish [name of a raid to remove]\n")
			fmt.Fprint(w, "\t• ping [name of a raid to ping people for]\n\n")
			fmt.Fprint(w, "This will only be for raids in #"+channel+". ")
			fmt.Fprint(w, "To find and use raids in other channels you'll want to ")
			fmt.Fprint(w, "use the "+r.Form.Get("command")+" command from those channels")
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
						"you would like to \""+r.Form.Get("command")+" host\" one?",
					channel)
			} else {
				fmt.Fprintf(w, "The following raids are being hosted on #%s:\n", channel)
				for _, v := range raidDb.list(channel) {
					var isMember = false
					for _, m := range v.Members {
						if m == username {
							isMember = true
							break
						}
					}
					var link = ""
					if isMember == false {
						link = fmt.Sprintf(
							"<http://%s%s?a=rj&u=%s&c=%s&r=%s&h=%s&cmd=%s|join>",
							r.Host,
							r.RequestURI,
							url.QueryEscape(username),
							url.QueryEscape(channel),
							url.QueryEscape(v.UUID),
							url.QueryEscape(v.hmacForUser(username)),
							url.QueryEscape(r.Form.Get("command")),
						)
					} else {
						link = fmt.Sprintf(
							"<http://%s%s?a=rp&u=%s&c=%s&r=%s&h=%s&cmd=%s|leave>",
							r.Host,
							r.RequestURI,
							url.QueryEscape(username),
							url.QueryEscape(channel),
							url.QueryEscape(v.UUID),
							url.QueryEscape(v.hmacForUser(username)),
							url.QueryEscape(r.Form.Get("command")),
						)
					}
					fmt.Fprintf(
						w,
						"• \"%s\" with: _%s_ %s\n",
						v.Name,
						strings.Join(v.Members, "_, _"),
						link,
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
					"@%s is hosting a new raid: \"%s\".\nUse \""+r.Form.Get("command")+" join %s\" to sign up!",
					username, subcommand, subcommand))
			}
		case "join":
			raidJoin(w, username, channel, subcommand, r.Form.Get("command"))
		case "leave":
			raidLeave(w, channel, subcommand, username, r.Form.Get("command"))
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
			fmt.Fprint(w, "Try '/"+r.Form.Get("command")+" help' to get a list of things I can do for you")
			return
		}
	default:
		doHTTPStatus(w, http.StatusNotImplemented)
		return
	}
}
