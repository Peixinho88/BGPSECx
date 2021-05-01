#!/bin/bash
#
#Exit on first error, print all commands.
set -e

cd network
./generateCrypto.sh

rm -r ../network_resources/config
rm -r ../network_resources/crypto-config
mv -f config ../network_resources/config
mv -f crypto-config ../network_resources/crypto-config

CURRENT_DIR=$PWD
cd ../network_resources/crypto-config/peerOrganizations/org1.example.com/ca/
PRIV_KEY=$(ls *_sk)
cd "$CURRENT_DIR"
sed -Ei "s/\w+_sk/${PRIV_KEY}/g" docker-compose.yml
