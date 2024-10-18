package graylog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	cfg *ClientConfig
}

type ClientConfig struct {
	Username string
	Password string
	Server   string
	Filter   string
	Offset   uint
	Limit    uint
	Sort     string
}

func NewClient(cfg *ClientConfig) *Client {
	client := new(Client)
	client.cfg = cfg

	return client
}

func (client *Client) Absolute(from, to time.Time, query string) (resp *Response, err error) {
	vars := url.Values{}
	vars.Set("query", query)
	vars.Set("from", from.Format(time.RFC3339Nano))
	vars.Set("to", to.Format(time.RFC3339Nano))
	target := "/api/search/universal/absolute?"

	return client.request(target, vars)
}

func (client *Client) Relative(duration time.Duration, query string) (resp *Response, err error) {
	vars := url.Values{}
	vars.Set("query", query)
	vars.Set("range", fmt.Sprintf("%.0f", duration.Seconds()))
	target := "/api/search/universal/relative?"

	return client.request(target, vars)

}

func (client *Client) request(target string, vars url.Values) (resp *Response, err error) {
	vars.Set("offset", fmt.Sprintf("%d", client.cfg.Offset))
	vars.Set("limit", fmt.Sprintf("%d", client.cfg.Limit))
	vars.Set("sort", client.cfg.Sort)
	vars.Set("filter", client.cfg.Filter)

	url := client.cfg.Server + target + vars.Encode()

	httpClient := new(http.Client)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.SetBasicAuth(client.cfg.Username, client.cfg.Password)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode > 299 {
		b, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.Errorf("%v %v", httpResp.Status, string(b))
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}
