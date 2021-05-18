#!/bin/bash

set -e

#Install cURL
sudo apt update
sudo apt install curl

#Remove previous Docker installations
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get purge docker-ce docker-ce-cli containerd.io
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd
sudo rm /usr/local/bin/docker-compose
sudo apt remove docker-compose
sudo apt autoremove

#Remove previous Golang installations
sudo apt remove golang
sudo apt autoremove
sudo rm -rf /usr/local/go
sudo nano ~/.bashrc
source ~/.bashrc

#Install Docker
sudo apt-get update
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

#Do I need Docker Compose?
sudo apt update
sudo curl -L "https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

#Install Docker Engine
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io

#Install Golang
sudo wget https://golang.org/dl/go1.15.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
source ~/.bashrc

#Alternative Golang installation process
#sudo snap install --classic --channel=1.15/stable go

#Install Node JS and NPM
sudo apt install nodejs
sudo apt install npm
#npm install npm@5.6.0 -g #Upgrade NPM

#Install Python
sudo apt-get install python

#Install Maven
sudo apt update
sudo apt install maven

#Install Java
sudo apt update
sudo apt install default-jdk
#Maybe install JRE specifically? Not sure...


