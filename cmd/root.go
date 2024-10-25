package cmd

import (
	"os"
	"time"

	"github.com/jeehoon/graylog-cli/pkg/graylog/client"
	"github.com/spf13/cobra"
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
	SearchFrom     = ""
	SearchTo       = ""
	SearchRange    = 8 * time.Hour
	ServerEndpoint = "https://127.0.0.1"
	Username       = ""
	Password       = ""
	Offset         = 0
	Limit          = 150
	Sort           = "timestamp:DESC"
	Verbose        = false

	DecoderConfig = &client.DecoderConfig{
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

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", Verbose, "")

	rootCmd.PersistentFlags().StringVar(&SearchFrom, "from", SearchFrom, "")
	rootCmd.PersistentFlags().StringVar(&SearchTo, "to", SearchTo, "")
	rootCmd.PersistentFlags().DurationVar(&SearchRange, "range", SearchRange, "")

	rootCmd.PersistentFlags().StringVar(&ServerEndpoint, "server", ServerEndpoint, "")
	rootCmd.PersistentFlags().StringVar(&Username, "username", Username, "")
	rootCmd.PersistentFlags().StringVar(&Password, "password", Password, "")
}
