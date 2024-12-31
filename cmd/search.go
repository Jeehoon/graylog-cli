/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jeehoon/graylog-cli/pkg/graylog/client"
	"github.com/spf13/cobra"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &client.Config{
			Verbose:  Verbose,
			Endpoint: ServerEndpoint,
			Username: Username,
			Password: Password,
		}

		graylog := client.NewClient(cfg)

		requestId, _ := randomHex(12)
		queryId := uuid.New().String()
		messageId := uuid.New().String()
		histogramId := uuid.New().String()
		termsId := uuid.New().String()

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

		// Search
		if err := graylog.Post("/api/views/search", req, nil); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v", err)
			os.Exit(1)
		}

		// Execute
		var execResp *client.ExecuteResponse
		if err := graylog.Post("/api/views/search/"+requestId+"/execute", nil, &execResp); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v", err)
			os.Exit(1)
		}

		var results map[string]*client.Result

		if execResp.Execution.Done {
			results = execResp.Results
		} else {
			for {
				var statusResp *client.ExecuteResponse

				path := fmt.Sprintf("/api/views/searchjobs/%v/%v/status", execResp.ExecutingNode, execResp.Id)
				if err := graylog.Get(path, nil, &statusResp); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v", err)
					os.Exit(1)
				}

				if statusResp.Execution.Done {
					results = statusResp.Results
					break
				}

				time.Sleep(500 * time.Millisecond)
			}
		}

		dec := client.NewDecoder(DecoderConfig)

		result, has := results[queryId]
		if !has {
			fmt.Fprintf(os.Stderr, "ERROR: not found query result of %v", queryId)
			os.Exit(1)
		}

		if typ, has := result.SearchTypes[messageId]; has {
			for idx := len(typ.Messages) - 1; idx >= 0; idx-- {
				msg := typ.Messages[idx]
				fmt.Println(client.Render(dec, true, msg.Message))
			}
			fmt.Printf("========== Messages ==========\n")
			fmt.Printf("= Range: %v ~ %v\n", typ.EffectiveTimerange.From, typ.EffectiveTimerange.To)
			fmt.Printf("= Total: %v\n", typ.TotalResults)
		}

		if typ, has := result.SearchTypes[termsId]; has {

			labels := []string{}
			data := []float64{}

			for _, row := range typ.Rows {
				if len(row.Key) == 0 {
					continue
				}

				key := row.Key[0]
				value := row.Values[0].Value
				labels = append(labels, key)
				data = append(data, value)
			}

			Chart(labels, data, Tick)
			fmt.Printf("========== Top Values of [%v] field ==========\n", TermsTop)
			fmt.Printf("= Range: %v ~ %v\n", typ.EffectiveTimerange.From, typ.EffectiveTimerange.To)
			fmt.Printf("= Total: %v\n", typ.Total)
		}

		if typ, has := result.SearchTypes[histogramId]; has {

			labels := []string{}
			data := []float64{}

			for _, row := range typ.Rows {
				if len(row.Key) == 0 {
					continue
				}

				key := row.Key[0]
				value := row.Values[0].Value
				labels = append(labels, key)
				data = append(data, value)
			}

			Chart(labels, data, Tick)
			fmt.Printf("========== Histogram ==========\n")
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
	searchCmd.Flags().StringSliceVarP(&DecoderConfig.FieldKeys, "fields", "F", DecoderConfig.FieldKeys, "")
	searchCmd.Flags().StringSliceVar(&DecoderConfig.SkipFieldKeys, "skip-fields", DecoderConfig.SkipFieldKeys, "")

	// Histogram
	searchCmd.Flags().BoolVarP(&Histogram, "histogram", "H", Histogram, "")
	searchCmd.Flags().StringVarP(&TermsTop, "top", "T", TermsTop, "")
	searchCmd.Flags().StringVar(&Tick, "tick", Tick, "")
}

func Chart(labels []string, data []float64, tick string) {
	length := len(labels)
	if len(labels) > len(data) {
		length = len(data)
	}

	var file = os.Stdout
	var maxLabelLength int
	var maxValue float64

	for i := 0; i < length; i++ {
		label := labels[i]
		value := data[i]
		if maxLabelLength < len(label) {
			maxLabelLength = len(label)
		}

		if maxValue < value {
			maxValue = value
		}
	}

	maxBarLength := float64(50)
	labelFmt := fmt.Sprintf("%%%ds", maxLabelLength)

	for i := 0; i < length; i++ {
		label := labels[i]
		value := data[i]

		barLength := (value / maxValue) * maxBarLength
		bar := strings.Repeat(tick, int(barLength))

		s := fmt.Sprintf(labelFmt+":%s %.3f", label, bar, value)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		fmt.Fprintln(file, s)
	}
}
