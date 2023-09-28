package clients

import "fabric-edgenode/sdkInit"

var (
	Orgs_userinfoChannel = []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org1MSPanchors_userinfoChannel.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org2MSPanchors_userinfoChannel.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org3MSPanchors_userinfoChannel.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org4MSPanchors_userinfoChannel.tx",
		},
	}

	Orgs_accessChannel = []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org1MSPanchors_accessChannel.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org2MSPanchors_accessChannel.tx",
		}, {
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org3MSPanchors_accessChannel.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org4MSPanchors_accessChannel.tx",
		},
	}

	Orgs_nodeinfoChannel = []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org1MSPanchors_nodeinfoChannel.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org2MSPanchors_nodeinfoChannel.tx",
		}, {
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org3MSPanchors_nodeinfoChannel.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/fixtures/channel-artifacts/Org4MSPanchors_nodeinfoChannel.tx",
		},
	}

	//init sdk env info
	UserinfoChannel_info = sdkInit.SdkEnvInfo{
		ChannelID:        "myuserinfochannel",
		ChannelConfig:    "/fixtures/channel-artifacts/userinfoChannel.tx",
		Orgs:             Orgs_userinfoChannel,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "user_information",
		ChaincodePath:    "./chaincode/",
		ChaincodeVersion: "1.0.0",
	}
	AccessChannel_info = sdkInit.SdkEnvInfo{
		ChannelID:        "myaccesschannel",
		ChannelConfig:    "/fixtures/channel-artifacts/accessChannel.tx",
		Orgs:             Orgs_accessChannel,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "access_information",
		ChaincodePath:    "./chaincode/",
		ChaincodeVersion: "1.0.0",
	}
	NodeinfoChannel_info = sdkInit.SdkEnvInfo{
		ChannelID:        "mynodeinfochannel",
		ChannelConfig:    "/fixtures/channel-artifacts/nodeinfoChannel.tx",
		Orgs:             Orgs_nodeinfoChannel,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "nodeinfo_channel",
		ChaincodePath:    "./chaincode/",
		ChaincodeVersion: "1.0.0",
	}
)
