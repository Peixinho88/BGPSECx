#!/bin/bash

service docker start

cd network

./stop.sh

./teardown.sh

./build.sh

cd ../java

mvn install

cd target

cp blockchain-java-sdk-0.0.1-SNAPSHOT-jar-with-dependencies.jar blockchain-client.jar

cp blockchain-client.jar ../../network_resources

#cd ../../network_resources

#java -cp blockchain-client.jar main.java.org.example.network.CreateChannel

#java -cp blockchain-client.jar main.java.org.example.network.DeployInstantiateChaincode

#java -cp blockchain-client.jar main.java.org.example.user.RegisterEnrollUser
