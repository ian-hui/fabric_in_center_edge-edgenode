package NodeUtils

import (
	"encoding/json"
	"fabric-edgenode/clients"
	"fabric-edgenode/sdkInit"
	"fmt"
)

func register(nodestru Nodestructure, msg []byte) (err error) {
	var userif sdkInit.UserInfo
	err = json.Unmarshal(msg, &userif)
	if err != nil {
		return fmt.Errorf("register unmarshal error:%v", err)
	}
	// prikeyPath := sdkInit.GenerateKey(userif.UserId)
	pubkey := &sdkInit.GetClientPrivateKey("/kafka_crypto/ianhui.private.pem").PublicKey
	userif.PublicKey = sdkInit.ConversionEcdsaPub2MyPub(pubkey)
	ret, err := clients.GetPeerFabric(nodestru.PeerNodeName, "user").SetUserInfo(userif)
	if err != nil {
		return fmt.Errorf("register set error:%v", err)
	}
	fmt.Println(userif.Username, " set success,the transactionID is ", ret, " and userid is ", userif.UserId)
	return nil
}
