#!/bin/bash

set -eu

######################################
# "TERRA" | "SIF"
SEND_CHAIN="SIF"

SEND_CONNECTION="21" ## NOTE: 21 for Sif Terra-Sif connection; 19 for Terra Terra-Sif connection
RCV_CONNECTION="19" ## NOTE: 21 for Sif Terra-Sif connection; 19 for Terra Terra-Sif connection
#####################################

case $SEND_CHAIN in

  SIF)
    SEND_CHAIN_DIR=sif
    RCV_CHAIN_DIR=terra
    ;;

  TERRA)
    SEND_CHAIN_DIR=terra
    RCV_CHAIN_DIR=sif
    ;;

  *)
    echo -n "Unknown chain: $SEND_CHAIN"
    exit 1
    ;;
esac

SEND_CHAIN_DIR=$SEND_CHAIN_DIR/$SEND_CONNECTION
RCV_CHAIN_DIR=$RCV_CHAIN_DIR/$RCV_CONNECTION


sort clean/$RCV_CHAIN_DIR/rcv_packet_seq.data > tmp_rcvs
sort clean/$SEND_CHAIN_DIR/send_packet_seq.data  > tmp_send
sort clean/$SEND_CHAIN_DIR/timeout_packet_seq.data > tmp_timeout

comm -23 tmp_send tmp_rcvs > clean/$SEND_CHAIN_DIR/missing_rcvs.data
comm -23 clean/$SEND_CHAIN_DIR/missing_rcvs.data tmp_timeout > clean/$SEND_CHAIN_DIR/missing_txs.data

sort -no clean/$SEND_CHAIN_DIR/missing_rcvs.data clean/$SEND_CHAIN_DIR/missing_acks.data
sort -no clean/$SEND_CHAIN_DIR/missing_txs.data clean/$SEND_CHAIN_DIR/missing_txs.data

rm tmp_rcvs tmp_send tmp_timeout

