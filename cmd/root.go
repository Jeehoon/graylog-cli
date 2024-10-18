package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeehoon/graylog-cli/pkg/graylog"
	"github.com/jeehoon/graylog-cli/pkg/render"
)

var rootCmd = &cobra.Command{
	Use:   "graylog-cli",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		client := graylog.NewClient(ClientConfig)

		query := "*"
		if len(args) != 0 {
			query = args[0]
		}

		var from time.Time
		var to time.Time
		var relative = false
		var err error

		if SearchFrom != "" && SearchTo != "" {
			from, err = time.Parse(time.RFC3339Nano, SearchFrom)
			to, err = time.Parse(time.RFC3339Nano, SearchTo)
		} else if SearchFrom == "" && SearchTo != "" {
			to, err = time.Parse(time.RFC3339Nano, SearchTo)
			from = to.Add(-SearchRange)
		} else if SearchFrom != "" && SearchTo == "" {
			from, err = time.Parse(time.RFC3339Nano, SearchFrom)
			to = from.Add(SearchRange)
		} else {
			relative = true
		}

		var resp *graylog.Response

		if relative {
			resp, err = client.Relative(SearchRange, query)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
		} else {
			resp, err = client.Absolute(from, to, query)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		finfo, err := os.Stdout.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		useColor := finfo.Mode()&os.ModeCharDevice == os.ModeCharDevice

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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	SearchFrom  = ""
	SearchTo    = ""
	SearchRange = 8 * time.Hour

	ClientConfig = &graylog.ClientConfig{
		Server:   "https://127.0.0.1",
		Username: "",
		Password: "",
		Filter:   "streams:000000000000000000000001",
		Offset:   0,
		Limit:    150,
		Sort:     "timestamp:desc",
	}

	DecoderConfig = &graylog.DecoderConfig{
		HostnameKeys: []string{
			"hostname",
			"source",
		},
		TimestampKeys: []string{
			"timestamp",
		},
		LevelKeys: []string{
			"level",
		},
		TextKeys: []string{
			"message",
		},
		SkipFieldKeys: []string{
			"streams",
			"hostname",
			"input",
			"gl2_source_input",
			"gl2_remote_ip",
			"gl2_accounted_message_size",
			"gl2_message_id",
			"gl2_source_node",
			"gl2_remote_port",
			"file",
			"function",
			"line",
			"timestamp",
			"_id",
			"source",
			"message",
			"level",
			"caller",
		},
		FieldKeys: []string{},
	}
)

func init() {
	rootCmd.Flags().SortFlags = false

	rootCmd.Flags().StringVar(&SearchFrom, "from", SearchFrom, "")
	rootCmd.Flags().StringVar(&SearchTo, "to", SearchTo, "")
	rootCmd.Flags().DurationVar(&SearchRange, "range", SearchRange, "")

	rootCmd.Flags().StringVar(&ClientConfig.Server, "server", ClientConfig.Server, "")
	rootCmd.Flags().StringVar(&ClientConfig.Username, "username", ClientConfig.Username, "")
	rootCmd.Flags().StringVar(&ClientConfig.Password, "password", ClientConfig.Password, "")
	rootCmd.Flags().StringVar(&ClientConfig.Filter, "filter", ClientConfig.Filter, "")
	rootCmd.Flags().UintVar(&ClientConfig.Offset, "offset", ClientConfig.Offset, "")
	rootCmd.Flags().UintVar(&ClientConfig.Limit, "limit", ClientConfig.Limit, "")
	rootCmd.Flags().StringVar(&ClientConfig.Sort, "sort", ClientConfig.Sort, "")

	rootCmd.Flags().StringSliceVar(&DecoderConfig.HostnameKeys, "hostname", DecoderConfig.HostnameKeys, "")
	rootCmd.Flags().StringSliceVar(&DecoderConfig.TimestampKeys, "timestamp", DecoderConfig.TimestampKeys, "")
	rootCmd.Flags().StringSliceVar(&DecoderConfig.LevelKeys, "level", DecoderConfig.LevelKeys, "")
	rootCmd.Flags().StringSliceVar(&DecoderConfig.TextKeys, "text", DecoderConfig.TextKeys, "")
	rootCmd.Flags().StringSliceVar(&DecoderConfig.FieldKeys, "fields", DecoderConfig.FieldKeys, "")
	rootCmd.Flags().StringSliceVar(&DecoderConfig.SkipFieldKeys, "skip-fields", DecoderConfig.SkipFieldKeys, "")
}
