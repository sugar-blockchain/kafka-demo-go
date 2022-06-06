package producer

import (
	"fmt"
	"my_sarama/conf"
	"testing"
)

func TestSyncProducer(t *testing.T) {
	server := &Service{}
	server.SyncProducerCollector = server.NewSyncProducerCollector(conf.Brokers)
	partition, offset, err := server.SendSyncProducer(conf.TopicA, conf.Key, &Entry{Data: "This is a message."})
	if err != nil {
		panic(err)
	}
	fmt.Printf("partition: %v. offset: %v\n", partition, offset)
}

func TestAsyncProducer(t *testing.T) {
	server := &Service{}
	server.AsyncProducerCollector = server.NewAsyncProducerCollector(conf.Brokers)
	server.SendAsyncProducer(conf.TopicA, conf.Key, &Entry{Data: "This is a message."})
	fmt.Printf("TestAsyncProducer")
}
