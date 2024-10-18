package graylog

import (
	"encoding/json"
	"fmt"
)

type HistogramResponse struct {
	Results          map[uint64]uint64 `json:"results"`
	Time             int               `json:"time"`
	BuiltQuery       string            `json:"built_query"`
	QueriedTimerange Timerange         `json:"queried_timerange"`
	Interval         string            `json:"interval"`
}

type Timerange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (client *Client) Histogram(query *Query, interval string) (resp *HistogramResponse, err error) {

	vars, err := client.parseQuery(query)
	if err != nil {
		return nil, err
	}
	vars.Set("interval", interval)

	var path string
	if vars.Has("range") {
		path = fmt.Sprintf("/api/search/universal/relative/histogram?%v", vars.Encode())
	} else {
		path = fmt.Sprintf("/api/search/universal/absolute/histogram?%v", vars.Encode())
	}

	httpResp, err := client.request(path, query)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}
