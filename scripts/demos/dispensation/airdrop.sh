


# Multisig Key - It is a key composed of two or more keys (N) , with a signing threshold (K) ,such that the transaction needs K out of N votes to go through.

# create airdrop
# mkey = multisig key
# ar1 = name for airdrop , needs to be unique for every airdrop . If not the tx gets rejected
# input.json list of funding addresses  -  Input address must be part of the multisig key
# output.json list of airdrop receivers.
sifnodecli tx dispensation create mkey ar1 Airdrop input.json output.json --gas 200064128 --generate-only > offlinetx.json
# First user signs
sifnodecli tx sign --multisig $(sifnodecli keys show mkey -a) --from $(sifnodecli keys show sif -a)  offlinetx.json > sig1.json
# Second user signs
sifnodecli tx sign --multisig $(sifnodecli keys show mkey -a) --from $(sifnodecli keys show akasha -a)  offlinetx.json > sig2.json
# Multisign created from the above signatures
sifnodecli tx multisign offlinetx.json mkey sig1.json sig2.json > signedtx.json
# transaction broadcast , distribution happens
sifnodecli tx broadcast signedtx.json
sleep 8
sifnodecli q dispensation distributions-all
sifnodecli q dispensation records-by-name-all ar1 >> all.json
sifnodecli q dispensation records-by-name-pending ar1 >> pending.json
sifnodecli q dispensation records-by-name-completed ar1 >> completed.json
sifnodecli q dispensation records-by-addr sif1cp23ye3h49nl5ty35vewrtvsgwnuczt03jwg00
rm -rf offlinetx.json sig1.json sig2.json signedtx.json


