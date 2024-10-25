package client

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jeehoon/graylog-cli/pkg/timeutil"
)

type DecoderConfig struct {
	HostnameKeys  []string
	TimestampKeys []string
	LevelKeys     []string
	TextKeys      []string
	FieldKeys     []string
	SkipFieldKeys []string
}

type Decoder struct {
	cfg           *DecoderConfig
	skipFieldKeys map[string]struct{}
}

func DefaultDecoderConfig() *DecoderConfig {
	cfg := new(DecoderConfig)
	cfg.HostnameKeys = []string{
		"hostname",
		"source",
	}
	cfg.TimestampKeys = []string{
		"timestamp",
	}
	cfg.LevelKeys = []string{
		"level",
	}
	cfg.TextKeys = []string{
		"message",
	}
	cfg.SkipFieldKeys = []string{
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
	}
	cfg.FieldKeys = []string{}

	return cfg
}

func NewDecoder(cfg *DecoderConfig) *Decoder {
	dec := new(Decoder)
	dec.cfg = cfg

	dec.skipFieldKeys = map[string]struct{}{}
	for _, key := range dec.cfg.SkipFieldKeys {
		dec.skipFieldKeys[key] = struct{}{}
	}

	return dec
}

func (dec *Decoder) Hostname(msg map[string]any) string {
	hostname := ""
	for _, key := range dec.cfg.HostnameKeys {
		if v, has := msg[key]; has {
			hostname = v.(string)
			break
		}
	}

	return hostname
}

func (dec *Decoder) Timestamp(msg map[string]any) (ts time.Time) {
	timestamp := ""
	for _, key := range dec.cfg.TimestampKeys {
		if v, has := msg[key]; has {
			timestamp = v.(string)
			break
		}
	}

	if timestamp != "" {
		ts, _ = timeutil.Parse(timestamp)
	}

	return ts
}

func (dec *Decoder) Level(msg map[string]any) Level {
	level := LevelUnkown
	for _, key := range dec.cfg.LevelKeys {
		if v, has := msg[key]; has {
			switch v.(type) {
			case float64:
				level = Level(v.(float64))
			case string:
				if i, err := strconv.ParseUint(v.(string), 10, 64); err == nil {
					level = Level(i)
				}
			}

			break
		}
	}

	return level
}

func (dec *Decoder) Text(msg map[string]any) string {
	message := "-----"
	for _, key := range dec.cfg.TextKeys {
		if v, has := msg[key]; has {
			message = v.(string)
			break
		}
	}

	return message
}

func (dec *Decoder) Fields(msg map[string]any) (keys []string, values []string) {
	// find keys
	if len(dec.cfg.FieldKeys) != 0 {
		for _, key := range dec.cfg.FieldKeys {
			if _, has := msg[key]; !has {
				continue
			}
			keys = append(keys, key)
		}
	} else {
		for key := range msg {
			if _, has := dec.skipFieldKeys[key]; has {
				continue
			}
			keys = append(keys, key)
		}
		sort.Sort(sort.StringSlice(keys))
	}

	// find values
	for _, key := range keys {
		value, has := msg[key]
		if !has {
			continue
		}

		switch value.(type) {
		case float64:
			s := fmt.Sprintf("%f", value)
			s = strings.TrimRight(s, "0")
			s = strings.TrimRight(s, ".")
			values = append(values, s)
		default:
			values = append(values, fmt.Sprintf("%v", value))
		}
	}

	return keys, values
}
