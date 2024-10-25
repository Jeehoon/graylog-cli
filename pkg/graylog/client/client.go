package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Config struct {
	Endpoint string
	Username string
	Password string
}

type Client struct {
	cfg *Config
}

func NewClient(cfg *Config) *Client {
	client := new(Client)
	client.cfg = cfg
	return client
}

func (client *Client) Post(path string, req any) (httpResp *http.Response, err error) {
	var buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return nil, err
	}

	url := client.cfg.Endpoint + path

	httpClient := new(http.Client)
	httpReq, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}

	httpReq.SetBasicAuth(client.cfg.Username, client.cfg.Password)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Requested-By", "XMLHttpRequest")
	httpReq.Header.Set("X-Requested-With", "XMLHttpRequest")

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

	return httpResp, nil
}
