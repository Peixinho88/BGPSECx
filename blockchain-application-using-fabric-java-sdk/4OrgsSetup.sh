#!/bin/bash

set -e

./cleanSetup.sh
./4OrgsTest/createChannel.sh
./4OrgsTest/updateAnchor.sh
./4OrgsTest/packageInstallCC.sh
./instantiateChaincode.sh
