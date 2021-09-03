#!/bin/bash 

addr=$1
shift

sifnoded q auth account ${addr:=${VALIDATOR1_ADDR}} -o json | jq
