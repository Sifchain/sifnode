
# end to end test
## start truffle, sifnoded and relayer
1. open a console
cd testnet-contracts
yarn develop

2. open other console 
cd testnet-contracts
yarn migrate
yarn peggy:setup
sifnoded start

3. open new console
cd testnet-contracts
ebrelayer init tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 user1 --chain-id=sifchain

### case 1: lock eth and send to cosmos test user from eth operator account
1. check the balance of operator before lock
yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
2. check the ballance of contract before lock
yarn peggy:getTokenBalance  0x30753E4A8aad7F8597332E813735Def5dD395028  eth
3. check the user balance before lock
sifnodecli query account $(sifnodecli keys show user2 -a)

yarn peggy:lock $(sifnodecli keys show user2 -a) 0x0000000000000000000000000000000000000000 1000000000000000000

4. check the balance of operator before lock
yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
5. check the ballance of contract before lock
yarn peggy:getTokenBalance  0x30753E4A8aad7F8597332E813735Def5dD395028  eth
6. check the user2 balance before lock
sifnodecli query account $(sifnodecli keys show user2 -a)

### case 2: burn user2's eth in cosmos then asset to back to ethereum's validator account
1. check the validator's balance before burn
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth
sifnodecli query account $(sifnodecli keys show user2 -a)

2. send burn tx in cosmos
sifnodecli tx ethbridge burn $(sifnodecli keys show user2 -a) 0xf17f52151EbEF6C7334FAD080c5704D77216b732 1000000000000000000 peggyeth --ethereum-chain-id=5777 --from=user2 --yes

3. check user2's account 
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth
sifnodecli query account $(sifnodecli keys show user2 -a)

### case 3: lock atom in cosmos then issue the token in ethereum
sifnodecli tx ethbridge lock $(sifnodecli keys show user2 -a) 0xf17f52151EbEF6C7334FAD080c5704D77216b732 1 rwn  --ethereum-chain-id=5777 --from=user2 --yes

1. check the balance of validator peggyatom in ethereum
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732  0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
sifnodecli query account $(sifnodecli keys show user2 -a)

### case 4: burn atom in ethereum and atom will be back to cosmos
yarn peggy:burn $(sifnodecli keys show user2 -a) 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA 1
1. check balance after burn 
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732  0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
sifnodecli query account $(sifnodecli keys show user2 -a)