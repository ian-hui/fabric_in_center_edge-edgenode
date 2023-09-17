package sdkInit

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *Application) GetUserInfo(userid string, endpoint string) (UserInfo, error) {
	var userif UserInfo
	response, err := t.SdkEnvInfo.ChClient.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(userid)}},
		channel.WithTargetEndpoints(endpoint))
	if err != nil {
		return userif, fmt.Errorf("failed to query: %v", err)
	}
	// 对查询到的状态进行反序列化
	err = json.Unmarshal(response.Payload, &userif)
	if err != nil {
		return userif, err
	}
	return userif, nil
}

func (t *Application) SetUserInfo(usrinfo UserInfo) (string, error) {
	b, err := json.Marshal(usrinfo)
	if err != nil {
		return "", fmt.Errorf("指定的userinfo对象序列化时发生错误")
	}
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(usrinfo.UserId), b}}
	response, err := t.SdkEnvInfo.ChClient.Execute(request)
	if err != nil {
		// set失败
		return "", err
	}

	//fmt.Println("============== response:",response)

	return string(response.TransactionID), nil
}
