#!/bin/bash
#
#Exit on first error, print all commands.
set -e

#Put the cryptogen and configtxgen tools in the PATH environment variable
export PATH=../fabric-samples/bin/:$PATH

#Generate the artifacts
cryptogen generate --config=./crypto-config.yaml
mkdir config
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./config/genesis.block
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/channel.tx -channelID mychannel
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org3MSPanchors.tx -channelID mychannel -asOrg Org3MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org4MSPanchors.tx -channelID mychannel -asOrg Org4MSP
