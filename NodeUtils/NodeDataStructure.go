package NodeUtils

import (
	"encoding/base64"
	"fabric-edgenode/models"
	"fabric-edgenode/sdkInit"
)

type Nodestructure struct {
	KafkaIp          string
	Couchdb_addr     string
	ZookeeperAddr    string
	PeerNodeName     string
	OrgID            string
	KeyPath          string
	CenterAddr       string
	ConfigPath       string
	NodeInfo         *models.NodeInfo
	UserChannel_info *sdkInit.SdkEnvInfo
}

type PositionInfo struct {
	// AreaId   string `json:"AreaId"`
	FileId        string   `json:"FileId"`
	FilePosition  string   `json:"FilePosition"`
	KeyGroupAddrs []string `json:"KeyGroupAddrs"`
}

type FileInfo struct {
	FileId     string `json:"FileId"`
	Ciphertext string `json:"Ciphertext"`
}

type FileRequest struct {
	FileId     string `json:"FileId"`
	UserId     string `json:"UserId"`
	Kafka_addr string `json:"Kafka_addr"`
}

type FileRequestDTO struct {
	FileId     string `json:"FileId"`
	Ciphertext string `json:"Ciphertext"`
	Id         string `json:"_id"`
	Rev        string `json:"_rev"`
}

type FilePositionInfoDTO struct {
	FileId   string `json:"FileId"`
	Position string `json:"Position"`
	AreaId   string `json:"AreaId"`
	Rev      string `json:"_rev"`
	Id       string `json:"_id"`
}

type KeyDetailInfoDTO struct {
	FileId string `json:"FileId"`
	Key    string `json:"Key"`
	Rev    string `json:"_rev"`
	Id     string `json:"_id"`
}

type KeyDetailInfo struct {
	FileId  string `json:"FileId"`
	Signlen int    `json:"Signlen"`
	Key     []byte `json:"Key"`
	UserId  string `json:"UserId"`
}

func (k *KeyDetailInfo) SetKeyFromBase64(encoded string) error {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	k.Key = decoded
	return nil
}

// type KeyPostionUploadInfo struct {
// 	FileId     string   `json:"FileId"`
// 	GroupAddrs []string `json:"GroupAddrs"`
// }

type DataSend2clientInfo struct {
	TransferFlag bool   `json:"TransferFlag"`
	Data         []byte `json:"Data"`
	FileId       string `json:"FileId"`
	UserId       string `json:"UserId"`
}

type KeyUploadInfo struct {
	Upload_Infomation *map[string]KeyDetailInfo
	Attribute         []string
}
