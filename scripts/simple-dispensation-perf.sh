#!/bin/zsh

#Uncomment the following lines in the address `sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd` is not present in your keyring
#echo "Generating deterministic account - sif"
#echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnoded keys add sif --recover --keyring-backend=test

for i in {1..500}
do
   sifnoded tx dispensation run test_dist Airdrop --from=sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --keyring-backend=test --yes --chain-id=localnet --gas=1000000000 --gas-prices=0.5rowan
   sleep 6
done