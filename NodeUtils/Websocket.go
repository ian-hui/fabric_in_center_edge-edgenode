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
	welcomeMessage := "请选择:1.上传密钥2.请求密文密钥"
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
	if len(msg) == 0 {
		log.Println("failed to readmessage : ", err)
		return
	}
	Option := string(msg)
	if Option == "1" {
		fmt.Println("websocket start")
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("failed to readmessage : ", err)
			return
		}
		var test_D Test_data
		err = json.Unmarshal(msg, &test_D)
		if err != nil {
			log.Println("failed to unmarshal : ", err)
			return
		}

		Attribute := test_D.Attribute
		User_id := test_D.Username
		user_key_path := "/kafka_crypto/ianhui.private.pem"
		file_id := test_D.Fileid
		upload_infomation := make(map[string]KeyDetailInfo)

		//注册
		userinformation := sdkInit.UserInfo{
			UserId:    test_D.Username,
			Username:  test_D.Username,
			Attribute: test_D.Attribute,
		}
		res, err := json.Marshal(userinformation)
		if err != nil {
			fmt.Printf("fail to Serialization, err:%v\n", err)
			return
		}
		if err = ProducerAsyncSending(res, "register", os.Getenv("KAFKA_IP")); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		const kb = 100
		file := strings.Repeat("a", 100*kb) // 5MB string
		//构建一个user信息结构
		fileinfomation := FileInfo{
			FileId:     test_D.Fileid,
			Ciphertext: file,
		}

		// 连接kafka
		r, err := json.Marshal(fileinfomation)
		if err != nil {
			fmt.Printf("fail to Serialization, err:%v\n", err)
			return
		}
		if err = ProducerAsyncSending(r, "upload", os.Getenv("KAFKA_IP")); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		//生成测试key
		str := strings.Repeat("a", 16)
		secrets, err := shamir.Split([]byte(str), 8, 3)
		if err != nil {
			fmt.Println(err)
		}

		//test
		for i := range test_D.Keypath {
			pubkey := &sdkInit.GetNodePrivateKey(test_D.Keypath[i]).PublicKey
			encrypted, sign_len := sdkInit.ClientEncryptionByPubECC(user_key_path, pubkey, secrets[i])
			upload_infomation[test_D.Addr[i]] = KeyDetailInfo{
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
		res, err = json.Marshal(keyuploadinfo)
		if err != nil {
			log.Printf("fail to Serialization, err:%v\n", err)
			return
		}
		err = ProducerAsyncSending(res, "KeyUpload", os.Getenv("KAFKA_IP"))
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
		//⌛️计时器
		startTime := time.Now()
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("failed to readmessage : ", err)
			return
		}
		userid := strings.Split(string(msg), ",")[0]
		fileid := strings.Split(string(msg), ",")[1]
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
		err = ProducerAsyncSending(res, "filereq", os.Getenv("KAFKA_IP"))
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
							ws.WriteMessage(websocket.TextMessage, []byte(duration.String()))
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
