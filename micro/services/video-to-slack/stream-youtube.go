package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	youtubeConfigFile          = "./youtube.json"
	youtubeCFG                 = youtubeConfig{}
	youtubeAPIKey              = ""
	youtubeBaseURL             = "https://www.googleapis.com/youtube/v3"
	youtubeLiveStreamSearchURL = "%s/search?part=snippet&eventType=live&type=video&channelId=%s&key=%s"
	youtubeLinkURL             = "https://gaming.youtube.com/watch?v=%s"
	youtubeStreams             = make(map[string][]string)
)

type youtubeConfig struct {
	Key      string   `json:"api_key"`
	Channels []string `json:"channels"`
}

type youtubeChannelsResponse struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

type youtubeLiveStreamsResponseItem struct {
	ID struct {
		VideoId string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		Title                string `json:"title"`
		Description          string `json:"description"`
		ChannelTitle         string `json:"channelTitle"`
		LiveBroadcastContent string `json:"liveBroadcastContent"`
		Thumbnails           struct {
			Default struct {
				URL string `json:"url"`
			} `json:"default"`
		} `json:"thumbnails"`
	} `json:"snippet"`
}

type youtubeLiveStreamsResponse struct {
	Items []youtubeLiveStreamsResponseItem `json:"items"`
}

func init() {
	flag.StringVar(&youtubeConfigFile, "youtube", youtubeConfigFile, "path to the youtube config file")
}

func doYoutubeMessage(liveStreamItem youtubeLiveStreamsResponseItem) {
	messageParams := slack.NewPostMessageParameters()
	messageParams.AsUser = true
	messageParams.Parse = "full"
	messageParams.LinkNames = 1
	messageParams.UnfurlMedia = true
	messageParams.UnfurlLinks = true
	messageParams.EscapeText = false
	var channelTitle = liveStreamItem.Snippet.ChannelTitle
	if channelTitle == "" {
		channelTitle = "Someone Forgot To Set A Channel Title"
	}
	messageParams.Attachments = append(messageParams.Attachments, slack.Attachment{
		Title: fmt.Sprintf(
			"Watch %s play %s",
			channelTitle,
			liveStreamItem.Snippet.Title,
		),
		TitleLink: fmt.Sprintf(youtubeLinkURL, liveStreamItem.ID.VideoId),
		ThumbURL:  liveStreamItem.Snippet.Thumbnails.Default.URL,
	})
	_, _, err := slackClient.PostMessage(
		slackCFG.Channel,
		fmt.Sprintf(
			"*%s* has begun streaming *%s*\n\n%s",
			channelTitle,
			liveStreamItem.Snippet.Title,
			liveStreamItem.Snippet.Description,
		),
		messageParams,
	)
	if err != nil {
		log.Printf("error sending message to channel: %s", err.Error())
	}
}

func mindYoutube() {
	if buf, err := ioutil.ReadFile(youtubeConfigFile); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(buf, &youtubeCFG); err != nil {
			log.Fatal(err)
		}
	}

	ticker := time.Tick(5 * time.Minute)
	for {
		for _, channel := range youtubeCFG.Channels {
			var liveStreams youtubeLiveStreamsResponse
			url := fmt.Sprintf(youtubeLiveStreamSearchURL, youtubeBaseURL, channel, youtubeCFG.Key)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error fetching channel information for %s: %s", channel, err.Error())
				continue
			}
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&liveStreams); err != nil {
				resp.Body.Close()
				log.Printf("Error parsing channel information for %s: %s", channel, err.Error())
				continue
			}
			var newLiveStreams = []string{}
			for _, liveStreamItem := range liveStreams.Items {
				newLiveStreams = append(newLiveStreams, liveStreamItem.ID.VideoId)
				var found = false
				for _, old := range youtubeStreams[channel] {
					if old == liveStreamItem.ID.VideoId {
						found = true
						break
					}
				}
				if found == true {
					continue
				}
				doYoutubeMessage(liveStreamItem)
			}
			youtubeStreams[channel] = newLiveStreams
		}
		<-ticker
	}
}
