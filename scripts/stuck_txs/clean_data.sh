#!/bin/bash

set -eu

######################################
# "SEND" | "ACK" | "TIMEOUT" | "UPDATE"
TYPE="ACK"

# "TERRA" | "SIF"
CHAIN="SIF"

CONNECTION="18" ## NOTE: 21 for Sif Terra-Sif connection; 19 for Terra Terra-Sif connection
#####################################

case $CHAIN in

  SIF)
    CHAIN_DIR=sif
    ;;

  TERRA)
    CHAIN_DIR=terra
    ;;

  *)
    echo -n "Unknown chain: $CHAIN"
    exit 1
    ;;
esac

CHAIN_DIR=$CHAIN_DIR/$CONNECTION

case $TYPE in

  SEND)
    OUT_FILE="clean/$CHAIN_DIR/send_packet_seq.data"
    FILES="data/$CHAIN_DIR/send/*"
    QUERY='.txs[].logs[].events[]|select(.type=="send_packet").attributes[]|select(.key=="packet_sequence").value|tonumber'
    ;;

  ACK)
    OUT_FILE="clean/$CHAIN_DIR/ack_packet_seq.data"
    FILES="data/$CHAIN_DIR/ack/*"
    QUERY='.txs[].logs[].events[]|select(.type=="acknowledge_packet").attributes[]|select(.key=="packet_sequence").value|tonumber'
    ;;

  TIMEOUT)
    OUT_FILE="clean/$CHAIN_DIR/timeout_packet_seq.data"
    FILES="data/$CHAIN_DIR/timeout/*"
    QUERY='.txs[].logs[].events[]|select(.type=="timeout_packet").attributes[]|select(.key=="packet_sequence").value|tonumber'
    ;;

  *)
    echo -n "Unknown type: $TYPE"
    exit 1
    ;;
esac

mkdir -p clean/$CHAIN_DIR


rm -f $OUT_FILE

for f in $FILES
do
  jq $QUERY $f >> $OUT_FILE
done

sort -no $OUT_FILE $OUT_FILE



