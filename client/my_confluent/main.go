package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"my_confluent/consumer"
	"my_confluent/producer"
	"time"
)

func main() {
	for true {
		producer.Producer(&kafka.ConfigMap{"bootstrap.servers": "127.0.0.1"}, "topicA", "keyA", "words")
		time.Sleep(2 * time.Second)
	}
	consumer.Consumer(&kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1",
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"group.id":              "group",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	}, []string{"topicA"})
}
