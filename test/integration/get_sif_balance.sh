#!/bin/bash 

addr=$1
shift

sifnodecli q auth account ${addr:=${OWNER_ADDR}} -o json | jq
