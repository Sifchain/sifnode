#!/bin/bash -x

sifnoded keys list --keyring-backend test --output json | jq -r ".[] | {name: .name, address: .address}"