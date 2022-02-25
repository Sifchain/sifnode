#!/bin/bash

# This script returns the sequence numbers of stuck packets. Stuck packets are 
# packets for which there is still a commitment on the originating chain, and 
# which have never been received on the destination chain.
# The fact that there is still a commitment on the originating chain is 
# equivalent to the packet having been sent but never acknowledged or timedout,
# because when a packet is acknowledged or timedout, the corresponding 
# commitment is deleted

from_node="http://rpc.sifchain.finance:80"
to_node="http://public-node.terra.dev:26657"

# to find channel numbers, use this query:
# sifnoded q ibc channel connections connection-21 --node [node-rpc-url] 
from_channel="channel-18"
to_channel="channel-7"

packets=()
page=1
x=1

# Get the list of packet commitments from the originating chain. These are 
# packets that were sent but never acknowledged or timed out
while [ $x -gt 0 ]
do
  res=$(sifnoded q ibc channel packet-commitments transfer $from_channel --node $from_node --page $page --output json)
  commitments=($(echo $res | jq '.commitments[].sequence|tonumber')) 
  packets+=(${commitments[@]})
  x=${#commitments[@]}
  page=$(( $page + 1 ))
done

printf -v joined '%s,' "${packets[@]}"

# Out of the list of packet commitments retrieved aboce, check which ones were 
# not received on the destination chain.
sifnoded q ibc channel unreceived-packets transfer $to_channel --sequences="${joined%,}" --node $to_node --output json | jq '.sequences[]|tonumber'
