#!/usr/bin/env bash

sudo apt update && sudo apt-get install -y git jq

sudo chown -R kitchen. /home/kitchen/go
#
#sudo wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
#sudo tar -xvf go1.13.3.linux-amd64.tar.gz
#sudo mv go /usr/local

#cd /home/kitchen/go/src/

echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile