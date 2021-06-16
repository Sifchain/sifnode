


# Multisig Key - It is a key composed of two or more keys (N) , with a signing threshold (K) ,such that the transaction needs K out of N votes to go through.

# create airdrop
# mkey = multisig key
# ar1 = name for airdrop , needs to be unique for every airdrop . If not the tx gets rejected
# input.json list of funding addresses  -  Input address must be part of the multisig key
# output.json list of airdrop receivers.
sifnodecli tx dispensation create Airdrop output.json $(sifnodecli keys show sif -a) --from $(sifnodecli keys show sif -a) --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan
sleep 8
sifnodecli q dispensation distributions-all
#sifnodecli q dispensation records-by-name-all ar1 >> all.json
#sifnodecli q dispensation records-by-name-pending ar1 >> pending.json
#sifnodecli q dispensation records-by-name-completed ar1 >> completed.json
#sifnodecli q dispensation records-by-addr sif1cp23ye3h49nl5ty35vewrtvsgwnuczt03jwg00

sifnodecli tx dispensation create Airdrop output.json --gas 90128 --from $(sifnodecli keys show sif -a) --yes --broadcast-mode async --sequence 26 --account-number 3 --chain-id localnet
sifnodecli tx dispensation create Airdrop output.json --gas 90128 --from $(sifnodecli keys show sif -a) --yes --broadcast-mode async --sequence 27 --account-number 3 --chain-id localnet



