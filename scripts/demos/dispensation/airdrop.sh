


# Multisig Key - It is a key composed of two or more keys (N) , with a signing threshold (K) ,such that the transaction needs K out of N votes to go through.

# create airdrop
# mkey = multisig key
# ar1 = name for airdrop , needs to be unique for every airdrop . If not the tx gets rejected
# input.json list of funding addresses  -  Input address must be part of the multisig key
# output.json list of airdrop receivers.
sifnoded tx dispensation create mkey ar1 input.json output.json --gas 200064128 --generate-only >> offlinetx.json
# First user signs
sifnoded tx sign --multisig $(sifnoded keys show mkey -a) --from $(sifnoded keys show sif -a)  offlinetx.json >> sig1.json
# Second user signs
sifnoded tx sign --multisig $(sifnoded keys show mkey -a) --from $(sifnoded keys show akasha -a)  offlinetx.json >> sig2.json
# Multisign created from the above signatures
sifnoded tx multisign offlinetx.json mkey sig1.json sig2.json >> signedtx.json
# transaction broadcast , distribution happens
sifnoded tx broadcast signedtx.json
sleep 8
sifnoded q dispensation distributions-all
sifnoded q dispensation records-by-addr sif1cp23ye3h49nl5ty35vewrtvsgwnuczt03jwg00
rm -rf offlinetx.json sig1.json sig2.json signedtx.json


