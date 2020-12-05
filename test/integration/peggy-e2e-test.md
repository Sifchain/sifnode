
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

ebrelayer init tcp://localhost:26657 ws://localhost:7545/ \
0xEcFcaB0A285d3380E488A39B4BB21e777f8A4EaC \ # BridgeRegistry contract address
akasha <akasha's mnemonic> --chain-id=sifchain

### case 1: lock eth and send to cosmos user2 from eth operator account
1. check the balance of operator before lock
yarn peggy:getTokenBalance \
0xf17f52151EbEF6C7334FAD080c5704D77216b732  \ # any ethereum account
eth
2. check the ballance of contract before lock
yarn peggy:getTokenBalance \
0x75c35C980C0d37ef46DF04d31A140b65503c0eEd \ # BridgeBank contract address
eth
3. check the user balance before lock
sifnodecli query account $(sifnodecli keys show akasha -a)

yarn peggy:lock $(sifnodecli keys show akasha -a) 0x0000000000000000000000000000000000000000 100 0xf17f52151EbEF6C7334FAD080c5704D77216b732

4. check the balance of operator after lock
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 \ # any ethereum account
eth
5. check the ballance of contract after lock
yarn peggy:getTokenBalance 0x75c35C980C0d37ef46DF04d31A140b65503c0eEd \ # BridgeBank contract address
eth
6. check the cosmos user balance after lock
sifnodecli query account $(sifnodecli keys show akasha -a)

### case 2: burn user's ceth in cosmos then unlock eth back to ethereum's receiver account
1. check the receiver's balance before burn
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 \ # any ethereum receiver account
eth
sifnodecli query account $(sifnodecli keys show akasha -a)

2. send burn tx in cosmos
sifnodecli tx ethbridge burn $(sifnodecli keys show akasha -a) \ # cosmos account from
0xf17f52151EbEF6C7334FAD080c5704D77216b732 \ # ethereum receiver account
100 ceth --ethereum-chain-id=5777 \ # web3.eth.getChainId()
--from=akasha --yes

3. check both accounts
yarn peggy:getTokenBalance 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
sifnodecli query account $(sifnodecli keys show akasha -a)

### case 3: lock rowan in cosmos then issue the token in ethereum
1. sifnodecli tx ethbridge lock $(sifnodecli keys show akasha -a) 0xf17f52151EbEF6C7334FAD080c5704D77216b732 \ # receiver account in ethereum network
1 rwn --ethereum-chain-id=5777 --from=akasha --yes

2. get newly created "eRWN" token address
yarn peggy:getTokenAddress eRWN

3. check both balances rwn and eRWN
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eRWN
sifnodecli query account $(sifnodecli keys show akasha -a)

### case 4: burn rowan in ethereum and rowan will be back to cosmos
1. yarn peggy:burn $(sifnodecli keys show akasha -a) eRWN 1 0xf17f52151EbEF6C7334FAD080c5704D77216b732
2. check balance after burn
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eRWN
sifnodecli query account $(sifnodecli keys show akasha -a)