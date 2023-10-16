package clients

import (
	"encoding/json"
	"fabric-edgenode/models"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t Channel_client) GetUserInfo(userid string, endpoint string) (userinfo models.UserInfo, err error) {
	response, err := t.ChClient.Query(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(userid)}},
		channel.WithTargetEndpoints(endpoint))
	if err != nil {
		return userinfo, fmt.Errorf("failed to query: %v", err)
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(response.Payload, &userinfo)
	if err != nil {
		return userinfo, fmt.Errorf("指定的userinfo对象反序列化时发生错误:%v", err)
	}
	return
}

func (t Channel_client) SetUserInfo(userinfo models.UserInfo) (string, error) {
	b, err := json.Marshal(userinfo)
	if err != nil {
		return "", fmt.Errorf("指定的userinfo对象序列化时发生错误")
	}
	request := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(userinfo.UserId), b}}
	response, err := t.ChClient.Execute(request)
	if err != nil {
		// set失败
		return "", err
	}

	//fmt.Println("============== response:",response)

	return string(response.TransactionID), nil
}
