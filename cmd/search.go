/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/daoleno/tgraph"
	"github.com/jeehoon/graylog-cli/pkg/graylog"
	"github.com/jeehoon/graylog-cli/pkg/render"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		client := graylog.NewClient(ClientConfig)

		var query = new(graylog.Query)
		if len(args) != 0 {
			query.Query = args[0]
		} else {
			query.Query = "*"
		}

		if SearchFrom != "" {
			from, err := time.Parse(time.RFC3339Nano, SearchFrom)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
			query.From = from
		}

		if SearchTo != "" {
			to, err := time.Parse(time.RFC3339Nano, SearchTo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
			query.To = to
		}

		if SearchRange != 0 {
			query.Range = SearchRange
		}

		if Histogram {
			resp, err := client.Histogram(query, HistogramInterval)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}

			labels := []string{}
			data := [][]float64{}
			var keys UintSlice

			for ts := range resp.Results {
				keys = append(keys, ts)
			}
			keys.Sort()

			for _, ts := range keys {
				at := time.Unix(int64(ts), 0).UTC()
				n := resp.Results[ts]

				labels = append(labels, at.Format(time.RFC3339Nano))
				data = append(data, []float64{float64(n)})
			}

			tgraph.Chart("", labels, data, nil, nil, 50, false, Tick)

			fmt.Println()
			fmt.Printf("= Query: %v\n", query.Query)
			fmt.Printf("= Range: %v ~ %v\n", resp.QueriedTimerange.From, resp.QueriedTimerange.To)
			fmt.Printf("= Time:  %v\n", resp.Time)
		} else {
			resp, err := client.Query(query)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}

			useColor := UseColor()
			dec := graylog.NewDecoder(DecoderConfig)
			for idx := len(resp.Messages) - 1; idx >= 0; idx-- {
				msg := resp.Messages[idx]
				fmt.Println(render.Render(dec, useColor, msg))
			}

			fmt.Println()
			fmt.Printf("= Query: %v\n", resp.Query)
			fmt.Printf("= Range: %v ~ %v\n", resp.From, resp.To)
			fmt.Printf("= Total: %v\n", resp.TotalResults)
			fmt.Printf("= Time:  %v\n", resp.Time)
		}

	},
}

type UintSlice []uint64

func (x UintSlice) Len() int           { return len(x) }
func (x UintSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x UintSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x UintSlice) Sort()              { sort.Sort(x) }

var (
	Histogram         = false
	HistogramInterval = "minute"
	Tick              = "■"
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().SortFlags = false

	searchCmd.Flags().UintVar(&ClientConfig.Offset, "offset", ClientConfig.Offset, "")
	searchCmd.Flags().UintVar(&ClientConfig.Limit, "limit", ClientConfig.Limit, "")
	searchCmd.Flags().StringVar(&ClientConfig.Sort, "sort", ClientConfig.Sort, "")

	// Search
	searchCmd.Flags().StringSliceVar(&DecoderConfig.HostnameKeys, "hostname", DecoderConfig.HostnameKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.TimestampKeys, "timestamp", DecoderConfig.TimestampKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.LevelKeys, "level", DecoderConfig.LevelKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.TextKeys, "text", DecoderConfig.TextKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.FieldKeys, "fields", DecoderConfig.FieldKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.SkipFieldKeys, "skip-fields", DecoderConfig.SkipFieldKeys, "")

	// Histogram
	searchCmd.Flags().BoolVarP(&Histogram, "histogram", "H", Histogram, "")
	searchCmd.Flags().StringVar(&HistogramInterval, "interval", HistogramInterval, "")
	searchCmd.Flags().StringVar(&Tick, "tick", Tick, "")
}
