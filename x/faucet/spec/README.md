# README

### Summary

For the Sifchain MVP, the faucet module provides the following functionality
- Initialize a module account in genesis with token balances
- Transfer tokens from the faucet module to a requesting account
- Refill the faucet module balance from an external account
- Query the faucet module balance

## Trasactions supported

 - **Request coins**
    - sifnodecli tx faucet request-coins 1000ceth --from shadowfiend

 - **Query module account balance**
    - sifnodecli query faucet balance 

 - **Query signer account balance**
    - sifnodecli query account $(sifnodecli keys show shadowfiend -a)
    
 - **Add to faucet**
    - sifnodecli tx faucet add-coins 10000ceth --from shadowfiend