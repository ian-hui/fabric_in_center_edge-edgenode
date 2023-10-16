package clients

import (
	"encoding/json"
	"fabric-edgenode/models"
	"fabric-edgenode/sdkInit"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	sdk     *fabsdk.FabricSDK
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

func init() {
	temp_sdk, err := fabsdk.New(config.FromFile("/conf/config.yaml"))
	if err != nil {
		log.Printf("Failed to create new SDK: %s\n", err)
		panic(fmt.Sprintf(">>channel %s SDK setup error: %v", "/conf/config.yaml", err))
	}
	sdk = temp_sdk
}

type Channel_client struct {
	ChaincodeID string
	ChClient    *channel.Client
}
type App struct {
	channelInfo *sdkInit.SdkEnvInfo
	cclient     Channel_client
}

func Constructor(PeerNodeAddr string, orgID string, path string, channel_info *sdkInit.SdkEnvInfo, app *App) error {
	orgnum, err := strconv.Atoi(orgID)
	if err != nil {
		return err
	}
	// sdk, err := sdkInit.Setup(path, channel_info)
	// if err != nil {
	// 	panic(fmt.Sprintf(">>channel %s SDK setup error: %v", path, err))
	// }
	// if err := channel_info.InitService(channel_info.ChaincodeID, channel_info.ChannelID, channel_info.Orgs[orgnum-1], sdk); err != nil {
	// 	panic(fmt.Sprintf(">>channel %s InitService unsuccessful: %v", path, err))
	// }

	clientChannelContext1 := sdk.ChannelContext(channel_info.ChannelID, fabsdk.WithUser(channel_info.OrdererAdminUser), fabsdk.WithOrg(channel_info.Orgs[orgnum-1].OrgName))
	client, err := channel.New(clientChannelContext1)
	if err != nil {
		log.Printf("Failed to create new channel client: %s\n", err)
		return err
	}
	*app = App{
		channelInfo: channel_info,
		cclient:     Channel_client{ChaincodeID: channel_info.ChaincodeID, ChClient: client},
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

func GetPeerFabric(PeerNodeName string, app_type string) Channel_client {
	switch app_type {
	case "user":
		return userApp.cclient
	case "access":
		return accessApp.cclient
	case "node":
		return nodeApp.cclient
	default:
		panic(fmt.Sprintln(">>channel", PeerNodeName, "GetPeerFabric error: unknown app type"))
	}
}

func FabricClose() {
	//close sdk
	// userApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.BlockListener(userApp.app.SdkEnvInfo.EvClient))
	// userApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.ChainCodeEventListener(userApp.app.SdkEnvInfo.EvClient, userApp.app.SdkEnvInfo.ChaincodeID))
	// accessApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.BlockListener(accessApp.app.SdkEnvInfo.EvClient))
	// accessApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.ChainCodeEventListener(accessApp.app.SdkEnvInfo.EvClient, accessApp.app.SdkEnvInfo.ChaincodeID))
	// nodeApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.BlockListener(nodeApp.app.SdkEnvInfo.EvClient))
	// nodeApp.app.SdkEnvInfo.EvClient.Unregister(sdkInit.ChainCodeEventListener(nodeApp.app.SdkEnvInfo.EvClient, nodeApp.app.SdkEnvInfo.ChaincodeID))
	// userApp.sdk.Close()
	// accessApp.sdk.Close()
	// nodeApp.sdk.Close()
	sdk.Close()
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
