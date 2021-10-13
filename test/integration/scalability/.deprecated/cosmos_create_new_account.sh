#!/bin/bash -x

yes | gaiad keys add $1 --keyring-backend test --output json | jq -r .mnemonic