#!/bin/bash -x

SIFCHAIN_LOCK_ADDRESS=$(sifnoded keys show testnet-peggy-loadtest-burn-erowan --keyring-backend test -a)
for i in {1..10}
do
    yarn \
        -s \
        --cwd /sifnode/smart-contracts \
        integrationtest:sendBurnTx \
        --sifchain_address ${SIFCHAIN_LOCK_ADDRESS} \
        --bridgebank_address 0xDC959a7cad365F22DeB2d70Dfbb2f4974BdAcDbf \
        --symbol 0xd3c2e9b5539A056EE072Dc8Dd1D3Bbb1A9215B88 \
        --ethereum_private_key_env_var ETHEREUM_PRIVATE_KEY \
        --json_path /sifnode/smart-contracts/deployments/sifchain-testnet-042-ibc \
        --gas estimate \
        --ethereum_network ropsten \
        --ethereum_address 0x5171050beb52148aB834Fb21E3E30FA429470c46 \
        --amount 1
        # --symbol 0xE4A3869e481F2C0A964d739929d65f6627b99324 \
        # --bridgebank_address 0xB75849afEF2864977a858073458Cb13F9410f8e5 \
done