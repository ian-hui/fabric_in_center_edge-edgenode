package clients

import (
	"fabric-edgenode/sdkInit"
)

var (
	Orgs_userinfoChannel = []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org1MSPanchors_ChannelOne.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org2MSPanchors_ChannelOne.tx",
		}, {
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org3MSPanchors_ChannelOneAreaTwo.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org4MSPanchors_ChannelOneAreaTwo.tx",
		},
	}

	Orgs_accessChannel = []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org1MSPanchors_ChannelOne.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org2MSPanchors_ChannelOne.tx",
		}, {
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org3MSPanchors_ChannelOneAreaTwo.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org4MSPanchors_ChannelOneAreaTwo.tx",
		},
	}

	Orgs_nodeinfoChannel = []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org1MSPanchors_ChannelOne.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org2MSPanchors_ChannelOne.tx",
		}, {
			OrgAdminUser:  "Admin",
			OrgName:       "Org3",
			OrgMspId:      "Org3MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org3MSPanchors_ChannelOneAreaTwo.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org4",
			OrgMspId:      "Org4MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "./fixtures/channel-artifacts/Org4MSPanchors_ChannelOneAreaTwo.tx",
		},
	}

	//init sdk env info
	userinfoChannel_info = sdkInit.SdkEnvInfo{
		ChannelID:        "myUserinfoChannel",
		ChannelConfig:    "./fixtures/channel-artifacts/channel1.tx",
		Orgs:             Orgs_userinfoChannel,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "user_information",
		ChaincodePath:    "./chaincode/",
		ChaincodeVersion: "1.0.0",
	}
	accessChannel_info = sdkInit.SdkEnvInfo{
		ChannelID:        "myAccessChannel",
		ChannelConfig:    "./fixtures/channel-artifacts/channel2.tx",
		Orgs:             Orgs_accessChannel,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "access_information",
		ChaincodePath:    "./chaincode/",
		ChaincodeVersion: "1.0.0",
	}
	nodeinfoChannel_info = sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel3",
		ChannelConfig:    "./fixtures/channel-artifacts/channel3.tx",
		Orgs:             Orgs_nodeinfoChannel,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      "user_information_2",
		ChaincodePath:    "./chaincode/",
		ChaincodeVersion: "1.0.0",
	}
)
