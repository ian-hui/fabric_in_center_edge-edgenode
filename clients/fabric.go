package clients

import (
	"encoding/json"
	"fabric-edgenode/models"
	"fabric-edgenode/sdkInit"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	userApp = &App{
		channelInfo: &UserinfoChannel_info,
	}
	accessApp = &App{
		channelInfo: &AccessChannel_info,
	}
	nodeApp = &App{
		channelInfo: &NodeinfoChannel_info,
	}
)

type App struct {
	channelInfo *sdkInit.SdkEnvInfo
	app         *sdkInit.Application
	sdk         *fabsdk.FabricSDK
}

func Constructor(PeerNodeAddr string, orgID string, path string, channel_info *sdkInit.SdkEnvInfo, app *App) error {
	orgnum, err := strconv.Atoi(orgID)
	if err != nil {
		return err
	}
	sdk, err := sdkInit.Setup(path, channel_info)
	if err != nil {
		panic(fmt.Sprintf(">>channel %s SDK setup error: %v", path, err))
	}
	if err := channel_info.InitService(channel_info.ChaincodeID, channel_info.ChannelID, channel_info.Orgs[orgnum-1], sdk); err != nil {
		panic(fmt.Sprintf(">>channel %s InitService unsuccessful: %v", path, err))
	}
	*app = App{
		channelInfo: channel_info,
		app:         &sdkInit.Application{SdkEnvInfo: channel_info},
		sdk:         sdk,
	}
	return nil
}

func InitPeerSdk(PeerNodeAddr string, orgID string, path string) error {
	apps := []*App{
		userApp,
		accessApp,
		nodeApp,
	}
	for i := range apps {
		fmt.Println(">>channel", PeerNodeAddr, "InitPeerSdk begin")
		if err := Constructor(PeerNodeAddr, orgID, path, apps[i].channelInfo, apps[i]); err != nil {
			return err
		}
		fmt.Println(">>channel", PeerNodeAddr, "InitPeerSdk successful")
	}
	return nil
}

func GetPeerFabric(PeerNodeName string, app_type string) *sdkInit.Application {
	switch app_type {
	case "user":
		a := userApp.app
		return a
	case "access":
		return accessApp.app
	case "node":
		return nodeApp.app
	default:
		panic(fmt.Sprintln(">>channel", PeerNodeName, "GetPeerFabric error: unknown app type"))
	}
}

func FabricClose() {
	//close sdk
	userApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.BlockListener(userApp.app.SdkEnvInfo.EvClient))
	userApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.ChainCodeEventListener(userApp.app.SdkEnvInfo.EvClient, userApp.app.SdkEnvInfo.ChaincodeID))
	accessApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.BlockListener(accessApp.app.SdkEnvInfo.EvClient))
	accessApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.ChainCodeEventListener(accessApp.app.SdkEnvInfo.EvClient, accessApp.app.SdkEnvInfo.ChaincodeID))
	nodeApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.BlockListener(nodeApp.app.SdkEnvInfo.EvClient))
	nodeApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.ChainCodeEventListener(nodeApp.app.SdkEnvInfo.EvClient, nodeApp.app.SdkEnvInfo.ChaincodeID))
	userApp.sdk.Close()
	accessApp.sdk.Close()
	nodeApp.sdk.Close()
}

func InitNodeInfo(nodeinfo models.NodeInfo) error {
	//初始化节点信息
	nodeInfo_byte, err := json.Marshal(nodeinfo)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %v", err)
	}
	s, err := GetPeerFabric(nodeinfo.PeerNodeName, "node").SetNodeInfo(nodeinfo.KafkaAddr, nodeInfo_byte)
	if err != nil {
		return fmt.Errorf("SetNodeInfo error: %v", err)
	}
	log.Println(">>channel", nodeinfo.PeerNodeName, "InitNodeInfo successful", s)
	return nil

}
