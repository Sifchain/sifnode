


# Multisig Key - It is a key composed of two or more keys (N) , with a signing threshold (K) ,such that the transaction needs K out of N votes to go through.
sifnoded tx dispensation claim ValidatorSubsidy --from akasha --keyring-backend test --yes --chain-id=localnet -o json
sifnoded tx dispensation claim ValidatorSubsidy --from sif --keyring-backend test --yes --chain-id=localnet -o json
# create airdrop
# mkey = multisig key
# ar1 = name for airdrop , needs to be unique for every airdrop . If not the tx gets rejected
# input.json list of funding addresses  -  Input address must be part of the multisig key
# output.json list of airdrop receivers.
sleep 8
sifnoded q dispensation claims-by-type ValidatorSubsidy -o json
sleep 8
sifnoded tx dispensation create ValidatorSubsidy output.json sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --gas 200064128 --from=sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --keyring-backend=test --fees 100000rowan --yes --chain-id=sifchain-devnet-042 --node=tcp://rpc-devnet-042.sifchain.finance:80

sleep 8
sifnoded q dispensation distributions-all --chain-id=localnet -o json
#sifnoded q dispensation records-by-name-all ar1 >> all.json
#sifnoded q dispensation records-by-name-pending ar1 >> pending.json
#sifnoded q dispensation records-by-name-completed ar1 >> completed.json
#sifnoded q dispensation records-by-addr sif1cp23ye3h49nl5ty35vewrtvsgwnuczt03jwg00


