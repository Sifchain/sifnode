#!/bin/zsh

sed -i -s 's/        "voting_period": "172800s"/        "voting_period": "60s"/g' ~/.sifnode-1/config/genesis.json
sed -i -s 's/        "voting_period": "172800s"/        "voting_period": "60s"/g' ~/.sifnode-2/config/genesis.json
sed -i -s 's/        "voting_period": "172800s"/        "voting_period": "60s"/g' ~/.sifnode-3/config/genesis.json
