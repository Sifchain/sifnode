#!/bin/bash

# convert from contracts/BridgeBank/BridgeBank.sol to artifacts/contracts/BridgeBank/BridgeBank.sol/BridgeBank.json
for f in $*
do
  base=$(basename $f .sol) # extracts BrideBank from contracts/BridgeBank/BridgeBank.sol
  echo artifacts/$f/${base}.json
done
