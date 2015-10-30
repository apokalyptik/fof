package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
)

var (
	templateURL    = "https://slack.com/api/files.list?token=%s&page=%d"
	templateDelete = "https://slack.com/api/files.delete?token=%s&file=%s"
	apiKey         = ""
	maxAge         = 672 * time.Hour
)

type fileListResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	Files []struct {
		Name      string `json:"name"`
		MimeType  string `json:"mimetype"`
		ID        string `json:"id"`
		Created   int64  `json:"created"`
		Timestamp int64  `json:"timestamp"`
		Size      uint64 `json:"size"`
	} `json:"files"`
}

func init() {
	flag.StringVar(&apiKey, "api", apiKey, "slack API Key")
	flag.DurationVar(&maxAge, "max", maxAge, "maximum age of files to keep")
}

func main() {
	flag.Parse()
	newest := time.Now().Add(0 - maxAge).Unix()
	var kept uint64
	var pruned uint64
	var unknown uint64
	for i := 0; i < 1000; i++ {
		if rsp, err := http.Get(fmt.Sprintf(templateURL, apiKey, i)); err != nil {
			log.Fatal(err)
		} else {
			var data fileListResponse
			var dec = json.NewDecoder(rsp.Body)
			if err := dec.Decode(&data); err != nil {
				log.Fatal(err)
			}
			rsp.Body.Close()
			if !data.OK {
				log.Fatalf("Error fetching list: %s", data.Error)
			}
			if data.Files == nil {
				break
			}
			if len(data.Files) == 0 {
				break
			}
			for _, file := range data.Files {
				if file.Created >= newest {
					kept = kept + file.Size
					log.Println(
						"skipping",
						humanize.Bytes(file.Size),
						file.MimeType,
						time.Now().Sub(time.Unix(file.Timestamp, 0)).String(),
						fmt.Sprintf("'%s'", file.Name))
					continue
				}
				if d, err := http.Get(fmt.Sprintf(templateDelete, apiKey, file.ID)); err != nil {
					unknown = unknown + file.Size
					log.Println(
						"delete failure",
						humanize.Bytes(file.Size),
						file.MimeType,
						time.Now().Sub(time.Unix(file.Timestamp, 0)).String(),
						fmt.Sprintf("'%s'", file.Name),
						fmt.Sprintf("'%s'", err.Error()),
					)
				} else {
					pruned = pruned + file.Size
					log.Println(
						"delete success",
						humanize.Bytes(file.Size),
						file.MimeType,
						time.Now().Sub(time.Unix(file.Timestamp, 0)).String(),
						fmt.Sprintf("'%s'", file.Name))
					d.Body.Close()
				}
			}
		}
	}
	log.Println("kept", humanize.Bytes(kept), "deleted", humanize.Bytes(pruned), "unknown", humanize.Bytes(unknown))
}
