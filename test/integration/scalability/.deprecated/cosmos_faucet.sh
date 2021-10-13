#!/bin/bash -x

curl \
    -X POST \
    -d "{\"address\": \"$(./get_cosmos_address.sh $1)\"}" \
    https://faucet.testnet.cosmos.network
