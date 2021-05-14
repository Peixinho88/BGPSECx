#!/bin/bash

set -e

docker exec -it cli bash -c '
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n fabcar -l golang -v 1.0 -c "{\"Args\":[\"init\"]}"
exit ; bash'
