#!/usr/bin/env bash

set -x

sifnoded tx clp create-pool \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol ceth \
  --nativeAmount 75326087364173030274309330 \
  --externalAmount 2249685578540845097982 \
  --fees 100000000000000000rowan \
  --node ${SIFNODE_NODE} \
  --chain-id $SIFNODE_CHAIN_ID \
  --broadcast-mode block \
  -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cusdc \
#   --nativeAmount 52798591956187184978275830 \
#   --externalAmount 5940239555604 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol cusdt \
#   --nativeAmount 1550459183129248235861408 \
#   --externalAmount 174248776094 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2 \
#   --nativeAmount 200501596725333601567765449 \
#   --externalAmount 708998027178 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA \
#   --nativeAmount 32788415426458039601937058 \
#   --externalAmount 139140831718 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y

# sifnoded tx clp create-pool \
#   --from $SIF_ACT \
#   --keyring-backend test \
#   --symbol ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D \
#   --nativeAmount 29315228314524379224549414 \
#   --externalAmount 29441954962 \
#   --fees 100000000000000000rowan \
#   --node ${SIFNODE_NODE} \
#   --chain-id $SIFNODE_CHAIN_ID \
#   --broadcast-mode block \
#   -y