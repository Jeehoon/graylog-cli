/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// histogramCmd represents the histogram command
var histogramCmd = &cobra.Command{
	Use:   "histogram",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var (
	HistogramInterval = "minute"
)

func init() {
	rootCmd.AddCommand(histogramCmd)

	histogramCmd.Flags().StringVar(&HistogramInterval, "interval", HistogramInterval, "")
}
