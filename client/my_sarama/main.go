package main

import (
	"fmt"
	"my_sarama/conf"
	"my_sarama/consumer"
	"my_sarama/producer"
)

func main() {
	server := &producer.Service{}
	server.SyncProducerCollector = server.NewSyncProducerCollector(conf.Brokers)
	partition, offset, err := server.SendSyncProducer(conf.TopicA, conf.Key, &producer.Entry{Data: "This is a message."})
	if err != nil {
		panic(err)
	}
	fmt.Printf("partition: %v. offset: %v\n", partition, offset)
	consumer.StartConsumer(conf.Verbose, conf.Oldest, conf.Version, conf.Assignor, conf.GroupA, conf.Topics, conf.Brokers)
}
