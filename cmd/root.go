package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeehoon/graylog-cli/pkg/graylog"
)

var rootCmd = &cobra.Command{
	Use:   "graylog-cli",
	Short: "A brief description of your application",
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
	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.PersistentFlags().StringVar(&SearchFrom, "from", SearchFrom, "")
	rootCmd.PersistentFlags().StringVar(&SearchTo, "to", SearchTo, "")
	rootCmd.PersistentFlags().DurationVar(&SearchRange, "range", SearchRange, "")

	rootCmd.PersistentFlags().StringVar(&ClientConfig.Server, "server", ClientConfig.Server, "")
	rootCmd.PersistentFlags().StringVar(&ClientConfig.Username, "username", ClientConfig.Username, "")
	rootCmd.PersistentFlags().StringVar(&ClientConfig.Password, "password", ClientConfig.Password, "")
	rootCmd.PersistentFlags().StringVar(&ClientConfig.Filter, "filter", ClientConfig.Filter, "")
}
