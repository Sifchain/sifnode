#!/bin/bash

set -eu

######################################
TYPE="SEND" # "SEND" | "UPDATE"
#TYPE="UPDATE"
#####################################

case $TYPE in

  SEND)
    QUERY="send_packet.packet_connection=connection-19"
    OUTPUT_DIR=send
    ;;

#   ACK)
#     QUERY="acknowledge_packet.packet_connection=connection-19"
#     OUTPUT_DIR=update_client
#     ;;

  UPDATE)
    QUERY="update_client.client_id=07-tendermint-19"
    OUTPUT_DIR=update_client
    ;;

  *)
    echo -n "Unknown query type: $TYPE"
    exit 1
    ;;
esac

NODE="http://public-node.terra.dev:26657"

get_num_pages () {
    echo "Calculating number of pages"
    terrad query txs --events $QUERY --node $NODE --output json  --page 1 --limit=1 > tmp
    TOTAL_COUNT=$(cat tmp | jq '.total_count|tonumber')
    NUM_PAGES=$(( ($TOTAL_COUNT + (30 - 1)) / 30))
    ###############rm tmp
}

get_num_pages
echo "Total number of pages: $NUM_PAGES"

mkdir -p data/$OUTPUT_DIR

############################
#NUM_PAGES=1
############################

for ((i=1; i <= $NUM_PAGES; i++));
    do
       echo "Getting page $i"
       terrad query txs --events $QUERY --node $NODE --page $i --output json > data/$OUTPUT_DIR/$i.json
 done


