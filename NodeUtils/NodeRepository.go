package NodeUtils

import (
	"context"
	"encoding/json"
	"fabric-edgenode/clients"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/cloudflare/cfssl/log"
	_ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v4"
)

// func GetAreaKafkaAddrInZookeeper(area_num string) []string {
// 	var res map[string]interface{}
// 	kafka_addr := make([]string, 0)
// 	conn := clients.ZookeeperConns
// 	children, _, err := conn.Children("/peer/area" + area_num)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, child := range children {
// 		endfilenames, _, err := conn.Children("/peer/area" + area_num + "/" + child + "/brokers/ids")
// 		if err != nil {
// 			panic(err)
// 		}
// 		for _, fname := range endfilenames {
// 			data, _, err := conn.Get("/peer/area" + area_num + "/" + child + "/brokers/ids/" + fname)
// 			if err != nil {
// 				panic(err)
// 			}
// 			json.Unmarshal(data, &res)
// 			temp := res["port"].(float64)
// 			port := strconv.FormatFloat(temp, 'f', 0, 64)
// 			kafka_addr = append(kafka_addr, "0.0.0.0:"+port)
// 			// fmt.Printf("Kafka node %s is running at %s\n", child, data)
// 		}

// 	}
// 	return kafka_addr
// }

// func GetAllPeerAddrInZookeeper() []string {
// 	var res map[string]interface{}
// 	kafka_addr := make([]string, 0)
// 	conn := clients.ZookeeperConns
// 	areas, _, err := conn.Children("/peer")
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, area := range areas {
// 		endpointnames, _, err := conn.Children("/peer/" + area)
// 		if err != nil {
// 			panic(err)
// 		}
// 		for _, endpointname := range endpointnames {
// 			endfilenames, _, err := conn.Children("/peer/" + area + "/" + endpointname + "/brokers/ids")
// 			if err != nil {
// 				panic(err)
// 			}
// 			for _, fname := range endfilenames {
// 				data, _, err := conn.Get("/peer/" + area + "/" + endpointname + "/brokers/ids/" + fname)
// 				if err != nil {
// 					panic(err)
// 				}
// 				json.Unmarshal(data, &res)
// 				temp := res["port"].(float64)
// 				port := strconv.FormatFloat(temp, 'f', 0, 64)
// 				kafka_addr = append(kafka_addr, "0.0.0.0:"+port)
// 				// fmt.Printf("Kafka node %s is running at %s\n", child, data)
// 			}
// 		}

// 	}
// 	return kafka_addr
// }

func ProducerAsyncSending(messages []byte, topic string, kafka_addr string) error {
	kafka_client := clients.GetProducer(kafka_addr)
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	err := AsyncSend2Kafka(kafka_client, msg, messages)
	if err != nil {
		return fmt.Errorf("AsyncSend2Kafka error: %v", err)
	}
	return nil
}

// 异步发函数函数
func AsyncSend2Kafka(client sarama.AsyncProducer, msg *sarama.ProducerMessage, content []byte) error {
	msg.Value = sarama.StringEncoder(content)

	// 发送消息
	client.Input() <- msg
	go func() {
		select {
		case success := <-client.Successes():
			fmt.Println("message sent successfully")
			fmt.Printf("partition:%v offset:%v\n", success.Partition, success.Offset)
		case err := <-client.Errors():
			panic(err)
		}
	}()
	return nil
}

func DeleteTargetInArrayStr(array_str []string, target string) []string {
	j := 0
	for _, val := range array_str {
		if val != target {
			array_str[j] = val
			j++
		}
	}
	return array_str[:j]
}

func IsIdentical(str1 []string, str2 []string) (t bool) {
	t = false
	if len(str1) == 0 || len(str2) == 0 {
		return
	}
	map1, map2 := make(map[string]int), make(map[string]int)
	for i := 0; i < len(str1); i++ {
		map1[str1[i]] = i
	}
	for i := 0; i < len(str2); i++ {
		map2[str2[i]] = i
	}
	for k := range map1 {
		if _, ok := map2[k]; ok {
			t = true
		}
	}
	return
}

// func NewTLSConfig() (*tls.Config, error) {
// 	// Load client cert
// 	cert, err := tls.LoadX509KeyPair("./kafka_crypto/client.cer.pem", "./kafka_crypto/client.key.pem")
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Load CA cert
// 	caCert, err := ioutil.ReadFile("./kafka_crypto/server.cer.pem")
// 	if err != nil {
// 		return nil, err
// 	}
// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)
// 	tlsConfig := &tls.Config{
// 		RootCAs:      caCertPool,
// 		Certificates: []tls.Certificate{cert},
// 	}
// 	return tlsConfig, err
// }

// Functions that transfer data at regular intervals
func CipherPosTransfer(FileId string, nodestru Nodestructure, addr string) error {
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	client.Mu.Lock()
	defer client.Mu.Unlock()
	db := client.C.DB("ciphertext_info", nil) //connect to position_info
	resultSet := db.Get(context.TODO(), FileId)
	//deal with error
	if resultSet.Err() != nil {
		//not found
		if kivik.StatusCode(resultSet.Err()) == 404 {
			fmt.Println("<------file ", FileId, " Not found in this node------>")
		} else {
			fmt.Println("db.get error:", resultSet.Err())
		}
		return resultSet.Err()
	}
	//get ciphertext info
	var fileinfotempdto FileRequestDTO
	err = resultSet.ScanDoc(&fileinfotempdto)
	if err != nil {
		fmt.Println("scandoc error: ", err)
		return err
	}
	//ciphertext info
	fileinfostru := FileInfo{
		FileId:     fileinfotempdto.FileId,
		Ciphertext: fileinfotempdto.Ciphertext,
	}
	res, err := json.Marshal(fileinfostru)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}
	topic := "UploadCiText" //操作名
	err = ProducerAsyncSending(res, topic, addr)
	if err != nil {
		return fmt.Errorf("ProducerAsyncSending error: %v", err)
	}

	//delete local ciphertext
	_, err = db.Delete(context.TODO(), FileId, resultSet.Rev())
	if err != nil {
		return fmt.Errorf("delete local ciphertext error: %v", err)
	}
	return nil
}

func consumerTopic(consumer sarama.Consumer, nodestru Nodestructure, wg *sync.WaitGroup, TopicName string, f func(nodestru Nodestructure, msg []byte) (err error), errorMessage string) {
	partitonConsumer, err := consumer.ConsumePartition(TopicName, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	defer wg.Done()
	for {
		select {
		case msg := <-partitonConsumer.Messages():
			err := f(nodestru, msg.Value)
			if err != nil {
				log.Errorf("f error: %v", err)
			}
		case <-partitonConsumer.Errors():
			fmt.Println(errorMessage)
		}
	}
}
