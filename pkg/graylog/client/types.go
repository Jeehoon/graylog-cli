package client

type Filter struct {
	Type    string          `json:"type"`
	Filters []*FilterStream `json:"filters"`
}

type FilterStream struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Query struct {
	Type        string `json:"type"`
	QueryString string `json:"query_string"`
}

type Timerange struct {
	Type  string `json:"type"`
	Range int    `json:"range,omitempty"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
}

type SearchType interface{}

type SearchTypeMessage struct {
	Id     string  `json:"id"`
	Type   string  `json:"type"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
	Sort   []*Sort `json:"sort"`
}

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type SearchTypePivot struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	Rollup    bool        `json:"rollup"`
	Series    []*Series   `json:"series"`
	RowGroups []*RowGroup `json:"row_groups"`
}

type Series struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type RowGroup struct {
	Type     string    `json:"type"`
	Field    string    `json:"field"`
	Limit    int       `json:"limit,omitempty"`
	Interval *Interval `json:"interval,omitempty"`
}

type Interval struct {
	Type    string `json:"type"`
	Scaling int    `json:"scaling"`
}

type Result struct {
	Errors         []any                        `json:"errors"`
	State          string                       `json:"state"`
	Query          *SearchQuery                 `json:"query"`
	ExecutionStats map[string]any               `json:"execution_stats"`
	SearchTypes    map[string]*SearchTypeResult `json:"search_types"`
}

type SearchTypeResult struct {
	Id                 string     `json:"id"`
	Type               string     `json:"type"`
	EffectiveTimerange *Timerange `json:"effective_timerange"`
	Rows               []*Row     `json:"rows"`
	Total              uint64     `json:"total"`
	TotalResults       uint64     `json:"total_results"`
	Messages           []*Message `json:"messages"`
}

type Message struct {
	DecorationStats any            `json:"decoration_stats"`
	HighlightRanges any            `json:"highlight_ranges"`
	Message         map[string]any `json:"message"`
	Index           string         `json:"index"`
}

type Row struct {
	Key    []string `json:"key"`
	Source string   `json:"source"`
	Values []*Value `json:"values"`
}

type Value struct {
	Key    []string `json:"key"`
	Rollup bool     `json:"rollup"`
	Source string   `json:"source"`
	Value  float64  `json:"value"`
}

type SearchRequest struct {
	Id      string         `json:"id"`
	Queries []*SearchQuery `json:"queries"`
}

func NewSearchRequest(id string) *SearchRequest {
	req := new(SearchRequest)
	req.Id = id
	req.Queries = nil

	return req
}

func (req *SearchRequest) AddQuery(query *SearchQuery) {
	req.Queries = append(req.Queries, query)
}

type SearchResponse struct {
	Execution map[string]any     `json:"execution"`
	Results   map[string]*Result `json:"results"`
	Id        string             `json:"id"`
	Owner     string             `json:"owner"`
	SearchId  string             `json:"search_id"`
}
