FROM debian:stable-slim
# Setup basic system dependencies
RUN apt-get update && apt-get install -y make wget nodejs npm git python3 python3-yaml
# Install Golang support
RUN wget https://golang.org/dl/go1.17.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
ENV GOBIN=/usr/local/go/bin
# Run the sifchain dev environment
RUN git clone https://github.com/Sifchain/sifnode.git && cd /sifnode && git checkout future/devenv-rebased
RUN cd /sifnode/smart-contracts && npm install
WORKDIR "/sifnode/smart-contracts" 
CMD npx hardhat run scripts/devenv.ts