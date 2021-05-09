#!/bin/bash
#
# Exit on first error, print all commands.
set -e

#Start from here
echo -e "\nStopping the previous network (if any)"
docker-compose -f docker-compose-main.yml down

# Create and Start the Docker containers for the network
echo -e "\nSetting up the Hyperledger Fabric 1.4.1 network"
docker-compose -f docker-compose-main.yml up -d
sleep 15
echo -e "\nNetwork setup completed!!\n"

