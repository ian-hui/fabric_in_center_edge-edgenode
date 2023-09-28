package NodeUtils

import (
	"context"
	"fabric-edgenode/clients"
	"fmt"
	"log"
	"sync"
	"time"

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
func (nodestru Nodestructure) InitPeerNode(topics []string) {
	//init kafka producer
	clients.InitProducer(nodestru.KafkaIp)

	//initpeer
	// clients.InitPeerSdk(nodestru.PeerNodeName, nodestru.OrgID, nodestru.ConfigPath)

	//create db in couchdb
	if err := nodestru.Create_cipherkey_info(); err != nil {
		fmt.Println("create cipherkey_info db error:", err)
	}
	if err := nodestru.Create_ciphertext_info(); err != nil {
		fmt.Println("create ciphertext_info db error:", err)
	}

	var wg sync.WaitGroup
	wg.Add(9)
	//创建consumer
	go consumeRegister(nodestru, &wg)
	go consumeUpload(nodestru, &wg)
	go consumeFileReq(nodestru, &wg)
	// consumer2, err := clients.InitConsumer(nodestru.KafkaIp)
	// if err != nil {
	// 	fmt.Printf("fail to start err:%v\n", err)
	// 	return
	// }
	// fmt.Println(nodestru.KafkaIp, "init peer-consumer1 begin")
	go consumeKeyUpload(nodestru, &wg)
	go consumeReceiveKeyUpload(nodestru, &wg)
	go consumeGroupChoose(nodestru, &wg)
	// consumer3, err := clients.InitConsumer(nodestru.KafkaIp)
	// if err != nil {
	// 	fmt.Printf("fail to start err:%v\n", err)
	// 	return
	// }
	fmt.Println(nodestru.KafkaIp, "init peer-consumer1 begin")
	go consumeReceiveKeyReq(nodestru, &wg)
	go consumeDataForwarding(nodestru, &wg)
	go consumeReceiveFileRequestFromCenter(nodestru, &wg)
	wg.Wait()
}

func consumeRegister(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "register")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "register", err)
		}
		err = register(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerRegister error:", err)
		}
	}

}

func consumeUpload(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "upload")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "upload", err)
		}
		err = upload(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerUpload error:", err)
		}
	}
}

func consumeFileReq(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "filereq")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "filereq", err)
		}
		err = filerequest(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerFileReq error:", err)
		}
	}
}

func consumeKeyUpload(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "KeyUpload")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "KeyUpload", err)
		}
		err = keyUpload(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerKeyUpload error:", err)
		}
	}
}

func consumeReceiveKeyUpload(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "ReceiveKeyUpload")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "ReceiveKeyUpload", err)
		}
		err = receivekeyUpload(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerReceiveKeyUpload error:", err)
		}
	}
}

func consumeReceiveKeyReq(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "ReceiveKeyReq")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "ReceiveKeyReq", err)
		}
		err = receivekeyReq(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerReceiveKeyReq error:", err)
		}
	}
}

func consumeDataForwarding(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "DataForwarding")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "DataForwarding", err)
		}
		err = dataForwarding(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerDataForwarding error:", err)
		}
	}
}

func consumeReceiveFileRequestFromCenter(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "ReceiveFileRequestFromCenter")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "ReceiveFileRequestFromCenter", err)
		}
		err = receiveFileRequestFormCenter(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerReceiveFileRequestFromCenter error:", err)
		}
	}
}

func consumeGroupChoose(nodestru Nodestructure, wg *sync.WaitGroup) {
	consumer := clients.InitConsumer(nodestru.KafkaIp, "GroupChoose")
	wg.Done()
	for {
		msg, err := consumer.FetchMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message from topic %s: %v\n", "GroupChoose", err)
		}
		err = chooseGroup(nodestru, msg.Value)
		if err != nil {
			fmt.Println("consumerGroupChoose error:", err)
		}
	}
}
