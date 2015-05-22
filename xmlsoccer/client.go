package xmlsoccer

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is to call webservice
type Client struct {
	client *http.Client
	// it will be zero value ("") while not in testing
	testURL string

	// API key to access the service
	APIKey string

	// the base url for webservice
	BaseURL string
}

var (
	// ErrMissingAPIKey represents error when client makes request without APIKey
	ErrMissingAPIKey = errors.New("APIKey is requried")
)

// DemoClient creates client which access to Demo webservice
func DemoClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "http://www.xmlsoccer.com/FootballDataDemo.asmx",
	}
}

// FullClient create Client which access to Full webservice
func FullClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "http://www.xmlsoccer.com/FootballData.asmx",
	}
}

func (c *Client) postURL(service string) string {
	if c.testURL != "" {
		return c.testURL
	}

	return c.BaseURL + service
}

// GetAllLeagues returns all published leagues
func (c *Client) GetAllLeagues() ([]League, error) {
	if len(c.APIKey) == 0 {
		return nil, ErrMissingAPIKey
	}

	resp, err := http.PostForm(c.postURL("/GetAllLeagues"),
		url.Values{"ApiKey": {c.APIKey}})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := xmlroot{}
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}

	return result.Leagues, nil
}
