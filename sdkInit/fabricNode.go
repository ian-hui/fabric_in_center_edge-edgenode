package sdkInit

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *Application) GetNodeInfo(nodeIp string, endpoint string) ([]byte, error) {
	response, err := t.SdkEnvInfo.ChClient.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(nodeIp)}},
		channel.WithTargetEndpoints(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return response.Payload, nil
}

func (t *Application) SetNodeInfo(PeerNodeName string, info []byte) (string, error) {
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(PeerNodeName), info}}
	response, err := t.SdkEnvInfo.ChClient.Execute(request)
	if err != nil {
		// set失败
		return "", err
	}
	//fmt.Println("============== response:",response)
	return string(response.TransactionID), nil
}

func (t *Application) GetNodeInfoAllRange(endpoint string) ([]byte, error) {
	response, err := t.SdkEnvInfo.ChClient.Query(channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: "rangeQuery", Args: [][]byte{[]byte("k"), []byte("l")}},
		channel.WithTargetEndpoints(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return response.Payload, nil
}
