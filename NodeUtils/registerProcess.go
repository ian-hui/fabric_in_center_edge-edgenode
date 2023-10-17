package NodeUtils

import (
	"encoding/json"
	"fabric-edgenode/clients"
	"fabric-edgenode/models"
	"fabric-edgenode/sdkInit"
	"fmt"

	"github.com/cloudflare/cfssl/log"
)

func register(nodestru Nodestructure, msg []byte) (err error) {
	var userinfo models.UserInfo
	err = json.Unmarshal(msg, &userinfo)
	if err != nil {
		return fmt.Errorf("register unmarshal error:%v", err)
	}
	// prikeyPath := sdkInit.GenerateKey(userif.UserId)
	pubkey := &sdkInit.GetClientPrivateKey("/kafka_crypto/ianhui.private.pem").PublicKey
	userinfo.PublicKey = models.ConversionEcdsaPub2MyPub(pubkey)
	ret, err := clients.GetPeerFabric(nodestru.PeerNodeName, "user").SetUserInfo(userinfo)
	if err != nil {
		return fmt.Errorf("register set error:%v", err)
	}
	log.Info(userinfo.Username, " set success,the transactionID is ", ret, " and userid is ", userinfo.UserId)
	return nil
}
