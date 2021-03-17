# Ceth Receiver Account Setting
If Sifchain account send lock/burn transaction, then relayers need unpeg the token and send a transaction to smart contract in Ethereum. So Sifchain need charge the transaction fee from lock/burn transaction in Sifchain. All Ceth in the lock/burn transaction will go to a special account. In the doccument, we will explain how to set the account in genesis, update the account via a transaction and some concerns like how to upgrade the network.

## The account in genesis
After the ceth receiver account feature deployed to network, all ceth will go to this account. If the account not set, then all lock/burn transaction will failed. We must set it in genesis, the command as following:

1. command format 
sifnoded  set-genesis-ceth-receiver-account [ceth-receiver-account]

2. command example
sifnoded  set-genesis-ceth-receiver-account $(sifnodecli keys show sif -a)

## Update the account 
At the runtime, we can update the account via a transaction. The transaction is priviledged, only admin account can set the account. The admin account is the same as the one can update the white list validators. The command as following:
1. command format
sifnodecli tx ethbridge update_ceth_receiver_account [cosmos-sender] [ceth-receiver-account]

2. command example
sifnodecli tx ethbridge update_ceth_receiver_account sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --from=sif --yes

## Upgrade
If there is alive network need this feature, we must stop the sifchain nodes. Then set the account in genesis and upgrade the software.
