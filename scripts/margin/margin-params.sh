#!/usr/bin/env bash

set -x

sifnoded q params subspace margin Pools \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID