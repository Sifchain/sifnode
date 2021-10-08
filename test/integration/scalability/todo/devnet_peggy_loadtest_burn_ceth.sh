#!/bin/bash -x

SIFCHAIN_BURN_ADDRESS=$(sifnoded keys show devnet-peggy-loadtest-burn-eth --keyring-backend test -a)
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
        --node tcp://rpc-devnet-042.sifchain.finance:80 \
        --keyring-backend test \
        --fees 100000rowan \
        --ethereum-chain-id=3 \
        --chain-id=sifchain-devnet-042 \
        --yes \
        --from ${SIFCHAIN_BURN_ADDRESS}
done