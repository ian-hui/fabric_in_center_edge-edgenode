rm -rf crypto-config channel-artifacts && mkdir crypto-config channel-artifacts
../../fabric-samples/bin/cryptogen generate --config=crypto-config.yaml
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID fabric-channel
#generate channel config-channel.tx
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelOne -outputCreateChannelTx ./channel-artifacts/channel1.tx -channelID mychannel1
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelTwo -outputCreateChannelTx ./channel-artifacts/channel2.tx -channelID mychannel2
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelOneAreaTwo -outputCreateChannelTx ./channel-artifacts/channel3.tx -channelID mychannel3
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelTwoAreaTwo -outputCreateChannelTx ./channel-artifacts/channel4.tx -channelID mychannel4
#generate anchorpeerfile for every channel
#channel1
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelOne -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors_ChannelOne.tx -channelID mychannel1 -asOrg Org1MSP
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelOne -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors_ChannelOne.tx -channelID mychannel1 -asOrg Org2MSP
#channel2
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelTwo -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors_ChannelTwo.tx -channelID mychannel2 -asOrg Org1MSP
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelTwo -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors_ChannelTwo.tx -channelID mychannel2 -asOrg Org2MSP

../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelOneAreaTwo -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPanchors_ChannelOneAreaTwo.tx -channelID mychannel3 -asOrg Org3MSP
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelOneAreaTwo -outputAnchorPeersUpdate ./channel-artifacts/Org4MSPanchors_ChannelOneAreaTwo.tx -channelID mychannel3 -asOrg Org4MSP

../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelTwoAreaTwo -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPanchors_ChannelTwoAreaTwo.tx -channelID mychannel4 -asOrg Org3MSP
../../github.com/hyperledger/fabric/build/bin/configtxgen -profile ChannelTwoAreaTwo -outputAnchorPeersUpdate ./channel-artifacts/Org4MSPanchors_ChannelTwoAreaTwo.tx -channelID mychannel4 -asOrg Org4MSP