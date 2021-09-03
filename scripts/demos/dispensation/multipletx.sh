#!/usr/bin/env bash

# Use sifnoded q account $(sifnoded keys show sif -a) to get seq
seq=1
sifnoded tx dispensation create Airdrop output.json --gas 90128 --from $(sifnoded keys show sif -a) --yes --broadcast-mode async --sequence $seq --account-number 3 --chain-id localnet
seq=$((seq+1))
sifnoded tx dispensation create ValidatorSubsidy output.json --gas 90128 --from $(sifnoded keys show sif -a) --yes --broadcast-mode async --sequence $seq --account-number 3 --chain-id localnet
seq=$((seq+1))
sifnoded tx dispensation create ValidatorSubsidy output.json --gas 90128 --from $(sifnoded keys show sif -a) --yes --broadcast-mode async --sequence $seq --account-number 3 --chain-id localnet