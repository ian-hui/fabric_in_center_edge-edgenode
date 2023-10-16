package NodeUtils

import (
	"encoding/json"
	"fabric-edgenode/clients"
	"fabric-edgenode/sdkInit"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
)

func upload(nodestru Nodestructure, msg []byte) (err error) {
	log.Println("begin to upload locally")

	//先上传密文到本地
	var fileif FileInfo
	err = json.Unmarshal(msg, &fileif)
	if err != nil {
		return fmt.Errorf("upload unmarshal error:%v", err)
	}
	//check the file not exist
	couchdbClient, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	//upload ciphertext in local couchdb
	err = couchdbClient.CouchdbPut(fileif.FileId, structs.Map(&fileif), "ciphertext_info")
	if err != nil {
		return fmt.Errorf("upload ciphertext error:%v", err)
	}

	//然后把信息上传到center
	log.Println("<----begin to upload position to center---->")

	//密文位置，密钥位置
	file_postion_info := PositionInfo{
		FileId:   fileif.FileId,
		Position: nodestru.KafkaIp,
		// GroupAddrs: group,
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

// 上传密钥
func keyUpload(nodestru Nodestructure, msg []byte) (err error) {
	//获取这次的数据
	var keyinfostru KeyUploadInfo
	err = json.Unmarshal(msg, &keyinfostru)
	if err != nil {
		return fmt.Errorf("keyupload unmarshal error:%v", err)
	}
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	//实例化一个groupchooser
	group, _, _, err := NewGroupChooser(nodestru.NodeInfo, 0.5, 0.5, 0.5).ChooseGroupWithLock(nodestru)
	if err != nil {
		return fmt.Errorf("keyupload choose group error:%v", err)
	}
	log.Println("group:", group)

	check_map := make(map[string]bool)
	for i := range group {
		check_map[group[i]] = true
	}

	//set access in fabric
	_, err = clients.GetPeerFabric(nodestru.PeerNodeName, "access").SetAccess(sdkInit.FileAccessInfo{
		FileId:    (*keyinfostru.Upload_Infomation)[nodestru.KafkaIp].FileId,
		Attribute: keyinfostru.Attribute,
	})
	if err != nil {
		return fmt.Errorf("keyupload set access error:%v", err)
	}
	//send position to center
	res, err := json.Marshal(KeyPostionUploadInfo{
		FileId:     (*keyinfostru.Upload_Infomation)[nodestru.KafkaIp].FileId,
		GroupAddrs: group,
	})
	if err != nil {
		return fmt.Errorf("marshal error:%v", err)
	}
	err = ProducerAsyncSending(res, "UploadKeyPosition", nodestru.CenterAddr)
	if err != nil {
		return fmt.Errorf("producer async sending err:%v", err)
	}

	for KafkaAddr, key_info := range *keyinfostru.Upload_Infomation {
		var (
			err              error
			user_information sdkInit.UserInfo
		)
		if KafkaAddr == nodestru.KafkaIp && check_map[nodestru.KafkaIp] {
			//check the key not exist
			log.Println("uploading key to local couchdb...")
			check := client.CheckNotExistence(key_info.FileId, "cipherkey_info")
			if !check {
				return fmt.Errorf("keyupload error: key already exist")
			}
			user_information, err = clients.GetPeerFabric(nodestru.PeerNodeName, "user").GetUserInfo(key_info.UserId, nodestru.PeerNodeName)
			if err != nil {
				return fmt.Errorf("keyupload get user info error:%v", err)
			}
			key_info.Key = sdkInit.DecryptionByPri(nodestru.KeyPath, user_information.PublicKey.ConversionMyPub2EcdsaPub(), []byte(key_info.Key), key_info.Signlen)
			//upload cipherkey in local peernode
			err = client.CouchdbPut(key_info.FileId, structs.Map(&key_info), "cipherkey_info")
			if err != nil {
				return fmt.Errorf("keyupload upload cipherkey error:%v", err)
			}
		} else if check_map[KafkaAddr] {
			log.Println("sending kafkaAddr to it node:", KafkaAddr, "...")
			res, err := json.Marshal(key_info)
			if err != nil {
				return fmt.Errorf("marshal error:%v", err)
			}
			topic := "ReceiveKeyUpload" //操作名
			err = ProducerAsyncSending(res, topic, KafkaAddr)
			if err != nil {
				return fmt.Errorf("producer async sending err:%v", err)
			}
		} else {
			continue
		}
	}
	return nil
}

//各个节点获取到各自的密钥片段
func receivekeyUpload(nodestru Nodestructure, msg []byte) (err error) {
	log.Println("receive keyupload request! begin to upload locally")
	var k_detail KeyDetailInfo
	json.Unmarshal(msg, &k_detail)
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	//check if it already existed
	if client.CheckNotExistence(k_detail.FileId, "cipherkey_info") {
		var user_information sdkInit.UserInfo
		//get user pubkey and verify
		user_information, err = clients.GetPeerFabric(nodestru.PeerNodeName, "user").GetUserInfo(k_detail.UserId, nodestru.PeerNodeName)
		if err != nil {
			return fmt.Errorf("keyupload get user info error:%v", err)
		}

		key := sdkInit.DecryptionByPri(nodestru.KeyPath, user_information.PublicKey.ConversionMyPub2EcdsaPub(), []byte(k_detail.Key), k_detail.Signlen)
		k_detail.Key = key
		if err = client.CouchdbPut(k_detail.FileId, structs.Map(&k_detail), "cipherkey_info"); err != nil {
			return fmt.Errorf("keyupload upload cipherkey error:%v", err)
		}

	}
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
