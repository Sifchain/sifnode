#!/bin/bash

set -eu

# TODO: this is incomplete - doesn't take into account timed out packets

######################################
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

comm -23 clean/$CHAIN_DIR/send_packet_seq.data clean/$CHAIN_DIR/ack_packet_seq.data > clean/$CHAIN_DIR/missing_txs.data