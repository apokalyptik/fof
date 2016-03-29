package main

import (
	"log"
	"strconv"
	"time"

	"github.com/apokalyptik/fof/lib/fof/stats"
	"github.com/nsqio/go-nsq"
)

var nsqMessage = make(chan slackmsg, 1024)

func mindNsq() {
	if conn, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig()); err != nil {
		log.Fatal("Unable to create NSQ Producer: %s", err.Error())
	} else {
		if err := conn.Ping(); err != nil {
			log.Fatalf("Error connecting to NSQ: %s", err.Error())
		}
		go func(conn *nsq.Producer) {
			for {
				select {
				case m := <-nsqMessage:
					unix, _ := strconv.ParseFloat(m.Timestamp, 64)
					fstat.Stat{
						Member:   m.UserID,
						Platform: "slack",
						Product:  "slack",
						Stat:     "messages",
						Sub1:     "user",
						Sub2:     m.Channel,
						When:     time.Unix(int64(unix), 0),
						Value:    1,
						Method:   "inc",
					}.Send(conn)
					fstat.Stat{
						Member:   m.UserID,
						Platform: "slack",
						Product:  "slack",
						Stat:     "messages",
						Sub1:     "user",
						When:     time.Unix(int64(unix), 0),
						Value:    1,
						Method:   "inc",
					}.Send(conn)
				}
			}
		}(conn)
	}
}
