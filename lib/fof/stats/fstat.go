package fstat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

var Topic = "fof-stats"

type Stat struct {
	Member   string    `json:"member"`
	Platform string    `json:"platform"`
	Product  string    `json:"product"`
	Stat     string    `json:"stat"`
	Sub1     string    `json:"sub1"`
	Sub2     string    `json:"sub2"`
	Sub3     string    `json:"sub3"`
	Info     string    `json:"info"`
	When     time.Time `json:"When"`
	Value    int       `json:"value"`
	Method   string    `json:"method"`
}

func (s Stat) Send(conn *nsq.Producer) error {
	payload, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("Error marshalling payload: %s", err.Error())
	}
	return conn.Publish(Topic, payload)
}
