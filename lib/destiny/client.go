package destiny

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	ErrDestinyAccountNotFound = fmt.Errorf("We were unable to find your Destiny account information. If you have a valid Destiny Account, let us know.")
)

var errorMap = map[int]error{
	1601: ErrDestinyAccountNotFound,
}

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
}

// Client provides an interface to BungieNet.Platform.DestinyServices
// see: https://www.bungie.net/platform/destiny/help/
type Client struct {
	apiKey    string
	client    *http.Client
	baseURL   string
	userAgent string
	debugURL  bool
	wait      *time.Time
	lock      sync.Mutex
}

// Request is the basic work unit of this API. It allows you to select
// the API endpoint that you desire and also add url parameters to the
// request and then finally make the request and get the data back out
type Request struct {
	*url.URL
	url.Values
	c *Client
}

// ToURL returns the request url as a string
func (r *Request) ToURL() string {
	if v := r.Values.Encode(); v != "" {
		return fmt.Sprintf("%s%s?%s", r.c.baseURL, r.URL.String()[1:], v)
	}
	return fmt.Sprintf("%s%s", r.c.baseURL, r.URL.String()[1:])
}

// DebugInto works as Into except it enables logging of the fetched URL
// to stderr
func (r *Request) DebugInto(into interface{}) error {
	r.c.debugURL = true
	err := r.Into(into)
	r.c.debugURL = false
	return err
}

// Into makes the request and unmarshals the response in one of three ways
// depending on the format of the response...
//
// First if there is a ->Response->data (as in CharacterSummary) Response->data
// will be unmarshalled into the data structure supplied
//
// Second if there is a ->Response->results
// (as in GetAdminsOfGroup ( http://bungienetplatform.wikia.com/wiki/GetAdminsOfGroup ) )
// then Response->results will be unmarshalled into the data structure supplied.
// Note there may not be any of these kinds of api calles added to this client as of yes
// but when they are this is here to ensure that they integrate smoothly
//
// Lastly if there is nither a ->data nor a ->results (as in SearchDestinyPlayer)
// then the ->Response itself will be unmarshalled into the provided data structure
//
// Example:
//		type bungieSearchResponse []struct {
//			MembershipId string `json:"membershipId"`
//		}
//		var d bungieSearchResponse
//		req, _ := destinyClient.SearchDestinyPlayer(destiny.PlatformXBL, user.GamerTag)
//		_ := req.Into(&d)
//		if len(d) {
//			user.DestinyID = d[0].MembershipId
//		}
func (r *Request) Into(into interface{}) error {
	if v := r.Values.Encode(); v != "" {
		return r.c.get(r.URL.String()[1:]+"?"+v, into)
	}
	return r.c.get(r.URL.String()[1:], into)
}

// CustomURL Allows you to input any URL you desire into the client and make a request against it
// in this way you can use API endpoints not yet added to the client. Also the helper endpoints
// all use this one function internally.
func (c *Client) CustomURL(URL string) (*Request, error) {
	u, e := url.ParseRequestURI(URL)
	return &Request{URL: u, c: c, Values: url.Values{}}, e
}

// SearchDestinyPlayer returns a list of Destiny memberships given a full Gamertag or PSN ID.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=SearchDestinyPlayer%2f%7bmembershipType%7d%2f%7bdisplayName%7d%2f
// see: http://bungienetplatform.wikia.com/wiki/SearchDestinyPlayer
func (c *Client) SearchDestinyPlayer(platform int, gamertag string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/SearchDestinyPlayer/%d/%s/", platform, gamertag))
}

// AccountSummary returns Destiny account information for the supplied membership
// in a compact summary form.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=%7bmembershipType%7d%2fAccount%2f%7bdestinyMembershipId%7d%2fSummary%2f
// see: http://bungienetplatform.wikia.com/wiki/GetDestinyAccountSummary
func (c *Client) AccountSummary(platform int, id string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/%d/Account/%s/Summary/", platform, id))
}

// AccountItems returns information about all items on the for the supplied
// Destiny Membership ID, and a minimal set of character information so that
// it can be used.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=%7bmembershipType%7d%2fAccount%2f%7bdestinyMembershipId%7d%2fItems%2f
func (c *Client) AccountItems(platform int, id string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/%d/Account/%s/Items/", platform, id))
}

// ActivityHistory returns activity history stats for indicated character.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=Stats%2fActivityHistory%2f%7bmembershipType%7d%2f%7bdestinyMembershipId%7d%2f%7bcharacterId%7d%2f
// see: http://bungienetplatform.wikia.com/wiki/GetActivityHistory
func (c *Client) ActivityHistory(platform int, id string, cid string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/Stats/ActivityHistory/%d/%s/%s/", platform, id, cid))
}

// CharacterSummary returns a character summary for the supplied membership.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=%7bmembershipType%7d%2fAccount%2f%7bdestinyMembershipId%7d%2fCharacter%2f%7bcharacterId%7d%2f
// see: http://bungienetplatform.wikia.com/wiki/GetCharacterSummary
func (c *Client) CharacterSummary(platform int, id string, cid string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/%d/Account/%s/Character/%s/", platform, id, cid))
}

// CharacterActivities returns activity progress for a given character.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=%7bmembershipType%7d%2fAccount%2f%7bdestinyMembershipId%7d%2fCharacter%2f%7bcharacterId%7d%2fActivities%2f
// see: http://bungienetplatform.wikia.com/wiki/GetCharacterActivities
func (c *Client) CharacterActivities(platform int, id string, cid string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/%d/Account/%s/Character/%s/Activities/", platform, id, cid))
}

// CharacterInventory returns summary information for the inventory for the supplied character.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=%7bmembershipType%7d%2fAccount%2f%7bdestinyMembershipId%7d%2fCharacter%2f%7bcharacterId%7d%2fInventory%2fSummary%2f
func (c *Client) CharacterInventory(platform int, id string, cid string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/%d/Account/%s/Character/%s/Inventory/Summary/", platform, id, cid))
}

// CharacterProgression returns the progression details for the supplied character.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=%7bmembershipType%7d%2fAccount%2f%7bdestinyMembershipId%7d%2fCharacter%2f%7bcharacterId%7d%2fProgression%2f
// see: http://bungienetplatform.wikia.com/wiki/GetCharacterProgression
func (c *Client) CharacterProgression(platform int, id string, cid string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/%d/Account/%s/Character/%s/Progression/", platform, id, cid))
}

// AggregateActivityStats Returns all activities the character has participated in together with aggregate statistics for those activities.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=Stats%2fAggregateActivityStats%2f%7bmembershipType%7d%2f%7bdestinyMembershipId%7d%2f%7bcharacterId%7d%2f
// see: http://bungienetplatform.wikia.com/wiki/GetDestinyAggregateActivityStats
func (c *Client) AggregateActivityStats(platform int, id string, cid string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/Stats/AggregateActivityStats/%d/%s/%s/", platform, id, cid))
}

// UniqueWeapons Returns details about unique weapon usage, including all exotic weapons.
// see: https://www.bungie.net/platform/destiny/help/HelpDetail/GET?uri=Stats%2fUniqueWeapons%2f%7bmembershipType%7d%2f%7bdestinyMembershipId%7d%2f%7bcharacterId%7d%2f
// see: http://bungienetplatform.wikia.com/wiki/GetUniqueWeaponHistory
func (c *Client) UniqueWeapons(platform int, id string, cid string) (*Request, error) {
	return c.CustomURL(fmt.Sprintf("/Stats/UniqueWeapons/%d/%s/%s/", platform, id, cid))
}

// Platform returns a platform client which will obviate the need to provide
// platform distinctions for helper methods (convenient when you're only going
// to be working against one platform)
func (c *Client) Platform(platformID int) *Platform {
	return &Platform{
		c:  c,
		id: platformID,
	}
}

func (c *Client) get(uri string, into interface{}) error {
	c.lock.Lock()
	if c.wait != nil {
		time.Sleep(time.Now().Sub(*c.wait))
		c.wait = nil
	}
	// Prepare our request
	if c.debugURL {
		log.Println(fmt.Sprintf("destiny.Client debug: %s%s", c.baseURL, uri))
	}
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
		if e, ok := errorMap[envelope.ErrorCode]; ok {
			return e
		}
		return fmt.Errorf("API returned error code %d: %s %s", envelope.ErrorCode, envelope.ErrorStatus, envelope.Message)
	}

	if envelope.ThrottleSeconds > 0 {
		t := time.Now().Add(time.Duration(envelope.ThrottleSeconds) * time.Second)
		log.Println("Destiny API Servers requested throttling:", (time.Duration(envelope.ThrottleSeconds) * time.Second).String())
		c.wait = &t
	}

	var possibleResponse struct {
		Data    *json.RawMessage `json:"data,omitempty"`
		Results *json.RawMessage `json:"results,omitempty"`
	}
	if err := json.Unmarshal(*envelope.Response, &possibleResponse); err == nil {
		if possibleResponse.Data != nil {
			return json.Unmarshal(*possibleResponse.Data, &into)
		}
		if possibleResponse.Results != nil {
			return json.Unmarshal(*possibleResponse.Results, &into)
		}
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
