set -e

# Run from setup-linux-environment.sh.  Runs all the setup
# that needs to happen as root.
# (see setup-linux-environment-user.sh for more tools)

# We need to know what user to add to the docker group, since this file
# is run with sudo
dockeruser=$1

apt-get update && apt-get install -y curl sudo lsb-release software-properties-common wget

# yarn repository
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=1 apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list

# docker repository
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=1 apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# geth
sudo add-apt-repository -y ppa:ethereum/ethereum

# nodejs
curl -sL https://deb.nodesource.com/setup_15.x | sudo -E bash -
sudo apt-get install -y nodejs

apt-get update

apt-get install -y jq make rake docker-ce docker-ce-cli containerd.io libc6-dev gcc python3-venv python3-dev python3-pip parallel netcat uuid-runtime vim tmux rsync psmisc ethereum
apt-get install -y --no-install-recommends yarn

# don't want to require root to run docker
groupadd -f docker
usermod -aG docker ${dockeruser}

curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && chmod a+x /usr/local/bin/docker-compose

# install go
wget -O /tmp/go.tar.gz https://golang.org/dl/go1.15.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf /tmp/go.tar.gz
