package NodeUtils

import (
	"encoding/binary"
	"encoding/json"
	"fabric-edgenode/clients"
	"fabric-edgenode/models"
	"fabric-edgenode/sdkInit"
	"fmt"

	"github.com/cloudflare/cfssl/log"

	kivik "github.com/go-kivik/kivik/v4"
)

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
	//check if the file is in this node
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	resultSet, err := client.Getinfo(filereqstru.FileId, "ciphertext_info")
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
		log.Error("producer async sending err:", err)
	}
	return nil
}

// 边缘节点接收到中心节点转发的key请求
func receivekeyReq(nodestru Nodestructure, msg []byte) (err error) {
	log.Info("receive key request from center")
	var requestinfo FileRequest
	err = json.Unmarshal(msg, &requestinfo)
	if err != nil {
		return fmt.Errorf("unmarshal error:%v", err)
	}

	//check user access
	var user_information models.UserInfo
	user_information, err = clients.GetPeerFabric(nodestru.PeerNodeName, "user").GetUserInfo(requestinfo.UserId, nodestru.PeerNodeName)
	if err != nil {
		log.Error(err)
	}
	var file_access models.FileAccessInfo
	file_access, err = clients.GetPeerFabric(nodestru.PeerNodeName, "access").GetAccess(requestinfo.FileId, nodestru.PeerNodeName)
	if err != nil {
		log.Error(err)
	}
	if IsIdentical(user_information.Attribute, file_access.Attribute) {
		//take out the key if imformation is match
		client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
		if err != nil {
			return fmt.Errorf("get couchdb client error: %v", err)
		}
		rs, err := client.Getinfo(requestinfo.FileId, "cipherkey_info")
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
	client, err := clients.GetCouchdb(nodestru.Couchdb_addr)
	if err != nil {
		return fmt.Errorf("get couchdb client error: %v", err)
	}
	resultSet, err := client.Getinfo(filereqstru.FileId, "ciphertext_info")
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
			FileId:       fileif.FileId,
			FilePosition: filereqstru.Kafka_addr,
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

// 边缘节点收到了各个节点返回的信息
func dataForwarding(nodestru Nodestructure, msg []byte) (err error) {
	var dataSendingInfo DataSend2clientInfo
	err = json.Unmarshal(msg, &dataSendingInfo)
	if err != nil {
		return fmt.Errorf("unmarshal error:%v", err)
	}
	//send data to client
	SendData(dataSendingInfo.UserId, dataSendingInfo.Data)
	return nil
}
