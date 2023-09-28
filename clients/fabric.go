package clients

import (
	"fabric-edgenode/sdkInit"
	"fmt"
	"strconv"
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
		fmt.Println(userApp)
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
