#!/usr/bin/env bash

set -x

echo ${ADMIN_MNEMONIC} | sifnoded keys add ${ADMIN_KEY} --recover --keyring-backend=test