#!/usr/bin/env bash

set -x

sifnoded tx clp create-pool \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol ceth \
  --nativeAmount 65153731484359224820359378 \
  --externalAmount 386218358921701028829 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y

sifnoded tx clp create-pool \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdc \
  --nativeAmount 92424607216143501722331582 \
  --externalAmount 559730963318 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y

sifnoded tx clp create-pool \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdt \
  --nativeAmount 3511606198069702409404908 \
  --externalAmount 21260536574 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y

sifnoded tx clp create-pool \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol clink \
  --nativeAmount 455568472021115510794446 \
  --externalAmount 464241808330336890821 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y
