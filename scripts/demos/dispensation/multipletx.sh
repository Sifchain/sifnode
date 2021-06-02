#!/usr/bin/env bash

# Use sifnodecli q account $(sifnodecli keys show sif -a) to get seq
seq=31
sifnodecli tx dispensation create Airdrop output.json --gas 90128 --from $(sifnodecli keys show sif -a) --yes --broadcast-mode async --sequence $seq --account-number 3 --chain-id localnet
seq=$((seq+1))
sifnodecli tx dispensation create ValidatorSubsidy output.json --gas 90128 --from $(sifnodecli keys show sif -a) --yes --broadcast-mode async --sequence $seq --account-number 3 --chain-id localnet
seq=$((seq+1))
sifnodecli tx dispensation create ValidatorSubsidy output.json --gas 90128 --from $(sifnodecli keys show sif -a) --yes --broadcast-mode async --sequence $seq --account-number 3 --chain-id localnet