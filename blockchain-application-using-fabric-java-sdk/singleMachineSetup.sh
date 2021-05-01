#!/bin/bash

set -e

./cleanSetup.sh
./createChannel.sh
./updateAnchor.sh
./packageInstallCC.sh
./instantiateChaincode.sh
