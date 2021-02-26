#!/bin/bash 

addr=$1
shift

sifnodecli q auth account ${addr:=${VALIDATOR1_ADDR}} -o json | jq
