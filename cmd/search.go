/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/daoleno/tgraph"
	"github.com/jeehoon/graylog-cli/pkg/graylog/client"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &client.Config{
			Endpoint: ServerEndpoint,
			Username: Username,
			Password: Password,
		}

		graylog := client.NewClient(cfg)

		//requestId := "671a67509fa9a642c4ab9041"
		requestId := "aaaaaaaaaaaaaaaaaaaaaaaa"
		queryId := "11111111-1111-1111-1111-111111111111"
		messageId := "22222222-2222-2222-2222-222222222222"
		histogramId := "33333333-3333-3333-3333-333333333333"
		termsId := "44444444-4444-4444-4444-444444444444"

		q := "*"
		if len(args) != 0 {
			q = args[0]
		}

		req := client.NewSearchRequest(requestId)
		query := client.NewSearchQuery(queryId)
		query.SetQuery(q)

		if SearchFrom != "" && SearchTo != "" {
			query.SetTimerangeAbsolute(SearchFrom, SearchTo)
		} else {
			query.SetTimerangeRelative(int(SearchRange / time.Second))
		}

		if Histogram {
			query.AppendSearchHistogram(histogramId)
		} else if TermsTop != "" {
			query.AppendSearchTop(termsId, TermsTop, 20)
		} else {
			query.AppendSearchMessage(messageId, Limit, Offset, Sort)
		}

		req.AddQuery(query)

		if _, err := graylog.Post("/api/views/search", req); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v", err)
			os.Exit(1)
		}

		httpResp, err := graylog.Post("/api/views/search/"+requestId+"/execute", nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v", err)
			os.Exit(1)
		}

		var resp *client.SearchResponse
		if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v", err)
			os.Exit(1)
		}

		dec := client.NewDecoder(DecoderConfig)

		result, has := resp.Results[queryId]
		if !has {
			fmt.Fprintf(os.Stderr, "ERROR: not found query result of %v", queryId)
			os.Exit(1)
		}

		if typ, has := result.SearchTypes[messageId]; has {
			for idx := len(typ.Messages) - 1; idx >= 0; idx-- {
				msg := typ.Messages[idx]
				fmt.Println(client.Render(dec, true, msg.Message))
			}
			fmt.Printf("= Range: %v ~ %v\n", typ.EffectiveTimerange.From, typ.EffectiveTimerange.To)
			fmt.Printf("= Total: %v\n", typ.TotalResults)
		}

		if typ, has := result.SearchTypes[termsId]; has {

			labels := []string{}
			data := [][]float64{}

			for _, row := range typ.Rows {
				if len(row.Key) == 0 {
					continue
				}

				labels = append(labels, row.Key[0])
				data = append(data, []float64{row.Values[0].Value})
			}
			tgraph.Chart("", labels, data, nil, nil, 50, false, Tick)
			fmt.Printf("= Range: %v ~ %v\n", typ.EffectiveTimerange.From, typ.EffectiveTimerange.To)
			fmt.Printf("= Total: %v\n", typ.Total)
		}

		if typ, has := result.SearchTypes[histogramId]; has {

			labels := []string{}
			data := [][]float64{}

			for _, row := range typ.Rows {
				if len(row.Key) == 0 {
					continue
				}

				labels = append(labels, row.Key[0])
				data = append(data, []float64{row.Values[0].Value})
			}
			tgraph.Chart("", labels, data, nil, nil, 50, false, Tick)
			fmt.Printf("= Range: %v ~ %v\n", typ.EffectiveTimerange.From, typ.EffectiveTimerange.To)
			fmt.Printf("= Total: %v\n", typ.Total)
		}
		fmt.Printf("= State: %v\n", result.State)
		if len(result.Errors) != 0 {
			fmt.Printf("= Errors: %v\n", result.Errors)
		}
		fmt.Printf("= Query: %v\n", result.Query.Query.QueryString)
		fmt.Println()

	},
}

type UintSlice []uint64

func (x UintSlice) Len() int           { return len(x) }
func (x UintSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x UintSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x UintSlice) Sort()              { sort.Sort(x) }

var (
	Histogram = false
	TermsTop  = ""
	Tick      = "■"
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().SortFlags = false

	searchCmd.Flags().IntVar(&Offset, "offset", Offset, "")
	searchCmd.Flags().IntVar(&Limit, "limit", Limit, "")
	searchCmd.Flags().StringVar(&Sort, "sort", Sort, "")

	// Search
	searchCmd.Flags().StringSliceVar(&DecoderConfig.HostnameKeys, "hostname", DecoderConfig.HostnameKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.TimestampKeys, "timestamp", DecoderConfig.TimestampKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.LevelKeys, "level", DecoderConfig.LevelKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.TextKeys, "text", DecoderConfig.TextKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.FieldKeys, "fields", DecoderConfig.FieldKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.SkipFieldKeys, "skip-fields", DecoderConfig.SkipFieldKeys, "")

	// Histogram
	searchCmd.Flags().BoolVarP(&Histogram, "histogram", "H", Histogram, "")
	searchCmd.Flags().StringVarP(&TermsTop, "top", "T", TermsTop, "")
	searchCmd.Flags().StringVar(&Tick, "tick", Tick, "")
}
