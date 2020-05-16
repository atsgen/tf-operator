#!/bin/bash
set -ex

wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gz
rm -rf go1.13.4.linux-amd64.tar.gz
echo "export GOPATH=\$HOME/go" >> ~/.bashrc
echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> ~/.bashrc
source ~/.bashrc

# get pre-compiled operator-sdk v0.15.2
wget https://github.com/operator-framework/operator-sdk/releases/download/v0.15.2/operator-sdk-v0.15.2-x86_64-linux-gnu
# copy to bin
sudo mv operator-sdk-v0.15.2-x86_64-linux-gnu /usr/local/bin/operator-sdk
sudo chmod +x /usr/local/bin/operator-sdk
operator-sdk version
