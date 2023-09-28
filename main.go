package main

import (
	"fabric-edgenode/NodeUtils"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var peertopics = []string{"register", "upload", "filereq", "KeyUpload", "ReceiveKeyUpload", "ReceiveKeyReq", "DataForwarding", "ReceiveFileRequestFromCenter"}

func main() {

	var node = NodeUtils.Nodestructure{
		KafkaIp:      os.Getenv("KAFKA_IP"),
		Couchdb_addr: os.Getenv("COUCHDB_ADDR"),
		PeerNodeName: os.Getenv("PEER_NODE_NAME"),
		OrgID:        os.Getenv("ORG_ID"),
		KeyPath:      os.Getenv("KEY_PATH"),
		CenterAddr:   os.Getenv("CENTER_ADDR"),
		ConfigPath:   "/conf/config.yaml",
	}
	
	//init node
	node.InitPeerNode(peertopics)
	//start websocket
	go NodeUtils.InitWebsocket()

	//gin router
	r := gin.Default()
	r.POST("/register", NodeUtils.Register)       //http://10.0.0.144:8083/register
	r.POST("/upload", NodeUtils.Upload)           //http://10.0.0.144:8083/upload
	r.POST("/requestfile", NodeUtils.Filerequest) //http://10.0.0.144:8083/requestfile
	r.Run("0.0.0.0:8083")                         // 0.0.0.0:8083
	fmt.Println("start")
}
