#!/usr/bin/env bash

set -x

sifnoded tx margin whitelist $(sifnoded keys show tester1 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester2 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester3 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester4 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester5 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester6 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester7 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester8 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester9 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester10 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
sifnoded tx margin whitelist $(sifnoded keys show tester11 --keyring-backend=test -a) \
  --from $ADMIN_KEY \
  --keyring-backend test \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y
