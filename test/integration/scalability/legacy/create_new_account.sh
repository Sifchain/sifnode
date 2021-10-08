#!/bin/bash -x

yes | sifnoded keys add $1 --keyring-backend test --output json | jq -r .mnemonic