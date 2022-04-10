#!/usr/bin/env bash

set -x

echo ${ADMIN_MNEMONIC} | sifnoded keys add ${SIF_ACT} --recover --keyring-backend=test