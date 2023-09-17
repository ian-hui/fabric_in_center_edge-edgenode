package NodeUtils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fabric-edgenode/clients"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/fatih/structs"
	_ "github.com/go-kivik/couchdb/v4"   // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v4" //couchdb-go第三方库
)

func Create_ciphertext_info(nodestru Nodestructure) error { //create db in couchdb
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	err = client.C.CreateDB(context.TODO(), "ciphertext_info", nil)
	if err != nil {
		return fmt.Errorf("create ciphertext_info db error: %v", err)
	}
	fmt.Println("<-------", nodestru.Couchdb_addr, " couchdb ciphertext_info created!------>")
	return nil
}

func Create_cipherkey_info(nodestru Nodestructure) error { //create db in couchdb
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}

	err = client.C.CreateDB(context.TODO(), "cipherkey_info", nil)
	if err != nil {
		return fmt.Errorf("create ciphertext_info db error: %v", err)
	}
	fmt.Println("<-------", nodestru.Couchdb_addr, " couchdb cipherkey_info created!------>")
	return nil
}

func Create_position_info(nodestru Nodestructure) error { //create db in couchdb
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	err = client.C.CreateDB(context.TODO(), "position_info", nil)
	if err != nil {
		return fmt.Errorf("create ciphertext_info db error: %v", err)
	}
	fmt.Println("<-------", nodestru.Couchdb_addr, " couchdb position_info created!------>")
	return nil
}

// UploadOrUpdatePostion as the name means
func UploadPostion(f PositionInfo, nodestru Nodestructure) error {
	json_fileposif := structs.Map(&f) //转格式，详细看https://github.com/go-kivik/kivik
	//连接couchdb
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	client.Mu.Lock()
	defer client.Mu.Unlock()
	db := client.C.DB("position_info", nil)                   //连接couchdb中的positon_info数据库
	_, err = db.Put(context.TODO(), f.FileId, json_fileposif) //把数据info上传到db
	if err != nil {
		panic(err)
	}
	return nil
}

func UpdatePostion(fileid string, nodestru Nodestructure, targetdata []string, changedata []string) error {
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	// client.Mu.Lock()
	// defer client.Mu.Unlock()
	db := client.C.DB("position_info", nil) //connect to position_info
	resultSet := db.Get(context.TODO(), fileid)

	var doc map[string]interface{}
	err = resultSet.ScanDoc(&doc)
	if err != nil {
		return fmt.Errorf("scandoc error: ", err)
	}
	for i := range targetdata {
		//modify the position in center
		doc[targetdata[i]] = changedata[i]
	}
	_, err = db.Put(context.TODO(), fileid, doc) //把数据info上传到db
	if err != nil {
		UpdatePostion(fileid, nodestru, targetdata, changedata)
	}
	return nil
}

func UploadOrUpdateCipherText(f FileInfo, nodestru Nodestructure) error {
	json_fileposif := structs.Map(&f) //转格式，详细看https://github.com/go-kivik/kivik
	//连接couchdb
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	client.Mu.Lock()
	defer client.Mu.Unlock()
	db := client.C.DB("ciphertext_info", nil)                    //连接couchdb中的cipher_info数据库
	rev, err := db.Put(context.TODO(), f.FileId, json_fileposif) //把数据info上传到db
	if err != nil {
		// if kivik.StatusCode(err) == 409
		UploadOrUpdateCipherText(f, nodestru)
	}
	fmt.Printf("%s inserted with revision %s\n", f.FileId, rev)
	return nil
}

func Getinfo(fileid string, nodestru Nodestructure, dbname string) (kivik.ResultSet, error) {
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return nil, fmt.Errorf("get couchdb client error: %v", err)
	}
	db := client.C.DB(dbname, nil) //connect to position_info
	resultSet := db.Get(context.TODO(), fileid)
	if resultSet.Err() != nil {
		return nil, resultSet.Err()
	}
	return resultSet, nil

}

func UploadCipherKey(f KeyDetailInfo, nodestru Nodestructure) error {
	json_fileposif := structs.Map(&f) //转格式，详细看https://github.com/go-kivik/kivik
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	client.Mu.Lock()
	defer client.Mu.Unlock()
	db := client.C.DB("cipherkey_info", nil)                  //connect to ciphertext_info
	_, err = db.Put(context.TODO(), f.FileId, json_fileposif) //把数据info上传到db
	if err != nil {
		return fmt.Errorf("upload cipherkey error: %v", err)
	}
	// fmt.Printf("%s inserted with revision %s\n", f.FileId, rev)
	return nil
}

func GetAreaKafkaAddrInZookeeper(area_num string) []string {
	var res map[string]interface{}
	kafka_addr := make([]string, 0)
	conn := clients.ZookeeperConns
	children, _, err := conn.Children("/peer/area" + area_num)
	if err != nil {
		panic(err)
	}

	for _, child := range children {
		endfilenames, _, err := conn.Children("/peer/area" + area_num + "/" + child + "/brokers/ids")
		if err != nil {
			panic(err)
		}
		for _, fname := range endfilenames {
			data, _, err := conn.Get("/peer/area" + area_num + "/" + child + "/brokers/ids/" + fname)
			if err != nil {
				panic(err)
			}
			json.Unmarshal(data, &res)
			temp := res["port"].(float64)
			port := strconv.FormatFloat(temp, 'f', 0, 64)
			kafka_addr = append(kafka_addr, "0.0.0.0:"+port)
			// fmt.Printf("Kafka node %s is running at %s\n", child, data)
		}

	}
	return kafka_addr
}

func GetAllPeerAddrInZookeeper() []string {
	var res map[string]interface{}
	kafka_addr := make([]string, 0)
	conn := clients.ZookeeperConns
	areas, _, err := conn.Children("/peer")
	if err != nil {
		panic(err)
	}
	for _, area := range areas {
		endpointnames, _, err := conn.Children("/peer/" + area)
		if err != nil {
			panic(err)
		}
		for _, endpointname := range endpointnames {
			endfilenames, _, err := conn.Children("/peer/" + area + "/" + endpointname + "/brokers/ids")
			if err != nil {
				panic(err)
			}
			for _, fname := range endfilenames {
				data, _, err := conn.Get("/peer/" + area + "/" + endpointname + "/brokers/ids/" + fname)
				if err != nil {
					panic(err)
				}
				json.Unmarshal(data, &res)
				temp := res["port"].(float64)
				port := strconv.FormatFloat(temp, 'f', 0, 64)
				kafka_addr = append(kafka_addr, "0.0.0.0:"+port)
				// fmt.Printf("Kafka node %s is running at %s\n", child, data)
			}
		}

	}
	return kafka_addr
}

func GetCenterKafkaAddrInZookeeper() []string {
	var res map[string]interface{}
	center_kafka_addr := make([]string, 0, 4)
	conn := clients.ZookeeperConns
	children, _, err := conn.Children("/center/brokers/ids")
	if err != nil {
		panic(err)
	}

	for _, child := range children {
		data, _, err := conn.Get("/center/brokers/ids/" + child)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(data, &res)
		temp := res["port"].(float64)
		port := strconv.FormatFloat(temp, 'f', 0, 64)
		center_kafka_addr = append(center_kafka_addr, "0.0.0.0:"+port)
	}
	return center_kafka_addr
}

func CheckNotExistence(f string, nodestru Nodestructure, dbname string) bool {
	_, err := Getinfo(f, nodestru, dbname)
	if err != nil {
		if kivik.StatusCode(err) == 404 {
			return true
		}
		fmt.Println("check file existence getinfo error ", err)
	}
	return false
}

func ProducerAsyncSending(messages []byte, topic string, kafka_addr string) error {
	kafka_client, err := clients.GetProducer(kafka_addr)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	err = AsyncSend2Kafka(kafka_client, msg, messages)
	if err != nil {
		return err
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
	for k, _ := range map1 {
		if _, ok := map2[k]; ok {
			t = true
		}
	}
	return
}

func NewTLSConfig() (*tls.Config, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair("./kafka_crypto/client.cer.pem", "./kafka_crypto/client.key.pem")
	if err != nil {
		return nil, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile("./kafka_crypto/server.cer.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}
	return tlsConfig, err
}
