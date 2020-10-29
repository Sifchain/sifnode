#!/usr/bin/env bash

MNEMONIC=$(grep MNEMONIC .env | cut -d '=' -f 2-)
BLOCKSPEED=0
ganache-cli -m "$MNEMONIC"  -b "$BLOCKSPEED"
