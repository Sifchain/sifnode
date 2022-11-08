


# Multisig Key - It is a key composed of two or more keys (N) , with a signing threshold (K) ,such that the transaction needs K out of N votes to go through.

# create airdrop
# mkey = multisig key
# ar1 = name for airdrop , needs to be unique for every airdrop . If not the tx gets rejected
# input.json list of funding addresses  -  Input address must be part of the multisig key
# output.json list of airdrop receivers.
sifnoded tx dispensation create ValidatorSubsidy output.json sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --yes --fees 150000rowan --chain-id=localnet --keyring-backend=test
sifnoded tx dispensation run 29_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd ValidatorSubsidy--from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --yes --fees 150000rowan --chain-id=localnet --keyring-backend=test
sleep 8
sifnoded q dispensation distributions-all -chain-id localnet
#sifnoded q dispensation records-by-name-all ar1 >> all.json
#sifnoded q dispensation records-by-name-pending ar1 >> pending.json
#sifnoded q dispensation records-by-name-completed ar1 >> completed.json
#sifnoded q dispensation records-by-addr sif1cp23ye3h49nl5ty35vewrtvsgwnuczt03jwg00

sifnoded tx dispensation create Airdrop output.json --gas 90128 --from $(sifnoded keys show sif -a) --yes --broadcast-mode async --sequence 26 --account-number 3 --chain-id localnet
sifnoded tx dispensation create Airdrop output.json --gas 90128 --from $(sifnoded keys show sif -a) --yes --broadcast-mode async --sequence 27 --account-number 3 --chain-id localnet
sifnoded tx dispensation run 25_sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd ValidatorSubsidy --from sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --yes --gas auto --gas-adjustment=1.5 --gas-prices 1.0rowan --chain-id=localnet --keyring-backend=test



