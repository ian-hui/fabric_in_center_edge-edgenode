package NodeUtils

import (
	"fabric-edgenode/clients"
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	_ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
	//couchdb-go第三方库
)

// var startTime time.Time
var transmitTime = time.Minute * 30

// var ProducerConn = &kafka_producer_map{
// 	&sync.Map{},
// }
// var CouchDBConn = &kivik_client_map{
// 	&sync.Map{},
// }

// 初始化节点consumer producer
// topic：register，upload，filereq,KeyUpload
func InitPeerNode(topics []string, nodestru Nodestructure) {
	//创建ciphertext和cipherkey数据库
	clients.InitCouchdb(nodestru.Couchdb_addr)
	//init kafka producer
	clients.InitProducer(nodestru.KafkaIp)
	//init peer sdk
	clients.InitPeerSdk(nodestru.PeerNodeName, nodestru.OrgID, "./conf/config.yaml")
	//create db in couchdb
	Create_cipherkey_info(nodestru)
	Create_ciphertext_info(nodestru)
	var wg sync.WaitGroup
	wg.Add(9)
	//创建consumer
	consumer1, err := clients.InitConsumer(nodestru.KafkaIp)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	fmt.Println(nodestru.KafkaIp, "init peer-consumer1 begin")
	go consumeRegister(consumer1, nodestru, &wg)
	go consumeUpload(consumer1, nodestru, &wg)
	go consumeFileReq(consumer1, nodestru, &wg)
	consumer2, err := clients.InitConsumer(nodestru.KafkaIp)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	fmt.Println(nodestru.KafkaIp, "init peer-consumer1 begin")
	go consumeKeyUpload(consumer2, nodestru, &wg)
	go consumeReceiveKeyUpload(consumer2, nodestru, &wg)
	go consumeGroupChoose(consumer2, nodestru, &wg)
	consumer3, err := clients.InitConsumer(nodestru.KafkaIp)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	fmt.Println(nodestru.KafkaIp, "init peer-consumer1 begin")
	go consumeReceiveKeyReq(consumer3, nodestru, &wg)
	go consumeDataForwarding(consumer3, nodestru, &wg)
	go consumeReceiveFileRequestFromCenter(consumer3, nodestru, &wg)
	wg.Wait()
}

func consumeRegister(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("register", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := register(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerRegister error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerRegister error")
		}
	}
}

func consumeUpload(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("upload", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := upload(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerUpload error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerUpload error")
		}
	}
}

func consumeFileReq(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("filereq", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := filerequest(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerFileReq error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerFileReq error")
		}
	}
}

func consumeKeyUpload(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("KeyUpload", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := keyUpload(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerKeyUpload error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerKeyUpload error")
		}
	}
}

func consumeReceiveKeyUpload(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("ReceiveKeyUpload", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := receivekeyUpload(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerReceiveKeyUpload error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerReceiveKeyUpload error")
		}
	}
}

func consumeReceiveKeyReq(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("ReceiveKeyReq", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := receivekeyReq(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerReceiveKeyReq error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerReceiveKeyReq error")
		}
	}
}

func consumeDataForwarding(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("DataForwarding", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := dataForwarding(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerDataForwarding error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerDataForwarding error")
		}
	}
}

func consumeReceiveFileRequestFromCenter(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("ReceiveFileRequestFromCenter", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := receiveFileRequestFormCenter(nodestru, msg.Value)
			if err != nil {
				fmt.Println("consumerReceiveFileRequestFromCenter error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerReceiveFileRequestFromCenter error")
		}
	}
}

func consumeGroupChoose(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup) {
	partitonConsumer, err := consumer.ConsumePartition("GroupChoose", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := chooseGroup(nodestru, msg.Value)
			if err != nil {
				fmt.Println("chooseGroup error:", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println("consumerReceiveFileRequestFromCenter error")
		}
	}
}
