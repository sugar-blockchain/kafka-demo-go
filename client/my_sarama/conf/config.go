package conf

import "strings"

var (
	// Common
	Brokers = []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"}
	TopicA  = "Topic1"
	TopicB  = "Topic1"
	Key     = "key1"
	// Consumer
	Oldest   = true
	Version  = "2.8.1"
	GroupA   = "GroupA"
	GroupB   = "GroupB"
	GroupC   = "GroupC"
	Topics   = strings.Join([]string{TopicA, TopicB}, ",")
	Assignor = "range"
	Verbose  = false
)
