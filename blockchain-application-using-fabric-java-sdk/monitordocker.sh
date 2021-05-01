#!/bin/bash
# You can use this command to get this script
#
# curl -LOs https://gist.github.com/mbwhite/a32abc57a0a45ecc466977ceef67df1f/raw/monitordocker.sh && chmod +x monitordocker.sh
#
# This script uses the logspout and http stream tools to let you watch the docker containers
# in action.
#
# More information at https://github.com/gliderlabs/logspout/tree/master/httpstream

DOCKER_NETWORK=""
PARAMS=""
LOGNAME=docker.log
PORT=8500
LOG_TO_FILE=false
LOG_TO_STDOUT=true

while (( "$#" )); do
  case "$1" in
    -n|--docker-network)
      DOCKER_NETWORK=$2
      shift 2
      ;;
    -f|--log-to-file)
      LOG_TO_FILE=true
      shift 1
      ;;
    -s|--suppress-stdout)
      LOG_TO_STDOUT=false
      shift 1
      ;;
    --) # end argument parsing
      shift
      break
      ;;
    -*|--*=) # unsupported flags
      echo "Error: Unsupported flag $1" >&2
      exit 1
      ;;
    *) # preserve positional arguments
      PARAMS="$PARAMS $1"
      shift
      ;;
  esac
done
# set positional arguments in their proper place
eval set -- "$PARAMS"

if [ -z "${DOCKER_NETWORK}" ]; then
    echo "Please pick which docker network to monitor..."
    select NETWORK in $(docker network ls --format {{.Name}});
    do
      if [ "${NETWORK}" ]; then
        DOCKER_NETWORK=${NETWORK}
        break
      fi
    done
fi

echo Starting monitoring on all containers on the network ${DOCKER_NETWORK}
LOGDIR=$(pwd)/logs

docker kill logspout 2> /dev/null 1>&2 || true
docker rm logspout 2> /dev/null 1>&2 || true

docker run -d --name="logspout" \
    --volume=/var/run/docker.sock:/var/run/docker.sock \
    --volume=${LOGDIR}:/logs \
    --publish=127.0.0.1:${PORT}:80 \
    --network  ${DOCKER_NETWORK} \
    gliderlabs/logspout
sleep 3

if [ "${LOG_TO_FILE}" = true ]; then
  echo Logging to file at  ${LOGDIR}/${LOGNAME}
  docker exec -d logspout sh -c 'wget -q -O /logs/docker.log http://127.0.0.1:80/logs'
fi

if [ "${LOG_TO_STDOUT}" = true ]; then
  curl http://127.0.0.1:8500/logs
fi