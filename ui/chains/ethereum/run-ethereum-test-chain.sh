#!/usr/bin/env bash

ENV_MNEMONIC=$(grep MNEMONIC .env | cut -d '=' -f 2-)
BLOCKSPEED=0
MNEMONIC=${ENV_MNEMONIC:-race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow}
yarn ganache-cli -m "$MNEMONIC"  -b "$BLOCKSPEED"
