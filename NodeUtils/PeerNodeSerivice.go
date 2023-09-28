package NodeUtils

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fabric-edgenode/clients"
	"fabric-edgenode/sdkInit"
	"fmt"
	"time"

	kivik "github.com/go-kivik/kivik/v4"
)

func register(nodestru Nodestructure, msg []byte) (err error) {
	var userif sdkInit.UserInfo
	err = json.Unmarshal(msg, &userif)
	if err != nil {
		return fmt.Errorf("register unmarshal error:%v", err)
	}
	// prikeyPath := sdkInit.GenerateKey(userif.UserId)
	pubkey := &sdkInit.GetClientPrivateKey("/kafka_crypto/ianhui.private.pem").PublicKey
	userif.PublicKey = sdkInit.ConversionEcdsaPub2MyPub(pubkey)
	ret, err := clients.GetPeerFabric(nodestru.PeerNodeName, "user").SetUserInfo(userif)
	if err != nil {
		return fmt.Errorf("register set error:%v", err)
	}
	fmt.Println(userif.Username, " set success,the transactionID is ", ret, " and userid is ", userif.UserId)
	return nil
}

func upload(nodestru Nodestructure, msg []byte) (err error) {
	fmt.Println("begin to upload locally")

	//先上传密文到本地
	var fileif FileInfo
	err = json.Unmarshal(msg, &fileif)
	if err != nil {
		return fmt.Errorf("upload unmarshal error:%v", err)
	}
	err = UploadOrUpdateCipherText(fileif, nodestru)
	if err != nil {
		return fmt.Errorf("upload or update ciphertext error:%v", err)
	}

	//然后把信息上传到center
	fmt.Println("<----begin to upload position to center---->")
	//密文位置，密钥位置
	file_postion_info := PositionInfo{
		FileId:   fileif.FileId,
		Position: nodestru.KafkaIp,
		GroupIps: []string{},
	}
	res, err := json.Marshal(file_postion_info)
	if err != nil {
		return fmt.Errorf("marshal error:%v", err)
	}

	topic := "uploadposition" //操作名
	err = ProducerAsyncSending(res, topic, nodestru.CenterAddr)
	if err != nil {
		return fmt.Errorf("producer async sending err:%v", err)
	}

	time2transmit(fileif, nodestru)
	return nil
}

func chooseGroup(nodestru Nodestructure, msg []byte) (err error) {
	//实例化一个groupchooser
	gc := &groupChooser{
		storage_weight:            0.5,
		distance_weight:           0.5,
		standard_deviation_weight: 0.5,
		curNodeInfo:               nodestru.NodeInfo,
	}

	group, delay, sd, err := gc.chooseGroup(6)
	if err != nil {
		return fmt.Errorf("choose group error:%v", err)
	}
	marshaled_group, err := json.Marshal(group)
	if err != nil {
		return fmt.Errorf("marshal error:%v", err)
	}
	SendData(string(msg), marshaled_group)
	return
}

// 上传密钥
func keyUpload(nodestru Nodestructure, msg []byte) (err error) {
	var keyinfostru KeyUploadInfo
	err = json.Unmarshal(msg, &keyinfostru)
	if err != nil {
		return fmt.Errorf("keyupload unmarshal error:%v", err)
	}
	//获取这个组内的所有ip，用于存放在中心节点中作为元数据
	group_position := make([]string, 0, len(*keyinfostru.Upload_Infomation))
	for kafkaip := range *keyinfostru.Upload_Infomation {
		group_position = append(group_position, kafkaip)
	}
	for kafkaip, key_info := range *keyinfostru.Upload_Infomation {
		var (
			err              error
			user_information sdkInit.UserInfo
		)
		if kafkaip == nodestru.KafkaIp {
			//check the key not exist
			check := CheckNotExistence(key_info.FileId, nodestru, "cipherkey_info")
			if !check {
				return fmt.Errorf("keyupload error: key already exist")
			}
			k_access := sdkInit.FileAccessInfo{
				FileId:    key_info.FileId,
				Attribute: keyinfostru.Attribute,
			}
			k_pos := KeyPostionUploadInfo{
				FileId:   key_info.FileId,
				GroupIps: group_position,
			}
			user_information, err = clients.GetPeerFabric(nodestru.PeerNodeName, "user").GetUserInfo(key_info.UserId, nodestru.PeerNodeName)
			if err != nil {
				return fmt.Errorf("keyupload get user info error:%v", err)
			}
			data := sdkInit.DecryptionByPri(nodestru.KeyPath, user_information.PublicKey.ConversionMyPub2EcdsaPub(), []byte(key_info.Key), key_info.Signlen)
			key_info.Key = data
			//upload cipherkey in local peernode
			err = UploadCipherKey(key_info, nodestru)
			if err != nil {
				return fmt.Errorf("keyupload upload cipherkey error:%v", err)
			}
			//set access in fabric
			_, err := clients.GetPeerFabric(nodestru.PeerNodeName, "access").SetAccess(k_access)
			if err != nil {
				return fmt.Errorf("keyupload set access error:%v", err)
			}

			//send position to center

			res, err := json.Marshal(k_pos)
			if err != nil {
				return fmt.Errorf("marshal error:%v", err)
			}
			topic := "UploadKeyPosition"
			err = ProducerAsyncSending(res, topic, nodestru.CenterAddr)
			if err != nil {
				return fmt.Errorf("producer async sending err:%v", err)
			}
		} else {
			res, err := json.Marshal(key_info)
			if err != nil {
				return fmt.Errorf("marshal error:%v", err)
			}
			topic := "ReceiveKeyUpload" //操作名
			err = ProducerAsyncSending(res, topic, kafkaip)
			if err != nil {
				return fmt.Errorf("producer async sending err:%v", err)
			}
		}
	}
	return nil
}

func receivekeyUpload(nodestru Nodestructure, msg []byte) (err error) {
	fmt.Println("receive key upload")
	var k_detail KeyDetailInfo
	json.Unmarshal(msg, &k_detail)
	//check if it already existed
	if CheckNotExistence(k_detail.FileId, nodestru, "cipherkey_info") {
		var user_information sdkInit.UserInfo
		//get user pubkey and verify
		user_information, err = clients.GetPeerFabric(nodestru.PeerNodeName, "user").GetUserInfo(k_detail.UserId, nodestru.PeerNodeName)
		if err != nil {
			return fmt.Errorf("keyupload get user info error:%v", err)
		}

		key := sdkInit.DecryptionByPri(nodestru.KeyPath, user_information.PublicKey.ConversionMyPub2EcdsaPub(), []byte(k_detail.Key), k_detail.Signlen)
		k_detail.Key = key
		UploadCipherKey(k_detail, nodestru)
	}
	return nil
}

func filerequest(nodestru Nodestructure, msg []byte) (err error) {
	var filereqstru FileRequest
	err = json.Unmarshal(msg, &filereqstru)
	if err != nil {
		return fmt.Errorf("unmarshal error:%v", err)
	}

	//填充，用于之后传回到这个节点
	if filereqstru.Kafka_addr == "" {
		filereqstru.Kafka_addr = nodestru.KafkaIp
	}

	//ciphertext request
	res, err := json.Marshal(filereqstru)
	if err != nil {
		return fmt.Errorf("marshal error:%v", err)
	}
	resultSet, err := Getinfo(filereqstru.FileId, nodestru, "ciphertext_info")
	if err != nil {
		//file在本地不存在，转发到中心节点
		if kivik.StatusCode(err) == 404 {
			//QUERY CENTER
			topic := "FileReqestToCenter"
			err = ProducerAsyncSending(res, topic, nodestru.CenterAddr)
			if err != nil {
				return fmt.Errorf("producer async sending err:%v", err)
			}
		} else {
			return fmt.Errorf("getinfo error :%v", err)
		}
	}

	//if file exist send to client
	if resultSet != nil {
		var fileif FileRequestDTO
		err = resultSet.ScanDoc(&fileif)
		if err != nil {
			return fmt.Errorf("scan doc error:%v", err)
		}
		data := nodestru.PeerNodeName + " : " + fileif.Ciphertext
		//prove that the ciphertext is in this node
		SendData(filereqstru.UserId, data)
	}

	//key request
	//key request sended to center to forward
	topic := "KeyReqForwarding" //操作名
	err = ProducerAsyncSending(res, topic, nodestru.CenterAddr)
	if err != nil {
		fmt.Println("producer async sending err:", err)
	}
	return nil
}

// 边缘节点接收到中心节点转发的key请求
func receivekeyReq(nodestru Nodestructure, msg []byte) (err error) {
	fmt.Println(nodestru.KafkaIp, " --- recieve key req")
	var requestinfo FileRequest
	err = json.Unmarshal(msg, &requestinfo)
	if err != nil {
		return fmt.Errorf("unmarshal error:%v", err)
	}
	//check user access
	var user_information sdkInit.UserInfo

	user_information, err = clients.GetPeerFabric(nodestru.PeerNodeName, "user").GetUserInfo(requestinfo.UserId, nodestru.PeerNodeName)
	if err != nil {
		fmt.Println(err)
	}
	var file_access sdkInit.FileAccessInfo
	file_access, err = clients.GetPeerFabric(nodestru.PeerNodeName, "access").GetAccess(requestinfo.FileId, nodestru.PeerNodeName)
	if err != nil {
		fmt.Println(err)
	}
	if IsIdentical(user_information.Attribute, file_access.Attribute) {
		//take out the key if imformation is match
		rs, err := Getinfo(requestinfo.FileId, nodestru, "cipherkey_info")
		if err != nil {
			if kivik.StatusCode(err) == 404 {
				return fmt.Errorf(requestinfo.FileId, " in ", nodestru.Couchdb_addr, " not founded : ", err)
			} else {
				return fmt.Errorf("getinfo error :%v", err)
			}
		}
		var keydetail KeyDetailInfoDTO
		err = rs.ScanDoc(&keydetail)
		if err != nil {
			return fmt.Errorf("scan doc error:%v", err)
		}
		signed_encrypted_key, sign_len := sdkInit.NodeEncryptionByPubECC(nodestru.KeyPath, user_information.PublicKey.ConversionMyPub2EcdsaPub(), []byte(keydetail.Key))
		intBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(intBytes, uint32(sign_len)) // 转换为 4 字节的小端序字节切片
		//senddata if in the same area
		if requestinfo.Kafka_addr == nodestru.KafkaIp {
			SendData(requestinfo.UserId, append(signed_encrypted_key, intBytes...))
		} else {
			//send to if not in the same area
			keysendinginfo := DataSend2clientInfo{
				UserId: requestinfo.UserId,
				Data:   append(signed_encrypted_key, intBytes...),
			}
			res, err := json.Marshal(keysendinginfo)
			if err != nil {
				return fmt.Errorf("marshal error:%v", err)
			}
			topic := "DataForwarding" //操作名
			err = ProducerAsyncSending(res, topic, requestinfo.Kafka_addr)
			if err != nil {
				return fmt.Errorf("producer async sending err:%v", err)
			}
		}
	}
	return nil
}

// 边缘节点收到了来自中心节点的文件请求
func receiveFileRequestFormCenter(nodestru Nodestructure, msg []byte) (err error) {
	var filereqstru FileRequest
	err = json.Unmarshal(msg, &filereqstru)
	if err != nil {
		return fmt.Errorf("unmarshal error:%v", err)
	}
	resultSet, err := Getinfo(filereqstru.FileId, nodestru, "ciphertext_info")
	if err != nil {
		if kivik.StatusCode(err) == 404 {
			return fmt.Errorf(filereqstru.FileId, " in ", nodestru.Couchdb_addr, " not founded : ", err)
		} else {
			return fmt.Errorf("getinfo error :%v", err)
		}
	}
	var fileif FileRequestDTO
	err = resultSet.ScanDoc(&fileif)
	if err != nil {
		return fmt.Errorf("scan doc error:%v", err)
	}

	//forward to the node linking to client
	datasendinginfo := DataSend2clientInfo{
		UserId:       filereqstru.UserId,
		Data:         []byte(fileif.Ciphertext),
		FileId:       filereqstru.FileId,
		TransferFlag: false,
	}
	// if storageStrategy(nodestru, filereqstru) {
	// 	datasendinginfo.TransferFlag = true
	// }
	res, err := json.Marshal(datasendinginfo)
	if err != nil {
		return fmt.Errorf("marshal error:%v", err)
	}
	topic := "DataForwarding"
	err = ProducerAsyncSending(res, topic, filereqstru.Kafka_addr)
	if err != nil {
		return fmt.Errorf("producer async sending err:%v", err)
	}

	if datasendinginfo.TransferFlag {
		//change position in center
		fileposInfo := PositionInfo{
			FileId:   fileif.FileId,
			Position: filereqstru.Kafka_addr,
		}
		res, err := json.Marshal(fileposInfo)
		if err != nil {
			return fmt.Errorf("marshal error:%v", err)
		}
		//get center addr
		err = ProducerAsyncSending(res, "uploadposition", nodestru.CenterAddr)
		if err != nil {
			return fmt.Errorf("producer async sending err:%v", err)
		}

	}
	return nil
}

// 边缘节点收到了来自中心节点的key请求
func dataForwarding(nodestru Nodestructure, msg []byte) (err error) {
	var dataSendingInfo DataSend2clientInfo
	err = json.Unmarshal(msg, &dataSendingInfo)
	if err != nil {
		return fmt.Errorf("unmarshal error:%v", err)
	}
	//如果设定了转移标志，就把数据存储到本地
	if dataSendingInfo.TransferFlag {
		//store data
		fileif := FileInfo{
			FileId:     dataSendingInfo.FileId,
			Ciphertext: string(dataSendingInfo.Data),
		}
		err = UploadOrUpdateCipherText(fileif, nodestru)
		if err != nil {
			return fmt.Errorf("upload or update ciphertext error:%v", err)
		}
	}
	//send data to client
	SendData(dataSendingInfo.UserId, dataSendingInfo.Data)
	return nil
}

// 到时之后发送数据到中心节点
func time2transmit(fileif FileInfo, nodestru Nodestructure) {
	//transmit the cipherdata to center after 1 min(test)
	// 删除之前的计时器
	timer, ok := timers.Load(fileif.FileId)
	if ok {
		thetime := timer.(*time.Timer)
		thetime.Stop()
	}
	// 创建新的计时器
	timers.Store(fileif.FileId, time.AfterFunc(transmitTime, func() {
		// 到时间后执行的代码
		//begin to send to center
		CipherPosTransfer(fileif.FileId, nodestru, nodestru.CenterAddr)
		// 删除计时器
		timers.Delete(fileif.FileId)
	}))
}

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

// func storageStrategy(nodestru Nodestructure, fileReqStru FileRequest) bool {
// 	//判断是否转移到别的节点
// 	return fileReqStru.AreaId != nodestru.AreaId
// }
