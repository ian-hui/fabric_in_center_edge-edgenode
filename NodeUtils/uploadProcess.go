package NodeUtils

import (
	"encoding/json"
	"fabric-edgenode/clients"
	"fabric-edgenode/sdkInit"
	"fmt"
	"time"

	"github.com/fatih/structs"
)

func upload(nodestru Nodestructure, msg []byte) (err error) {
	fmt.Println("begin to upload locally")

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
	fmt.Println("<----begin to upload position to center---->")
	//实例化一个groupchooser
	gc := &groupChooser{
		storage_weight:            0.5,
		distance_weight:           0.5,
		standard_deviation_weight: 0.5,
		curNodeInfo:               nodestru.NodeInfo,
	}

	//获取分布式锁：TODO
	//计算出group
	group, _, _, err := gc.chooseGroup(5)
	if err != nil {
		return fmt.Errorf("choose group error:%v", err)
	}

	//修改nodeinfo：TODO
	//放回分布式锁：TODO

	//密文位置，密钥位置
	file_postion_info := PositionInfo{
		FileId:     fileif.FileId,
		Position:   nodestru.KafkaIp,
		GroupAddrs: group,
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

// func chooseGroup(nodestru Nodestructure, msg []byte) (err error) {
// 	//实例化一个groupchooser
// 	gc := &groupChooser{
// 		storage_weight:            0.5,
// 		distance_weight:           0.5,
// 		standard_deviation_weight: 0.5,
// 		curNodeInfo:               nodestru.NodeInfo,
// 	}

// 	group, delay, sd, err := gc.chooseGroup(6)
// 	if err != nil {
// 		return fmt.Errorf("choose group error:%v", err)
// 	}
// 	marshaled_group, err := json.Marshal(group)
// 	if err != nil {
// 		return fmt.Errorf("marshal error:%v", err)
// 	}
// 	SendData(string(msg), marshaled_group)
// 	return
// }

// 上传密钥
func keyUpload(nodestru Nodestructure, msg []byte) (err error) {
	var keyinfostru KeyUploadInfo
	err = json.Unmarshal(msg, &keyinfostru)
	if err != nil {
		return fmt.Errorf("keyupload unmarshal error:%v", err)
	}
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
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
			check := client.CheckNotExistence(key_info.FileId, "cipherkey_info")
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
			err = client.CouchdbPut(key_info.FileId, structs.Map(&key_info), "cipherkey_info")
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

//各个节点获取到各自的密钥片段
func receivekeyUpload(nodestru Nodestructure, msg []byte) (err error) {
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
		client.CouchdbPut(k_detail.FileId, structs.Map(&k_detail), "cipherkey_info")
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
