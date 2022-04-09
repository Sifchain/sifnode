#!/usr/bin/env bash

set -x

SIF_ACCT=$(echo ${ADMIN_MNEMONIC} | sifnoded keys add sif --recover --keyring-backend=test)
echo $SIF_ACCT