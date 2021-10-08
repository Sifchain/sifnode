#!/bin/bash -x

SIFCHAIN_LOCK_ADDRESS=$(sifnoded keys show testnet-peggy-loadtest-lock-rowan --keyring-backend test -a)
for i in {1..100}
do
    sifnoded \
        tx \
        ethbridge \
        lock \
        ${SIFCHAIN_LOCK_ADDRESS} \
        0x5171050beb52148aB834Fb21E3E30FA429470c46 \
        10000 \
        rowan \
        40000000000000000 \
        --node tcp://rpc-testnet-042-ibc.sifchain.finance:80 \
        --keyring-backend test \
        --fees 100000rowan \
        --ethereum-chain-id=3 \
        --chain-id=sifchain-testnet-042-ibc \
        --yes \
        --from ${SIFCHAIN_LOCK_ADDRESS}
done