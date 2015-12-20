package destiny

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Platform Definitions
const (
	PlatformXBL = 1
	PlatformPSN = 2
)

// ErrorCode Values (response envelope)
const (
	errorSuccess = 1
)

// The response envelope wrapping all api repsonses
type responseEnvelope struct {
	ErrorCode       int
	ErrorStatus     string
	Message         string
	ThrottleSeconds int64
	Response        *json.RawMessage
	// MessageData  interface{}
}

// Client provides an interface to BungieNet.Platform.DestinyServices
// see: https://www.bungie.net/platform/destiny/help/
type Client struct {
	apiKey    string
	client    *http.Client
	baseURL   string
	userAgent string
	wait      *time.Time
	lock      sync.Mutex
}

type Request struct {
	*url.URL
	url.Values
	c *Client
}

func (r *Request) ToURL() string {
	return fmt.Sprintf("%s%s?%s", r.c.baseURL, r.URL.String()[1:], r.Values.Encode())
}

func (r *Request) Into(into interface{}) error {
	return r.c.get(r.URL.String()[1:]+"?"+r.Values.Encode(), into)
}

func (c *Client) AccountSummary(platform int, id string) (*Request, error) {
	u, e := url.ParseRequestURI(fmt.Sprintf("/%d/Account/%s/Summary/", platform, id))
	return &Request{URL: u, c: c, Values: url.Values{}}, e
}

func (c *Client) ActivityHistory(platform int, id string, cid string) (*Request, error) {
	u, e := url.ParseRequestURI(fmt.Sprintf("/Stats/ActivityHistory/%d/%s/%s/", platform, id, cid))
	return &Request{URL: u, c: c, Values: url.Values{}}, e
}

func (c *Client) get(uri string, into interface{}) error {
	c.lock.Lock()
	if c.wait != nil {
		time.Sleep(time.Now().Sub(*c.wait))
		c.wait = nil
	}
	// Prepare our request
	//log.Println(fmt.Sprintf("%s%s", c.baseURL, uri))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.baseURL, uri), nil)
	c.lock.Unlock()
	if err != nil {
		return err
	}
	req.Header.Add("X-API-KEY", c.apiKey)
	req.Header.Add("User-Agent", c.userAgent)

	// Make our request
	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// Prepare to parse our response
	dec := json.NewDecoder(rsp.Body)
	defer rsp.Body.Close()
	var envelope responseEnvelope

	// Unmarshal the JSON
	if err := dec.Decode(&envelope); err != nil {
		return err
	}

	// Check for failure
	if envelope.ErrorCode != errorSuccess {
		return fmt.Errorf("API returned error code %d: %s %s", envelope.ErrorCode, envelope.ErrorStatus, envelope.Message)
	}

	if envelope.ThrottleSeconds > 0 {
		t := time.Now().Add(time.Duration(envelope.ThrottleSeconds) * time.Second)
		c.wait = &t
	}

	return json.Unmarshal(*envelope.Response, &into)
}

// New returns a new client with which you can make API calls to
// Bungies Destiny API
func New(apiKey, purpose string) *Client {
	return &Client{
		apiKey:    apiKey,
		client:    &http.Client{},
		baseURL:   "http://www.bungie.net/Platform/Destiny/",
		userAgent: fmt.Sprintf("Go (golang; net/http; github.com/apokalyptik/fof/lib/destiny; +%s)", purpose),
	}
}
