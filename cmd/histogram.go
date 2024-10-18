/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/daoleno/tgraph"
	"github.com/jeehoon/graylog-cli/pkg/graylog"
	"github.com/spf13/cobra"
)

// histogramCmd represents the histogram command
var histogramCmd = &cobra.Command{
	Use:   "histogram",
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

		resp, err := client.Histogram(query, HistogramInterval)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		labels := []string{}
		data := [][]float64{}

		for ts, n := range resp.Results {
			at := time.Unix(int64(ts), 0).UTC()

			labels = append(labels, at.Format(time.RFC3339Nano))
			data = append(data, []float64{float64(n)})
		}

		tgraph.Chart("", labels, data, nil, nil, 50, false, Tick)
		fmt.Println("= Time:", resp.Time)
		fmt.Println("= Range:", resp.QueriedTimerange)
	},
}

var (
	HistogramInterval = "minute"
	//Tick              = "▇"
	Tick = "■"
)

func init() {
	rootCmd.AddCommand(histogramCmd)

	histogramCmd.Flags().StringVar(&HistogramInterval, "interval", HistogramInterval, "")
	histogramCmd.Flags().StringVar(&Tick, "tick", Tick, "")
}
