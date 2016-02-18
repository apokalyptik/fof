package main

import (
	"fmt"
	"log"

	"github.com/influxdata/influxdb/client/v2"
)

var infxMessage = make(chan slackmsg, 1024)

func mindInfluxDB() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: infxAddr,
	})
	if err != nil {
		log.Fatalf("Error creating influxdb client: %s", err.Error())
	}
	go func(c client.Client) {
		defer c.Close()
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
					log.Fatalf("error in client.NewPoint: %s", err.Error())
				}
				bp.AddPoint(pt)
				bpCount++
				fmt.Println(pt.String())
				if bpCount >= 1 {
					if err := c.Write(bp); err != nil {
						log.Fatalf("error in client.Client.Write: %s", err.Error())
					}
					bpCount = 0
					bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
						Database:  "fof",
						Precision: "s",
					})
				}
			}
		}
	}(c)
}
