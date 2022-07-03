# siftool

Prerequisite
- gcc
- python3
- python3-dev
- docker
- abigen

# Setup on Ubuntu 22.04
- sudo apt update
- sudo apt upgrade -y
- sudo apt install -y gcc make python3-dev python3-venv golang
- curl -sL https://deb.nodesource.com/setup_16.x | sudo bash -
- sudo apt install nodejs
- Install geth (for peggy2)
  sudo add-apt-repository -y ppa:ethereum/ethereum
  sudo apt-get update
  sudo apt-get install ethereum
- Install Docker (for peggy2):
  sudo apt-get install ca-certificates curl gnupg lsb-release
  sudo mkdir -p /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
  sudo apt-get update
  sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
  sudo usermod -aG docker $USER

To start the local environment:

siftool run-env

It will automatically install Python dependencies upon first use. This command will detect if you are on Peggy1 or
Peggy2 branch, and will start local processes accordingly:
- For Peggy1, it will run ganache-cli, sifnoded and ebrelayer.
- For Peggy2, it will run hardhat, sifnoded and two instances of ebrelayer.

At the moment, the environment consists of Ethereum-compliant local node (ganache/hardhat), one `sifnode` validator and
a Peggy bridge implemented by `ebrelayer` binary.


Original design document: https://docs.google.com/document/d/1IhE2Y03Z48ROmTwO9-J_0x_lx2vIOFkyDFG7BkAIqCk/edit#
