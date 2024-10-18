package graylog

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	From            string     `json:"from"`
	To              string     `json:"to"`
	UsedIndices     any        `json:"used_indices"`
	Time            int        `json:"time"`
	TotalResults    int        `json:"total_results"`
	Fields          []string   `json:"fields"`
	DecorationStats any        `json:"decoration_stats"`
	Query           string     `json:"query"`
	BuiltQuery      string     `json:"built_query"`
	Messages        []*Message `json:"messages"`
}

type Message struct {
	DecorationStats any            `json:"decoration_stats"`
	HighlightRanges any            `json:"highlight_ranges"`
	Message         map[string]any `json:"message"`
	Index           string         `json:"index"`
}

func (client *Client) Query(query *Query) (resp *Response, err error) {

	vars, err := client.parseQuery(query)
	if err != nil {
		return nil, err
	}

	var path string
	if vars.Has("range") {
		path = fmt.Sprintf("/api/search/universal/relative?%v", vars.Encode())
	} else {
		path = fmt.Sprintf("/api/search/universal/absolute?%v", vars.Encode())
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
