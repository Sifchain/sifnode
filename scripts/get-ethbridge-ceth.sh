#!/usr/bin/env bash

# the script get the ceth amount of ethbridge module account.
# sdk.AccAddress(crypto.AddressHash([]byte("ethbridge")))
sifnodecli  q account sif1l3dftf499u4gvdeuuzdl2pgv4f0xdtnuuwlzp8 | jq '. | {coins: .value.coins} ' | jq '.[] | map(select(.denom == "ceth"))' | jq '.[]'
