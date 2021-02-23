#!/bin/bash

# Run from setup-linux-environment.sh.  Runs all the setup
# that needs to happen as non-root.
# (see setup-linux-environment-root.sh for more tools)

mkdir -p ~/.npm-global/lib
npm config set prefix '~/.npm-global'

# npm install of these succeeds, but then returns 1 as its exit value.  Just
# assume it worked; if it didn't, everything will die immediately
sudo npm install -g truffle @truffle/hdwallet-provider ganache-cli || true

# these npm packages were written correctly
sudo npm install -g dotenv

# set up environment vars in .bash_profile
echo 'export GOPATH=$HOME/go' >> ~/.bash_profile
echo 'export GOBIN=$GOPATH/bin' >> ~/.bash_profile
echo 'export GOROOT=/usr/local/go' >> ~/.bash_profile
echo 'export PATH=$GOROOT/bin:$PATH' >> ~/.bash_profile
echo 'export PATH=$GOBIN:$PATH' >> ~/.bash_profile
echo 'export PATH=$HOME/.npm-global/bin:$PATH' >> ~/.bash_profile

echo '. ~/.bash_profile' >> ~/.bashrc

. ~/.bash_profile
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.33.0

python3 -m pip install -U pytest
