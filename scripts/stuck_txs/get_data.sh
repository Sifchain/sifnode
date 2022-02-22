#!/bin/bash

set -eu

get_num_pages () {
    terrad query txs --events update_client.client_id='07-tendermint-19' --node http://public-node.terra.dev:26657 --output json  --page 1 --limit=1 > tmp
    TOTAL_COUNT=$(cat tmp | jq '.total_count|tonumber')
    NUM_PAGES=$(( ($TOTAL_COUNT + (30 - 1)) / 30))
}

get_num_pages
echo "Total number of pages: $NUM_PAGES"

mkdir -p data

for ((i=1; i <= $NUM_PAGES; i++));
    do
       echo "Getting page $i"
       terrad query txs --events update_client.client_id='07-tendermint-19' --node http://public-node.terra.dev:26657 --page $i --output json > data/terra_query_txs_update_client.client_id_07_tendermint_19_page_$i.json
 done


