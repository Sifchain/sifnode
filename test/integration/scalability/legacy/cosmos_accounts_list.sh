#!/bin/bash -x

gaiad keys list --keyring-backend test --output json | jq -r ".[] | {name: .name, address: .address}"