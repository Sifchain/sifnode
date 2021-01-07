
cd ../smart-contracts 

yarn peggy:lock $(sifnodecli keys show akasha -a) 0x0000000000000000000000000000000000000000 2000000000000000000

sleep 5

yarn advance 200

sleep 5

sifnodecli tx ethbridge burn $(sifnodecli keys show akasha -a) 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 2000000000000000000 ceth --ethereum-chain-id=5777 --from=akasha --yes

cd ../ui