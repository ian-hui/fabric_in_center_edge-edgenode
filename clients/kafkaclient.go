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

// func InitProducer(kafka_addr string) error {

// 	// 设置发送消息的超时时间
// 	fmt.Println("producer" + kafka_addr + "init begin")
// 	producerconfig := sarama.NewConfig()
// 	producerconfig.Producer.RequiredAcks = 1        // 发送完数据需要leader和follow都确认
// 	producerconfig.Producer.Return.Successes = true // 成功交付的消息将在success channel返回
// 	producerconfig.Producer.Return.Errors = true    // 发送失败的消息将在error channel返回
// 	producer, err := sarama.NewAsyncProducer([]string{kafka_addr}, producerconfig)
// 	if err != nil {
// 		return fmt.Errorf("init kafka producer error: %v", err)
// 	}
// 	fmt.Println("producer" + kafka_addr + "init success")
// 	producerConns = &producer
// 	return nil
// }

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

func InitConsumer(kafka_addr string) *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers:  []string{kafka_addr},
		MinBytes: 10e3, // 最小读取字节数
		MaxBytes: 10e6, // 最大读取字节数
	}

	// 创建Kafka Reader对象
	reader := kafka.NewReader(config)

	return reader
}
