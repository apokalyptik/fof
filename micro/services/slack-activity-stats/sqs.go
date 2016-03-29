package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type slackmsg struct {
	Token     string `json:"token"`
	Team      string `json:"team_domain"`
	TeamID    string `json:"team_id"`
	Channel   string `json:"channel_name"`
	ChannelID string `json:"channel_id"`
	User      string `json:"user_name"`
	UserID    string `json:"user_id"`
	Text      string `json:"text"`
	Timestamp string `json:"Timestamp"`
}

func (s slackmsg) time() time.Time {
	if f, err := strconv.ParseFloat(s.Timestamp, 64); err != nil {
		log.Println("error in strconv.ParseInt:", err.Error())
		return time.Unix(0, 0)
	} else {
		return time.Unix(int64(f), 0)
	}
}

func mindSQS() {
	for {
		svc := sqs.New(session.New(), &aws.Config{Region: aws.String(SQSRegion)})
		params := &sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(SQSURL),
			AttributeNames:        []*string{aws.String("All")},
			MaxNumberOfMessages:   aws.Int64(10),
			MessageAttributeNames: []*string{aws.String("All")},
			VisibilityTimeout:     aws.Int64(1),
			WaitTimeSeconds:       aws.Int64(10),
		}
		for {
			if resp, err := svc.ReceiveMessage(params); err != nil {
				log.Println("error in svc.ReceiveMessage():", err.Error())
				break
			} else {
				for _, message := range resp.Messages {
					var m slackmsg
					if err := json.Unmarshal([]byte(*message.Body), &m); err != nil {
						log.Println("error unmarshalling json:", err.Error())
					}
					log.Printf("Recieved message: %s/%s/%s/%s", m.Team, m.Channel, m.User, m.Timestamp)
					infxMessage <- m
					sqlMessage <- m
					nsqMessage <- m
					_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
						QueueUrl:      aws.String(SQSURL),
						ReceiptHandle: aws.String(*message.ReceiptHandle),
					})
					if err != nil {
						log.Println("error deleting message:", err.Error())
					}
				}
			}
		}
	}
}
