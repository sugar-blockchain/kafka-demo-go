package producer

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

type TlsConfiguration struct {
	certFile, keyFile, caFile string
	verifySsl                 bool
}

type Service struct {
	SyncProducerCollector  sarama.SyncProducer
	AsyncProducerCollector sarama.AsyncProducer
}

func createTlsConfiguration(certFile, keyFile, caFile string, verifySsl bool) (t *tls.Config) {
	if certFile != "" && keyFile != "" && caFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := os.ReadFile(caFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: verifySsl,
		}
	}
	// will be nil by default if nothing is provided
	return t
}

// SaramaLogEnable 开启sarama日志
func (s *Service) SaramaLogEnable() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
}

func (s *Service) NewSyncProducerCollector(brokerList []string) sarama.SyncProducer {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	tlsConfig := createTlsConfiguration("", "", "", false)
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func (s *Service) NewAsyncProducerCollector(brokerList []string) sarama.AsyncProducer {
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	tlsConfig := createTlsConfiguration("", "", "", false)
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	return producer
}

func (s *Service) Close() error {
	if err := s.SyncProducerCollector.Close(); err != nil {
		log.Println("Failed to shut down data collector cleanly", err)
	}

	if err := s.AsyncProducerCollector.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
	}
	return nil
}

// SendSyncProducer 发送同步消息
func (s *Service) SendSyncProducer(topic, key string, entry *Entry) (int32, int64, error) {
	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.
	partition, offset, err := s.SyncProducerCollector.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: entry,
	})
	if err != nil {
		return 0, 0, err
	} else {
		return partition, offset, nil
	}
}

// SendAsyncProducer 发送异步消息
func (s *Service) SendAsyncProducer(topic, key string, entry *Entry) {
	// We will use the client's IP address as key. This will cause
	// all the access log entries of the same IP address to end up
	// on the same partition.
	s.AsyncProducerCollector.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: entry,
	}
}

type Entry struct {
	Data interface{}

	encoded []byte
	err     error
}

func (ale *Entry) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale.Data)
	}
}

func (ale *Entry) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *Entry) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}
