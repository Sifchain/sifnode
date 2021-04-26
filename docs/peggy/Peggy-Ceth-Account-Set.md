# Ceth Receiver Account Setting
When a Sifchain account sends a lock/burn transaction, the relayers need to respond to that and send a transaction to an Ethereum smart contract.
When that happens, Sifchain needs to charge a transaction fee. The cEth in the lock/burn transaction will go to a defined account.
That account can be changed using sifnodecli commands, and the cEth in that account can be sent to other accounts.

## Set the account 
At runtime, we can update the account using transaction.
The transaction is priviledged, and only the admin account can set it.
The admin account is the same account that updates the white list validators. The command is as follows:

1. command format
```
sifnodecli tx ethbridge update_ceth_receiver_account $clp_admin_address $ceth-receiver-account
```

2. command example
sifnodecli tx ethbridge update_ceth_receiver_account sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd --from=sif --yes

## Rescue the Ceth
Before setting the account, some cEth can be locked in the ethbridge module, so we need a method to rescue the cEth.
Similar to the account setting method, the transaction is priviledged and only the admin account can call it.
It will transfer the cEth from ethbridge module to an specific account.

1. command format
```
sifnodecli tx ethbridge rescue_ceth $clp_admin_address $ceth-receiver-account $ceth_amount
```

2. command example
sifnodecli tx ethbridge rescue_ceth sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd  10000000 --from=sif --yes
