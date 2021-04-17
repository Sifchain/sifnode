#!/bin/bash
. ../credentials.sh

yarn && yarn ganache-cli \
  -m "$ETHEREUM_ROOT_MNEMONIC" \
  --db ~/.ganachedb \
  -p 7545 \
  --networkId 5777 \
  -g 20000000000 \
  --gasLimit 6721975