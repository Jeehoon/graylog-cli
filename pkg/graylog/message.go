package graylog

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
