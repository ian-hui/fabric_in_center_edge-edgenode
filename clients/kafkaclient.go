package clients

import (
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	producerConns *kafka.Writer
)

func InitProducer(kafka_addr string) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafka_addr},
		Async:   true,
	})
	producerConns = writer
}

func GetProducer(kafka_addr string) *kafka.Writer {
	if producerConns == nil {
		mu := sync.Mutex{}
		mu.Lock()
		defer mu.Unlock()
		if producerConns == nil {
			InitProducer(kafka_addr)
		}
	}
	return producerConns
}

func InitConsumer(kafka_addr string, topic string) *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers:  []string{kafka_addr},
		Topic:    topic,
		MinBytes: 10e3, // 最小读取字节数
		MaxBytes: 10e6, // 最大读取字节数
	}

	// 创建Kafka Reader对象
	reader := kafka.NewReader(config)

	return reader
}
