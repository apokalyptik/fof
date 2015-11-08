package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nlopes/slack"
)

var (
	youtubeConfigFile          = "./youtube.json"
	youtubeCFG                 = youtubeConfig{}
	youtubeAPIKey              = ""
	youtubeBaseURL             = "https://www.googleapis.com/youtube/v3"
	youtubeLiveStreamSearchURL = "%s/search?part=snippet&eventType=live&type=video&channelId=%s&key=%s&publishedAfter=%s"
	youtubeChannelDetailsURL   = "%s/channels?part=snippet&id=%s&key=%s"
	youtubeLinkURL             = "https://gaming.youtube.com/watch?v=%s"
	youtubeStreams             = make(map[string][]string)
	youtubeStreamDetails       = make(map[string]youtubeChannelDetails)
)

type youtubeConfig struct {
	Key      string   `json:"api_key"`
	Channels []string `json:"channels"`
}

type youtubeChannelDetails struct {
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"snippet"`
}

type youtubeChannelDetailsResponse struct {
	Items []youtubeChannelDetails `json:"items"`
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
		ChannelID            string `json:"channelId"`
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
		channelTitle = youtubeStreamDetails[liveStreamItem.Snippet.ChannelID].Snippet.Title
	}
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
		log.Printf("[yt] error sending message to channel: %s", err.Error())
	}
}

func mindYoutube() {
	log.Println("[yt] Starting up")
	if buf, err := ioutil.ReadFile(youtubeConfigFile); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(buf, &youtubeCFG); err != nil {
			log.Fatal(err)
		}
	}

	for _, channel := range youtubeCFG.Channels {
		var info youtubeChannelDetailsResponse
		url := fmt.Sprintf(youtubeChannelDetailsURL, youtubeBaseURL, channel, youtubeCFG.Key)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("[yt] Error fetching channel details for %s: %s", channel, err.Error())
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&info); err != nil {
			log.Fatalf("[yt] Error decoding channel details response for %s: %s", channel, err.Error())
		}
		if len(info.Items) < 1 {
			log.Fatal("[yt] No details returned for the channel: %s", channel)
		}
		log.Printf("[yt] Channel %s has title %s", channel, info.Items[0].Snippet.Title)
		youtubeStreamDetails[channel] = info.Items[0]
	}

	loopTime := 5 * time.Minute
	ticker := time.Tick(loopTime)
	lastTime := time.Now().Add(0 - loopTime).Format(time.RFC3339)
	nextLastTime := lastTime
	for {
		lastTime = nextLastTime
		nextLastTime = time.Now().Format(time.RFC3339)
		for _, channel := range youtubeCFG.Channels {
			var liveStreams youtubeLiveStreamsResponse
			log.Println("[yt] Checking for streams on", youtubeStreamDetails[channel].Snippet.Title, "since", lastTime)
			url := fmt.Sprintf(youtubeLiveStreamSearchURL, youtubeBaseURL, channel, youtubeCFG.Key, lastTime)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("[yt] Error fetching channel information for %s: %s", channel, err.Error())
				continue
			}
			dec := json.NewDecoder(resp.Body)
			if err := dec.Decode(&liveStreams); err != nil {
				resp.Body.Close()
				log.Printf("[yt] Error parsing channel information for %s: %s", channel, err.Error())
				continue
			}
			resp.Body.Close()
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
				log.Printf("[yt] Found new live stream %s on channel %s", liveStreamItem.Snippet.Title, liveStreamItem.Snippet.ChannelID)
				doYoutubeMessage(liveStreamItem)
			}
			youtubeStreams[channel] = newLiveStreams
		}
		<-ticker
	}
}
