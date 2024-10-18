/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

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
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().SortFlags = false

	searchCmd.Flags().UintVar(&ClientConfig.Offset, "offset", ClientConfig.Offset, "")
	searchCmd.Flags().UintVar(&ClientConfig.Limit, "limit", ClientConfig.Limit, "")
	searchCmd.Flags().StringVar(&ClientConfig.Sort, "sort", ClientConfig.Sort, "")

	searchCmd.Flags().StringSliceVar(&DecoderConfig.HostnameKeys, "hostname", DecoderConfig.HostnameKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.TimestampKeys, "timestamp", DecoderConfig.TimestampKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.LevelKeys, "level", DecoderConfig.LevelKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.TextKeys, "text", DecoderConfig.TextKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.FieldKeys, "fields", DecoderConfig.FieldKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.SkipFieldKeys, "skip-fields", DecoderConfig.SkipFieldKeys, "")
}
