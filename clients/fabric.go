package clients

import (
	"fabric-edgenode/sdkInit"
	"fmt"
	"strconv"
	"sync"
)

var (
	userApp   *UserApp
	accessApp *AccessApp
	nodeApp   *NodeApp
)

func init() {
	userApp = &UserApp{
		channelInfo: &userinfoChannel_info,
		appMap:      &sync.Map{},
	}
	accessApp = &AccessApp{
		channelInfo: &accessChannel_info,
		appMap:      &sync.Map{},
	}
	nodeApp = &NodeApp{
		channelInfo: &nodeinfoChannel_info,
		appMap:      &sync.Map{},
	}
}

// 三种类型的app
type App interface {
	Init(PeerNodeAddr string, orgID string, path string) error
}

type UserApp struct {
	channelInfo *sdkInit.SdkEnvInfo
	appMap      *sync.Map
}

func (u *UserApp) Init(PeerNodeAddr string, orgID string, path string) error {
	return initApp(PeerNodeAddr, orgID, path, u.channelInfo, u.appMap)
}

type AccessApp struct {
	channelInfo *sdkInit.SdkEnvInfo
	appMap      *sync.Map
}

func (a *AccessApp) Init(PeerNodeAddr string, orgID string, path string) error {
	return initApp(PeerNodeAddr, orgID, path, a.channelInfo, a.appMap)
}

type NodeApp struct {
	channelInfo *sdkInit.SdkEnvInfo
	appMap      *sync.Map
}

func (n *NodeApp) Init(PeerNodeAddr string, orgID string, path string) error {
	return initApp(PeerNodeAddr, orgID, path, n.channelInfo, n.appMap)
}

func InitPeerSdk(PeerNodeAddr string, orgID string, path string) error {
	apps := []App{
		userApp,
		accessApp,
		nodeApp,
	}
	for _, app := range apps {
		if err := app.Init(PeerNodeAddr, orgID, path); err != nil {
			return err
		}
	}

	return nil
}

func initApp(PeerNodeAddr string, orgID string, path string, channel_info *sdkInit.SdkEnvInfo, app *sync.Map) error {
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
	app.Store(PeerNodeAddr, &sdkInit.Application{
		SdkEnvInfo: channel_info,
	})
	return nil
}

func GetPeerFabric(PeerNodeName string, app_type string) *sdkInit.Application {
	switch app_type {
	case "user":
		if app, ok := userApp.appMap.Load(PeerNodeName); ok {
			return app.(*sdkInit.Application)
		}
	case "access":
		if app, ok := accessApp.appMap.Load(PeerNodeName); ok {
			return app.(*sdkInit.Application)
		}
	case "node":
		if app, ok := nodeApp.appMap.Load(PeerNodeName); ok {
			return app.(*sdkInit.Application)
		}
	default:
		panic(fmt.Sprintln(">>channel", PeerNodeName, "GetPeerFabric error: unknown app type"))
	}
	return nil
}
