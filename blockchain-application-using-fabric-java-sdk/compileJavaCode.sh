#!/bin/bash

set -e

cd java

mvn install

cd target

cp blockchain-java-sdk-0.0.1-SNAPSHOT-jar-with-dependencies.jar blockchain-client.jar

cp blockchain-client.jar ../../network_resources
