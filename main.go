package main

import (
	"fabric-edgenode/NodeUtils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var peertopics = []string{"register", "upload", "filereq", "KeyUpload", "ReceiveKeyUpload", "ReceiveKeyReq", "DataForwarding", "ReceiveFileRequestFromCenter"}

// var peer0_org1 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9092",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:7984",
// 	PeerNodeName: "peer0.org1.example.com:7051",
// 	OrgID:        "1",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }
// var peer1_org1 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9093",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:8984",
// 	PeerNodeName: "peer1.org1.example.com:8051",
// 	OrgID:        "1",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }
// var peer0_org2 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9094",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:9984",
// 	PeerNodeName: "peer0.org2.example.com:9051",
// 	OrgID:        "2",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }
// var peer1_org2 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9095",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:10984",
// 	PeerNodeName: "peer1.org2.example.com:10051",
// 	OrgID:        "2",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }
// var peer0_org3 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9096",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:11984",
// 	PeerNodeName: "peer0.org3.example.com:11051",
// 	OrgID:        "3",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }
// var peer1_org3 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9097",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:12984",
// 	PeerNodeName: "peer1.org3.example.com:12051",
// 	OrgID:        "3",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org3.example.com/peers/peer1.org3.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }
// var peer0_org4 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9098",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:13984",
// 	PeerNodeName: "peer0.org4.example.com:13051",
// 	OrgID:        "4",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org4.example.com/peers/peer0.org4.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }
// var peer1_org4 = NodeUtils.Nodestructure{
// 	KafkaIp:      "0.0.0.0:9099",
// 	Couchdb_addr: "http://admin:123456@0.0.0.0:14984",
// 	PeerNodeName: "peer1.org4.example.com:14051",
// 	OrgID:        "4",
// 	KeyPath:      "./fixtures/crypto-config/peerOrganizations/org4.example.com/peers/peer0.org4.example.com/msp/keystore/priv_sk",
// 	ConfigPath:   "./cfg/org1conf.yaml",
// }

func main() {
	// Init the peer node
	// NodeUtils.InitPeerNode(peertopics, peer0_org1)
	// NodeUtils.InitPeerNode(peertopics, peer0_org2)
	// NodeUtils.InitPeerNode(peertopics, peer1_org1)
	// NodeUtils.InitPeerNode(peertopics, peer1_org2)
	// NodeUtils.InitPeerNode(peertopics, peer0_org3)
	// NodeUtils.InitPeerNode(peertopics, peer0_org4)
	// NodeUtils.InitPeerNode(peertopics, peer1_org3)
	// NodeUtils.InitPeerNode(peertopics, peer1_org4)
	// var node = NodeUtils.Nodestructure{
	// 	KafkaIp:      os.Getenv("KAFKA_IP"),
	// 	Couchdb_addr: os.Getenv("COUCHDB_ADDR"),
	// 	PeerNodeName: os.Getenv("PEER_NODE_NAME"),
	// 	OrgID:        os.Getenv("ORG_ID"),
	// 	KeyPath:      os.Getenv("KEY_PATH"),
	// }
	//启动websocket
	go NodeUtils.InitWebsocket()

	//gin router
	r := gin.Default()
	r.POST("/register", NodeUtils.Register)       //http://10.0.0.144:8083/register
	r.POST("/upload", NodeUtils.Upload)           //http://10.0.0.144:8083/upload
	r.POST("/requestfile", NodeUtils.Filerequest) //http://10.0.0.144:8083/requestfile
	r.Run("0.0.0.0:8083")                         //在 0.0.0.0:8083 上启动服务
	time.Sleep(time.Second * 2)
	fmt.Println("start")
}
