package clients

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t Channel_client) GetNodeInfo(nodeIp string, endpoint string) ([]byte, error) {
	response, err := t.ChClient.Query(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(nodeIp)}},
		channel.WithTargetEndpoints(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return response.Payload, nil
}

func (t Channel_client) SetNodeInfo(PeerNodeName string, info []byte) (string, error) {
	request := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(PeerNodeName), info}}
	response, err := t.ChClient.Execute(request)
	if err != nil {
		// set失败
		return "", err
	}
	//fmt.Println("============== response:",response)
	return string(response.TransactionID), nil
}

func (t Channel_client) GetNodeInfoAllRange(endpoint string) ([]byte, error) {
	response, err := t.ChClient.Query(channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "rangeQuery", Args: [][]byte{[]byte("k"), []byte("l")}},
		channel.WithTargetEndpoints(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return response.Payload, nil
}
