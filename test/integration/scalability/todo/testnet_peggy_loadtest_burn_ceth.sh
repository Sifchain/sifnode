#!/bin/bash -x

SIFCHAIN_BURN_ADDRESS=$(sifnoded keys show testnet-peggy-loadtest-burn-eth --keyring-backend test -a)
for i in {1..100}
do
    sifnoded \
        tx \
        ethbridge \
        burn \
        ${SIFCHAIN_BURN_ADDRESS} \
        0x5171050beb52148aB834Fb21E3E30FA429470c46 \
        100 \
        ceth \
        40000000000000000 \
        --node tcp://rpc-testnet-042-ibc.sifchain.finance:80 \
        --chain-id=sifchain-testnet-042-ibc \
        --keyring-backend test \
        --fees 100000rowan \
        --ethereum-chain-id=3 \
        --yes \
        --from ${SIFCHAIN_BURN_ADDRESS}
done