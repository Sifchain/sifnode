#!/usr/bin/env bash

set -x

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol stake \
#   --nativeAmount 100000000000000000000 \
#   --externalAmount 100000000000000000000 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

sifnoded tx clp create-pool \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdc \
  --nativeAmount 96174925423929435353999282 \
  --externalAmount 488439008293 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol ceth \
#   --nativeAmount 70163907841446304439172325 \
#   --externalAmount 313805836067779187011 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cwbtc \
#   --nativeAmount 5945801596833170595473342 \
#   --externalAmount 146909052 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2 \
#   --nativeAmount 236384238961917282723975402 \
#   --externalAmount 149104332852 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y