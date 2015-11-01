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
	youtubeChannelSearchURL    = "%s/channels?part=id&forUsername=%s&key=%s"
	youtubeLiveStreamSearchURL = "%s/search?part=snippet&eventType=live&type=video&channelId=%s&key=%s"
	youtubeLinkURL             = "https://gaming.youtube.com/watch?v=%s"
	youtubeChannels            = make(map[string][]string)
	youtubeStreams             = make(map[string][]string)
)

type youtubeConfig struct {
	Key   string   `json:"api_key"`
	Users []string `json:"users"`
}

type youtubeChannelsResponse struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

type youtubeLiveStreamsResponse struct {
	Items []struct {
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
	} `json:"items"`
}

func init() {
	flag.StringVar(&youtubeConfigFile, "youtube", youtubeConfigFile, "path to the youtube config file")
}

func mindYoutube() {
	if buf, err := ioutil.ReadFile(youtubeConfigFile); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(buf, &youtubeCFG); err != nil {
			log.Fatal(err)
		}
	}

	for _, user := range youtubeCFG.Users {
		youtubeChannels[user] = []string{}
		var userChannels youtubeChannelsResponse
		url := fmt.Sprintf(youtubeChannelSearchURL, youtubeBaseURL, user, youtubeCFG.Key)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Unable to fetch channels for %s: %s", user, err.Error())
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&userChannels); err != nil {
			log.Fatal("Unable to parse channels response for %s: %s", user, err.Error())
		}
		resp.Body.Close()
		for _, channel := range userChannels.Items {
			log.Printf("Found channel for user %s: %s", user, channel.ID)
			youtubeChannels[user] = append(youtubeChannels[user], channel.ID)
		}
	}

	ticker := time.Tick(5 * time.Minute)
	for {
		for user, channels := range youtubeChannels {
			for _, channel := range channels {
				var liveStreams youtubeLiveStreamsResponse
				url := fmt.Sprintf(youtubeLiveStreamSearchURL, youtubeBaseURL, channel, youtubeCFG.Key)
				resp, err := http.Get(url)
				if err != nil {
					log.Printf("Error fetching channel information for %s (%s): %s", user, channel, err.Error())
					continue
				}
				dec := json.NewDecoder(resp.Body)
				if err := dec.Decode(&liveStreams); err != nil {
					resp.Body.Close()
					log.Printf("Error parsing channel information for %s (%s): %s", user, channel, err.Error())
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

					messageParams := slack.NewPostMessageParameters()
					messageParams.AsUser = true
					messageParams.Parse = "full"
					messageParams.LinkNames = 1
					messageParams.UnfurlMedia = true
					messageParams.UnfurlLinks = true
					messageParams.EscapeText = false
					messageParams.Attachments = append(messageParams.Attachments, slack.Attachment{
						Title: fmt.Sprintf(
							"Watch %s play %s",
							liveStreamItem.Snippet.ChannelTitle,
							liveStreamItem.Snippet.Title,
						),
						TitleLink: fmt.Sprintf(youtubeLinkURL, liveStreamItem.ID.VideoId),
						ThumbURL:  liveStreamItem.Snippet.Thumbnails.Default.URL,
					})
					_, _, err := slackClient.PostMessage(
						slackCFG.Channel,
						fmt.Sprintf(
							"*%s* has begun streaming *%s*",
							liveStreamItem.Snippet.ChannelTitle,
							liveStreamItem.Snippet.Title,
						),
						messageParams,
					)
					if err != nil {
						log.Printf("error sending message to channel: %s", err.Error())
					}
				}
				youtubeStreams[channel] = newLiveStreams
			}
		}
		<-ticker
	}
}
