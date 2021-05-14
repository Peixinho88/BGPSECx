#!/bin/bash

set -e

docker exec -it cli bash -c '
CORE_PEER_LOCALMSPID="Org3MSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
CORE_PEER_ADDRESS=peer0.org3.example.com:7051
peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org3MSPanchors.tx 
exit ; bash'
