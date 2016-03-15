package uplay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	username string
	password string
	appID    string
	ticket   string
}

func (c *Client) Authenticate() error {
	req, err := http.NewRequest(
		"POST",
		"https://uplayconnect.ubi.com/ubiservices/v2/profiles/sessions",
		bytes.NewBuffer([]byte("{\"rememberMe\":true}")),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Ubi-AppId", c.appID)
	req.Header.Set("Ubi-RequestedPlatformType", "uplay")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var data struct {
		HTTPCode int    `json:"httpCode"`
		Message  string `json:"message"`
		Ticket   string `json:"ticket"`
	}
	if err := dec.Decode(&data); err != nil {
		return nil
	}
	if data.Ticket == "" {
		return fmt.Errorf("Ubisoft login returned HTTP code %d: %s", data.HTTPCode, data.Message)
	}
	c.ticket = data.Ticket
	return nil
}

func New(username, password string) *Client {
	return &Client{
		appID:    "314d4fef-e568-454a-ae06-43e3bece12a6",
		username: username,
		password: password,
	}
}

func (c *Client) UserSearch(platform int, username string) ([]Profile, error) {
	var data = struct {
		Profiles []Profile
	}{}
	platformString, ok := uplayoverlayPlatforms[platform]
	if !ok {
		return nil, fmt.Errorf("Invalid platform type specified")
	}
	url := fmt.Sprintf(
		"https://uplayoverlay.ubi.com/ubiservices/v1/profiles?nameOnPlatform=%s&platformType=%s",
		url.QueryEscape(username),
		url.QueryEscape(platformString),
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Ubi_v1 t=%s", c.ticket))
	req.Header.Set("Ubi-AppId", "f35adcb5-1911-440c-b1c9-48fdc1701c68")
	client := &http.Client{}
	if rsp, err := client.Do(req); err != nil {
		return data.Profiles, fmt.Errorf("Error requesting data: %s", err.Error())
	} else {
		defer rsp.Body.Close()
		dec := json.NewDecoder(rsp.Body)
		if err := dec.Decode(&data); err != nil {
			return data.Profiles, fmt.Errorf("Error decoding data: %s", err.Error())
		}
	}
	return data.Profiles, nil
}

func (c *Client) DivisionStats(platform int, uuid string) ([]DivisionStat, error) {
	var data []struct {
		Stats []DivisionStat `json:"stats"`
	}
	platformString, ok := uplaywebcenterPlatforms[platform]
	if !ok {
		return nil, fmt.Errorf("Invalid platform type specified")
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://uplaywebcenter.ubi.com/v1/stats/playerStats/?game=TCTD&locale=en-GB&platform=%s&userId=%s", url.QueryEscape(platformString), url.QueryEscape(uuid)),
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Ubi_v1 t=%s", c.ticket))
	req.Header.Set("Ubi-AppId", "f35adcb5-1911-440c-b1c9-48fdc1701c68")
	client := &http.Client{}
	if rsp, err := client.Do(req); err != nil {
		return []DivisionStat{}, err
	} else {
		defer rsp.Body.Close()
		dec := json.NewDecoder(rsp.Body)
		dec.UseNumber()
		if err := dec.Decode(&data); err != nil {
			return []DivisionStat{}, err
		}
	}
	if len(data) > 0 {
		return data[0].Stats, nil
	}
	return []DivisionStat{}, nil
}
