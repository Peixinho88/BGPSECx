#!/bin/bash
#
# Exit on first error, print all commands.
set -ev

# Shut down the Docker containers for the system tests.
docker-compose -f docker-compose1.yml kill && docker-compose -f docker-compose1.yml down
if [ "$(docker ps -aq)" ]; then
	docker rm -f $(docker ps -aq)
fi

# remove chaincode docker images
if [ "$(docker images dev-* -q)" ]; then
	docker rmi $(docker images dev-* -q)
fi

docker system prune --volumes -f

# Your system is now clean
