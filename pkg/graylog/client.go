package graylog

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	cfg *ClientConfig
}

type Query struct {
	From  time.Time
	To    time.Time
	Range time.Duration
	Query string
}

type ClientConfig struct {
	Username string
	Password string
	Server   string
	Filter   string
	Offset   uint
	Limit    uint
	Sort     string
	Verbose  bool
}

func NewClient(cfg *ClientConfig) *Client {
	client := new(Client)
	client.cfg = cfg

	return client
}

func (client *Client) parseQuery(query *Query) (vars url.Values, err error) {
	if query.Range == 0 && (query.From.IsZero() || query.To.IsZero()) {
		return nil, errors.Errorf("invalid input: query input")
	}

	vars = url.Values{}
	vars.Set("query", query.Query)
	vars.Set("offset", fmt.Sprintf("%d", client.cfg.Offset))
	vars.Set("limit", fmt.Sprintf("%d", client.cfg.Limit))

	if client.cfg.Sort != "" {
		vars.Set("sort", client.cfg.Sort)
	}

	if client.cfg.Filter != "" {
		vars.Set("filter", client.cfg.Filter)
	}

	if query.From.IsZero() && query.To.IsZero() && query.Range != 0 {
		vars.Set("range", fmt.Sprintf("%.0f", query.Range.Seconds()))
	}

	if !query.From.IsZero() {
		vars.Set("from", query.From.Format(time.RFC3339Nano))
		if query.To.IsZero() {
			vars.Set("to", query.From.Add(query.Range).Format(time.RFC3339Nano))
		}
	}

	if !query.To.IsZero() {
		vars.Set("to", query.To.Format(time.RFC3339Nano))
		if query.From.IsZero() {
			vars.Set("from", query.To.Add(-query.Range).Format(time.RFC3339Nano))
		}
	}

	return vars, nil
}

func (client *Client) request(path string, query *Query) (httpResp *http.Response, err error) {

	url := client.cfg.Server + path

	httpClient := new(http.Client)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.SetBasicAuth(client.cfg.Username, client.cfg.Password)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	if client.cfg.Verbose {
		dump, err := httputil.DumpRequest(httpReq, false)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(os.Stderr, "%v\n", string(dump))
	}

	httpResp, err = httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode > 299 {
		b, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.Errorf("%v / %v", httpResp.Status, string(b))
	}

	return httpResp, err
}
