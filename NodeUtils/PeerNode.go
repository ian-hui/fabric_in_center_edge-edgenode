package NodeUtils

import (
	"fabric-edgenode/clients"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
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
	if err := clients.InitPeerSdk(nodestru.PeerNodeName, nodestru.OrgID, nodestru.ConfigPath); err != nil {
		fmt.Println("init peer sdk error:", err)
	}
	//create db in couchdb
	c, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		fmt.Println("get couchdb client error:", err)
	}
	if err := c.Create_ciphertext_info(); err != nil {
		fmt.Println("create position_info db error:", err)
	}
	if err := c.Create_cipherkey_info(); err != nil {
		fmt.Println("create cipherkey_info db error:", err)
	}
	var wg sync.WaitGroup
	wg.Add(8)
	consumer1, err := clients.InitConsumer(nodestru.KafkaIp)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	//创建consumer
	go consumerTopic(consumer1, nodestru, &wg, "register", register, "register error")
	go consumerTopic(consumer1, nodestru, &wg, "upload", upload, "upload error")
	go consumerTopic(consumer1, nodestru, &wg, "filereq", filerequest, "filereq error")
	consumer2, err := clients.InitConsumer(nodestru.KafkaIp)
	if err != nil {
		fmt.Printf("fail to start err:%v\n", err)
		return
	}
	go consumerTopic(consumer2, nodestru, &wg, "KeyUpload", keyUpload, "KeyUpload error")
	go consumerTopic(consumer2, nodestru, &wg, "ReceiveKeyUpload", receivekeyUpload, "ReceiveKeyUpload error")
	// go consumeGroupChoose(nodestru, &wg)
	consumer3, err := clients.InitConsumer(nodestru.KafkaIp)
	if err != nil {
		fmt.Printf("fail to start err:%v\n", err)
		return
	}
	go consumerTopic(consumer3, nodestru, &wg, "ReceiveKeyReq", receivekeyReq, "ReceiveKeyReq error")
	go consumerTopic(consumer3, nodestru, &wg, "DataForwarding", dataForwarding, "DataForwarding error")
	go consumerTopic(consumer3, nodestru, &wg, "ReceiveFileRequestFromCenter", receiveFileRequestFormCenter, "ReceiveFileRequestFromCenter error")
	wg.Wait()
}

// func consumeGroupChoose(nodestru Nodestructure, wg *sync.WaitGroup) {
// 	consumer := clients.InitConsumer(nodestru.KafkaIp, "GroupChoose")
// 	wg.Done()
// 	for {
// 		msg, err := consumer.FetchMessage(context.Background())
// 		if err != nil {
// 			log.Printf("failed to read message from topic %s: %v\n", "GroupChoose", err)
// 			return
// 		}
// 		err = chooseGroup(nodestru, msg.Value)
// 		if err != nil {
// 			fmt.Println("consumerGroupChoose error:", err)
// 		}
// 	}
// }

func GetNodeLoadService() (float64, error) {
	cmd := exec.Command("ps", "-p", "1", "-o", "%cpu")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("output error:", err)
		return 0, err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		fmt.Println("Invalid output")
		return 0, err
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 1 {
		fmt.Println("Invalid output")
		return 0, err
	}
	f, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		fmt.Println("Invalid output")
		return 0, err
	}
	return f, nil

}
