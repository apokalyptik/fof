package main

import (
	"fmt"
	"log"

	"github.com/influxdata/influxdb/client/v2"
)

var infxMessage = make(chan slackmsg, 1024)

func mindInfluxDB() {
	go func() {
		var c client.Client
		defer func(c client.Client) { c.Close() }(c)
		for {
			if cl, err := client.NewHTTPClient(client.HTTPConfig{Addr: infxAddr}); err != nil {
				log.Fatalf("Error creating influxdb client: %s", err.Error())
			} else {
				c = cl
			}
			bpCount := 0
			bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "fof",
				Precision: "s",
			})

			for {
				select {
				case message := <-infxMessage:
					pt, err := client.NewPoint(
						"public_mesages",
						map[string]string{
							"user_name":    message.User,
							"user_id":      message.UserID,
							"channel_name": message.Channel,
							"channel_id":   message.ChannelID,
						},
						map[string]interface{}{
							"message": message.Text,
						},
						message.time(),
					)
					if err != nil {
						log.Printf("error in client.NewPoint: %s", err.Error())
						break
					}
					bp.AddPoint(pt)
					bpCount++
					fmt.Println(pt.String())
					if bpCount >= 1 {
						if err := c.Write(bp); err != nil {
							log.Printf("error in client.Client.Write: %s", err.Error())
							break
						}
						bpCount = 0
						bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
							Database:  "fof",
							Precision: "s",
						})
					}
				}
			}
		}
	}()
}
