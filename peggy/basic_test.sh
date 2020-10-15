# test the cross chain asset transfer

## Common query
sifnodecli query account $(sifnodecli keys show validator -a)
sifnodecli query account $(sifnodecli keys show testuser -a)

1. transfer token
sifnodecli tx send validator $(sifnodecli keys show testuser -a) 10stake --yes


## Case 1
1. send tx to cosmos after get the lock event in ethereum
sifnodecli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 3 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
$(sifnodecli keys show testuser -a) $(sifnodecli keys show validator -a --bech val) 5 lock \
--token-contract-address=0x0000000000000000000000000000000000000000 --ethereum-chain-id=3 --from=validator --yes

2. query the tx
sifnodecli q tx

3. query the prophecy (error now)
sifnodecli query ethbridge prophecy 0x30753E4A8aad7F8597332E813735Def5dD395028 0 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 \
--ethereum-chain-id=3 --token-contract-address=0x0000000000000000000000000000000000000000

4. 
