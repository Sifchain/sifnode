FROM ubuntu:latest
RUN apt-get update -y && \
    apt-get install -y wget curl make build-essential unzip tar git && \
    curl -O https://storage.googleapis.com/golang/go1.13.5.linux-amd64.tar.gz && \
    tar -xvf go1.13.5.linux-amd64.tar.gz && \
    chmod a+xX go && \
    mv go /usr/local && \
    echo "export GOPATH=$HOME/work" >> ~/.profile && \
    echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> ~/.profile && \
    . ~/.profile && \
    go version && \
    git clone https://github.com/binance-chain/bsc && \
    cd bsc && \
    go get github.com/btcsuite/btcd@v0.20.1-beta && \
    make geth && \
    echo "export PATH=$(pwd)/build/bin:${PATH}" >> ~/.profile