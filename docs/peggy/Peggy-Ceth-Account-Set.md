# Ceth Receiver Account Setting
If Sifchain account send lock/burn transaction, then relayers need unpeg the token and send a transaction to smart contract in Ethereum. So Sifchain need charge the transaction fee from lock/burn transaction in Sifchain. All Ceth in the lock/burn transaction will go to a special account. In the doccument, we will explain how to set the account via a transaction.

## Set the account 
At the runtime, we can update the account via a transaction. The transaction is priviledged, only admin account can set the account. The admin account is the same as the one can update the white list validators. The command as following:
1. command format
sifnodecli tx ethbridge update_ceth_receiver_account [cosmos-sender] [ceth-receiver-account]

2. command example
sifnodecli tx ethbridge update_ceth_receiver_account sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --from=sif --yes

## Rescue the Ceth
Before set the account, some Ceth is locked in the ethbridge module, we need a method to rescue the Ceth. Similar to account setting method, the transaction is priviledged, only admin account can call it. It will transfer the Ceth from ethbridge module to an specific account. The command as following:
1. command format
sifnodecli tx ethbridge rescue_ceth [cosmos-sender-address] [ceth_receiver_account] [ceth_amount]

2. command example
sifnodecli tx ethbridge rescue_ceth sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd  10000000 --from=sif --yes
