package NodeUtils

import (
	"fabric-edgenode/sdkInit"
	"os"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/vault/shamir"
)

var (
	resu   int64
	timers sync.Map
	//maintain a connection map (connection,userid)
	connections      sync.Map
	sendDataChannels sync.Map
	upgrader         = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func InitWebsocket() {
	http.HandleFunc("/ws", WebsocketStarter)
	fmt.Println("websocket start")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func WebsocketStarter(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go HandleWebsocket(conn)
}

func HandleWebsocket(ws *websocket.Conn) {

	//welcome message
	welcomeMessage := "请选择:1.上传密钥2.请求密文密钥, 并输入连接的节点地址"
	err := ws.WriteJSON(welcomeMessage)
	if err != nil {
		log.Println(err)
		return
	}
	//first time receive (choose service)
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println("failed to readmessage : ", err)
		return
	}
	choose := strings.Split(string(msg), ",")
	if len(choose) != 2 {
		log.Println("failed to readmessage : ", err)
		return
	}
	Option := choose[0]
	kafka_ip := choose[1]
	if Option == "1" {
		Attribute := Test_data1.Attribute
		User_id := Test_data1.Username
		user_key_path := Test_data1.Userkeypath
		file_id := Test_data1.Fileid

		ws.WriteJSON("开始传输密钥")
		upload_infomation := make(map[string]KeyDetailInfo)

		//生成测试key
		str := strings.Repeat("a", 16)
		secrets, err := shamir.Split([]byte(str), 4, 2)
		if err != nil {
			fmt.Println(err)
		}

		//test
		for i := 0; i < 4; i++ {
			pubkey := &sdkInit.GetNodePrivateKey(Test_data1.Keypath[i]).PublicKey
			encrypted, sign_len := sdkInit.ClientEncryptionByPubECC(user_key_path, pubkey, secrets[i])
			upload_infomation[Test_data1.Ip[i]] = KeyDetailInfo{
				FileId:  file_id,
				Key:     encrypted,
				Signlen: sign_len,
				UserId:  User_id,
			}
		}
		keyuploadinfo := KeyUploadInfo{
			Upload_Infomation: &upload_infomation,
			Attribute:         Attribute,
		}
		res, err := json.Marshal(keyuploadinfo)
		if err != nil {
			log.Printf("fail to Serialization, err:%v\n", err)
			return
		}
		topic := "KeyUpload" //操作名
		err = ProducerAsyncSending(res, topic, kafka_ip)
		if err != nil {
			err := ws.WriteJSON(err)
			if err != nil {
				log.Println("writejson error :", err)
			}
		}
		ws.WriteJSON("SUCCESS")
		return

	}
	//request for file and key
	if Option == "2" {
		startTime := time.Now()
		userid := "ianhui"
		fileid := "1"
		//set userid and bind the websocket conn and userid
		// Create a new channel for the user if it doesn't exist
		if _, ok := sendDataChannels.Load(userid); !ok {
			sendDataChannels.Store(userid, make(chan interface{}, 100))
		}
		// send request to kafka
		FilerequestStruct := FileRequest{
			FileId: fileid,
			UserId: userid,
		}
		res, err := json.Marshal(FilerequestStruct)
		if err != nil {
			fmt.Printf("fail to Serialization, err:%v\n", err)
			return
		}
		topic := "filereq"
		err = ProducerAsyncSending(res, topic, os.Getenv("KAFKA_IP"))
		if err != nil {
			err := ws.WriteJSON(err)
			if err != nil {
				log.Println("writejson error :", err)
			}
			return
		}

		//receive data from kafka
		if conn, ok := sendDataChannels.Load(userid); ok {
			if conn, ok := conn.(chan interface{}); ok {
				temp := 0
				for {
					select {
					case <-conn:
						// fmt.Println(data)
						temp++
						if temp == 4 {
							duration := time.Since(startTime)
							fmt.Println("the duration is ", duration)
							resu = resu + int64(duration/time.Millisecond)
							return
						}
					case <-time.After(10 * time.Minute):
						fmt.Println("timeout")
						return
					}
				}
			}
		}
	}

}

// send data from the node to the client
func SendData(userid string, data interface{}) {
	// Create a new channel for the user if it doesn't exist
	conn, ok := sendDataChannels.Load(userid)
	if !ok {
		sendDataChannels.Store(userid, make(chan interface{}, 100))
	} else {
		// if exist ,send the data to func handleWebsocket()
		conn.(chan interface{}) <- data
		//if the channel not be used in 10min delete it
		// 删除之前的计时器
		timer, ok := timers.Load(userid)
		if ok {
			thetimer := timer.(*time.Timer)
			thetimer.Stop()
			// 创建新的计时器
			timers.Store(userid, time.AfterFunc(time.Duration(10)*time.Minute, func() {
				// 到时间后执行的代码
				sendDataChannels.Delete(userid)
				// 删除计时器
				timers.Delete(userid)
			}))
			defer func() {
				if err := recover(); err != nil {
					// log error and cleanup
					sendDataChannels.Delete(userid)
				}
			}()
		}

	}

}
