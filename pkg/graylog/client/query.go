package client

import "strings"

type SearchQuery struct {
	Id          string       `json:"id"`
	Query       *Query       `json:"query"`
	Timerange   *Timerange   `json:"timerange"`
	Filter      *Filter      `json:"filter"`
	SearchTypes []SearchType `json:"search_types"`
}

func NewSearchQuery(id string) *SearchQuery {
	query := new(SearchQuery)
	query.Id = id
	return query
}

func (query *SearchQuery) SetQuery(q string) {
	query.Query = &Query{
		Type:        "elasticsearch",
		QueryString: q,
	}
}

// "timerange": {"type":"relative","range":3600}
func (query *SearchQuery) SetTimerangeRelative(sec int) {
	query.Timerange = &Timerange{
		Type:  "relative",
		Range: sec,
	}
}

// "timerange": {type: "absolute", from: "2024-10-25T00:12:59.168Z", to: "2024-10-25T00:13:22.510Z"}
func (query *SearchQuery) SetTimerangeAbsolute(from, to string) {
	query.Timerange = &Timerange{
		Type: "absolute",
		From: from,
		To:   to,
	}
}

func (query *SearchQuery) AppendSearchMessage(id string, limit, offset int, sort string) {
	strs := strings.SplitN(sort, ":", 2)
	field := strs[0]
	order := strs[1]

	query.SearchTypes = append(query.SearchTypes, &SearchTypeMessage{
		Id:     id,
		Type:   "messages",
		Limit:  limit,
		Offset: offset,
		Sort: []*Sort{
			&Sort{
				Field: field,
				Order: order,
			},
		},
	})
}

func (query *SearchQuery) AppendSearchTop(id string, field string, limit int) {
	query.SearchTypes = append(query.SearchTypes, &SearchTypePivot{
		Id:     id,
		Type:   "pivot",
		Rollup: true,
		Series: []*Series{
			&Series{
				Type: "count",
				Id:   "count()",
			},
		},
		RowGroups: []*RowGroup{
			&RowGroup{
				Type:  "values",
				Field: field,
				Limit: limit,
			},
		},
	})
}

func (query *SearchQuery) AppendSearchHistogram(id string) {
	query.SearchTypes = append(query.SearchTypes, &SearchTypePivot{
		Id:     id,
		Type:   "pivot",
		Rollup: true,
		Series: []*Series{
			&Series{
				Type: "count",
				Id:   "count()",
			},
		},
		RowGroups: []*RowGroup{
			&RowGroup{
				Type:  "time",
				Field: "timestamp",
				Interval: &Interval{
					Type:    "auto",
					Scaling: 1,
				},
			},
		},
	})
}

func (query *SearchQuery) SetFilter(streams ...string) {

	query.Filter = &Filter{
		Type:    "or",
		Filters: nil,
	}

	for _, stream := range streams {
		query.Filter.Filters = append(query.Filter.Filters, &FilterStream{
			Type: "stream",
			Id:   stream,
		})
	}
}
