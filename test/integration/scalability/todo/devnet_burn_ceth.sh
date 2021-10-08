#!/bin/bash -x

SIFCHAIN_BURN_ADDRESS=$(sifnoded keys show peggy-burn --keyring-backend test -a)
sifnoded \
    tx \
    ethbridge \
    burn \
    ${SIFCHAIN_BURN_ADDRESS} \
    0x5171050beb52148aB834Fb21E3E30FA429470c46 \
    100 \
    ceth \
    40000000000000000 \
    --node tcp://rpc-devnet-042-ibc.sifchain.finance:80 \
    --chain-id=sifchain-devnet-042-ibc \
    --keyring-backend test \
    --fees 100000rowan \
    --ethereum-chain-id=3 \
    --yes \
    --from ${SIFCHAIN_BURN_ADDRESS}