


# Multisig Key - It is a key composed of two or more keys (N) , with a signing threshold (K) ,such that the transaction needs K out of N votes to go through.
sifnodecli tx dispensation claim ValidatorSubsidy --from akasha --keyring-backend test --yes
sifnodecli tx dispensation claim ValidatorSubsidy --from sif --keyring-backend test --yes
# create airdrop
# mkey = multisig key
# ar1 = name for airdrop , needs to be unique for every airdrop . If not the tx gets rejected
# input.json list of funding addresses  -  Input address must be part of the multisig key
# output.json list of airdrop receivers.
sleep 8
sifnodecli q dispensation claims-by-type ValidatorSubsidy
sleep 8
sifnoded tx dispensation create ValidatorSubsidy output.json sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --gas 200064128 --from=sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --yes --keyring-backend=test --chain-id=localnet

sleep 8
sifnoded q dispensation distributions-all --chain-id=localnet
#sifnodecli q dispensation records-by-name-all ar1 >> all.json
#sifnodecli q dispensation records-by-name-pending ar1 >> pending.json
#sifnodecli q dispensation records-by-name-completed ar1 >> completed.json
#sifnodecli q dispensation records-by-addr sif1cp23ye3h49nl5ty35vewrtvsgwnuczt03jwg00


