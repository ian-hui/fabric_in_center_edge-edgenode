package NodeUtils

import (
	"encoding/base64"
	"fabric-edgenode/sdkInit"
)

type Nodestructure struct {
	KafkaIp          string
	Couchdb_addr     string
	PeerNodeName     string
	OrgID            string
	KeyPath          string
	ConfigPath       string
	NodeInfo         *NodeInfo
	UserChannel_info *sdkInit.SdkEnvInfo
}

type PositionInfo struct {
	// AreaId   string `json:"AreaId"`
	FileId   string   `json:"FileId"`
	Position string   `json:"Position"`
	GroupIps []string `json:"GroupIps"`
}

type FileInfo struct {
	FileId     string `json:"FileId"`
	Ciphertext string `json:"Ciphertext"`
}

type FileRequest struct {
	FileId string `json:"FileId"`
	UserId string `json:"UserId"`
	// AreaId      string `json:"AreaId"`
	Kafka_addr  string `json:"kafka_addr"`
	storageFlag bool   `json:"storageFlag"`
}

type FileRequestDTO struct {
	FileId     string `json:"FileId"`
	Ciphertext string `json:"Ciphertext"`
	_id        string `json:"_id"`
	_rev       string `json:"_rev"`
}

type FilePositionInfoDTO struct {
	FileId   string `json:"FileId"`
	Position string `json:"Position"`
	AreaId   string `json:"AreaId"`
	_rev     string `json:"_rev"`
	_id      string `json:"_id"`
}

type KeyDetailInfoDTO struct {
	FileId string `json:"FileId"`
	Key    string `json:"Key"`
	_rev   string `json:"_rev"`
	_id    string `json:"_id"`
}

type KeyDetailInfo struct {
	FileId  string `json:"FileId"`
	Signlen int    `json:"Signlen"`
	Key     []byte `json:"Key"`
	UserId  string `json:UserId`
}

func (k *KeyDetailInfo) SetKeyFromBase64(encoded string) error {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	k.Key = decoded
	return nil
}

type KeyPostionUploadInfo struct {
	FileId   string   `json:"FileId"`
	GroupIps []string `json:"GroupIps"`
}

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

type NodeInfo struct {
	NodeIp      string
	LeftStorage string
	locationX   string
	locationY   string
}
