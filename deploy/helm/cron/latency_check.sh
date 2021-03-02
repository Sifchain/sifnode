#!/bin/sh

CHECK_LATENCY () {
#set the latest_block_timestamp
DOMAIN=$1
latest_block_time=$(curl -s $DOMAIN | jq '.result.sync_info.latest_block_time' | sed 's/"//g' | awk -F '.' '{print $1}' |sed 's/T/ /')

if [ -z "$latest_block_time" ]
then
      echo "didn't retrieve the latest_block_time"
      return
fi
latest_block_ts=$(date -u -d "$latest_block_time" +%s)
TAG=$(curl -s $DOMAIN | jq '.result.node_info.moniker' | sed 's/"//g')
#set the current timestamp
current_ts=$(date -u +%s)
#detect empty variable for latest_block_ts
if [ -z "$latest_block_ts" ]
then
      echo "didn't retrieve the latest_block_ts"
      exit 3
fi

NOW=$(date -u +%s)
latency=$(expr $current_ts - $latest_block_ts )


curl -X POST "https://api.datadoghq.com/api/v1/series?api_key=b4af28c08b859e010b40c39bf8f357a4" \
-H "Content-Type: application/json" \
-d @- << EOF
{
  "series": [
    {
      "interval": 60,
      "tags": [
        "environment:${TAG}"
      ],
      "type": "count",
      "unit": "second",
      "metric": "sifchain.rpc.block_latency",
      "points": [
        [
          "${NOW}",
          "${latency}"
        ]
      ]
    }
  ]
}
EOF
}


for endpoint in "$@"
do
  echo $endpoint
  CHECK_LATENCY $endpoint
done
