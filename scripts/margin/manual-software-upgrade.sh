#!/usr/bin/env bash

set -x

sifnoded tx gov submit-proposal software-upgrade "${NEW_VERSION}" \
  --from ${SIF_ACT} \
  --deposit "${DEPOSIT}" \
  --upgrade-height "${TARGET_BLOCK}" \
  --title "v${NEW_VERSION}" \
  --description "v${NEW_VERSION}" \
  --chain-id "${SIFNODE_CHAIN_ID}" \
  --node "${SIFNODE_NODE}" \
  --keyring-backend "test" \
  --fees 100000000000000000rowan \
  --broadcast-mode=block \
  -y