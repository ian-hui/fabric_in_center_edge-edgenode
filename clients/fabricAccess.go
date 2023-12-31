package clients

import (
	"encoding/json"
	"fabric-edgenode/models"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t Channel_client) GetAccess(fileid string, endpoint string) (fileaccessinfo models.FileAccessInfo, err error) {
	response, err := t.ChClient.Query(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(fileid)}},
		channel.WithTargetEndpoints(endpoint))
	if err != nil {
		return fileaccessinfo, fmt.Errorf("failed to query: %v", err)
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(response.Payload, &fileaccessinfo)
	if err != nil {
		return fileaccessinfo, err
	}
	return fileaccessinfo, nil
}

func (t Channel_client) SetAccess(fileaccessinfo models.FileAccessInfo) (string, error) {
	b, err := json.Marshal(fileaccessinfo)
	if err != nil {
		return "", fmt.Errorf("指定的fileaccessinfo对象序列化时发生错误")
	}

	request := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(fileaccessinfo.FileId), b}}
	response, err := t.ChClient.Execute(request)
	if err != nil {
		// set失败
		return "", err
	}

	//fmt.Println("============== response:",response)

	return string(response.TransactionID), nil
}
