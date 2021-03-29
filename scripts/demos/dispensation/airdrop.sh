#!/usr/bin/env bash


# Multisig Key - It is a key composed of two or more keys (N) , with a signing threshold (K) ,such that the transaction needs K out of N votes to go through.

# create airdrop
# mkey = multisig key
# ar1 = name for airdrop , needs to be unique for every airdrop . If not the tx gets rejected
# input.json list of funding addresses  -  Input address must be part of the multisig key
# output.json list of airdrop receivers.
sifnodecli tx dispensation airdrop mkey ar1 input.json output.json --generate-only >> offlinetx.json
# First user signs
sifnodecli tx sign --multisig $(sifnodecli keys show mkey -a) --from $(sifnodecli keys show sif -a)  offlinetx.json >> sig1.json
# Second user signs
sifnodecli tx sign --multisig $(sifnodecli keys show mkey -a) --from $(sifnodecli keys show akasha -a)  offlinetx.json >> sig2.json
# Multisign created from the above signatures
sifnodecli tx multisign offlinetx.json mkey sig1.json sig2.json >> signedtx.json
# transaction broadcast , distribution happens
sifnodecli tx broadcast signedtx.json
rm -rf offlinetx.json sig1.json sig2.json signedtx.json


