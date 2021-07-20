
# end to end test
## start truffle, sifnoded and relayer
1. open a console
cd smart-contracts
yarn develop

2. open other console 
cd smart-contracts
yarn migrate
sifnoded start

1. open new console
cd smart-contracts
ebrelayer init tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 user1 --chain-id=sifchain

### case 1: lock eth and send to cosmos user2 from eth operator account
1. check the balance of operator before lock
yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
2. check the ballance of contract before lock
yarn peggy:getTokenBalance  0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4  eth
3. check the user balance before lock
sifnoded query account $(sifnoded keys show user2 -a)

yarn peggy:lock $(sifnoded keys show user2 -a) 0x0000000000000000000000000000000000000000 1000000000000000000

4. check the balance of operator before lock
yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
5. check the ballance of contract before lock
yarn peggy:getTokenBalance  0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4  eth
6. check the user2 balance before lock
sifnoded query account $(sifnoded keys show user2 -a)

### case 2: burn user2's eth in cosmos then asset to back to ethereum's validator account
1. check the validator's balance before burn
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth
sifnoded query account $(sifnoded keys show user2 -a)

2. send burn tx in cosmos
sifnoded tx ethbridge burn $(sifnoded keys show user2 -a) 0xf17f52151EbEF6C7334FAD080c5704D77216b732 1000000000000000000 peggyeth --ethereum-chain-id=5777 --from=user2 --yes

3. check user2's account 
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth
sifnoded query account $(sifnoded keys show user2 -a)

### case 3: lock rowan in cosmos then issue the token in ethereum
sifnoded tx ethbridge lock $(sifnoded keys show user2 -a) 0xf17f52151EbEF6C7334FAD080c5704D77216b732 1 rwn  --ethereum-chain-id=5777 --from=user2 --yes

1. check the balance of user2 peggyatom in ethereum
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732  0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
sifnoded query account $(sifnoded keys show user2 -a)

### case 4: burn rowan in ethereum and rowan will be back to cosmos
yarn peggy:burn $(sifnoded keys show user2 -a) 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA 1
1. check balance after burn 
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732  0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
sifnoded query account $(sifnoded keys show user2 -a)