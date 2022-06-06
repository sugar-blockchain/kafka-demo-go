package consumer

import (
	"my_sarama/conf"
	"testing"
)

func TestConsumerGroupA1(t *testing.T) {
	StartConsumer(conf.Verbose, conf.Oldest, conf.Version, conf.Assignor, conf.GroupA, conf.Topics, conf.Brokers)
}

func TestConsumerGroupA2(t *testing.T) {
	StartConsumer(conf.Verbose, conf.Oldest, conf.Version, conf.Assignor, conf.GroupA, conf.Topics, conf.Brokers)
}

func TestConsumerGroupB(t *testing.T) {
	StartConsumer(conf.Verbose, conf.Oldest, conf.Version, conf.Assignor, conf.GroupB, conf.Topics, conf.Brokers)
}

func TestConsumerGroupC(t *testing.T) {
	StartConsumer(conf.Verbose, false, conf.Version, conf.Assignor, conf.GroupC, conf.Topics, conf.Brokers)
}
