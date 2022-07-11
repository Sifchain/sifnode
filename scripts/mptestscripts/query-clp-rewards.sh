#!/usr/bin/env bash

set -x

sifnoded q clp lp \
  ceth \
  sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \

sifnoded q clp lp \
  cusdc \
  sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \

sifnoded q clp lp \
  cusdt \
  sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \

sifnoded q clp lp \
  clink \
  sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \

sifnoded q clp lp \
  ceth \
  sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \

sifnoded q clp lp \
  cusdt \
  sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \

sifnoded q clp lp \
  cusdc \
  sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5 \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
