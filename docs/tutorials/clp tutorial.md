# Sifchain - Clp Basics Tutorial

#### demo video

* https://youtu.be/B2cn9Aag3sg

#### Previous tutorial 

* Peggy ethBridge: https://github.com/Sifchain/sifnode/blob/develop/docs/tutorials/peggy%20tutorial.md

#### Dependencies:

    0. `git clone git@github.com:Sifchain/sifnode.git`
        

#### What are they

Continuous liquidity pools are a way to pool assets that can then be used in a decentralised blockchain to enable the exchange/swapping from one asset to another without the need for a private off chain exchange. At the sametime providing a yield/return to the liquidity providers based on the pool units each provider has within a pool.

When used with the use of peg-zone as demonstrated a past video, this will enable cross chain swaps from one peg-zone to another. 

#### Setup 

1. Initialize the local chain run; `./scripts/init.sh`

2. Start the chain; `./scripts/run.sh`

3. Check to see you have two local accounts/keys setup; `sifnoded keys list`

```
[
  {
    "name": "akasha",
    "type": "local",
    "address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
    "pubkey": "sifpub1addwnpepqdycrc8usnjh0yk7cd532ushualgsderdqj8jr9m2rzy8stqrlpj5vymlww"
  },
  {
    "name": "shadowfiend",
    "type": "local",
    "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
    "pubkey": "sifpub1addwnpepqt6sfvz3mwetudyaxjn958kztxz9j8rvrlsu55fw6fjkjyac2s9z5sc8npe"
  }
]
```

4. Check your seed account balance/s;
   `sifnoded q account $(sifnoded keys show sif -a)`
   `sifnoded q account $(sifnoded keys show akasha -a)`
   
#### Create and query pools

1. Create the first pool; `sifnoded tx clp create-pool --from sif --symbol chot --nativeAmount 20000 --externalAmount 20000`

2. Create another pool with a different account `sifnoded tx clp create-pool --from akasha --symbol clink --nativeAmount 30000 --externalAmount 30000`    

3. Check funds left on first account; `sifnoded q account $(sifnoded keys show sif -a)`

4. Check funds left on second account; `sifnoded q account $(sifnoded keys show akasha -a)`

5. Query all clp pools; `sifnoded q clp pools`

6. Query the ceth pool; `sifnoded q clp pool ceth`

7. Query an accounts liquidity provider `sifnoded q clp lp chot $(sifnoded keys show sif -a)`

#### Add Extra liquidity  (Continuing from above)

1. Add more liquidity; `sifnoded tx clp add-liquidity --from sif --symbol chot --nativeAmount 10000 --externalAmount 10000` 

2. Add more liquidity from other account ; `sifnoded tx clp add-liquidity --from akasha --symbol clink --nativeAmount 10000 --externalAmount 10000`

#### Swap via the pools 

1.  Swap some chot for clink via the sif key/account; `sifnoded tx clp swap --from sif --sentSymbol chot --receivedSymbol clink --sentAmount 200` 

2. Swap some clink for chot via the akasha key/account; `sifnoded tx clp swap --from akasha --sentSymbol clink --receivedSymbol chot --sentAmount 200`

#### Removing liquidity (Continuing from above)

### Basic Options 
 
```--asymmetry```         -10000 = 100% Native Asset, 0 = 50% Native Asset 50% External Asset, 10000 = 100% External Asset

```--wBasis```            0 = 0%, 10000 = 100%

E.g

1. Remove 50% of your liquidity asymmetry (Equal rowan/chot); `sifnoded tx clp remove-liquidity --from sif --symbol chot --wBasis 5000 --asymmetry 0`

2. Remove 10% of your liquidity asymmetry; `sifnoded tx clp remove-liquidity --from akasha --symbol clink --wBasis 1000 --asymmetry 0`


#### Coming  

* Liquidity fees model  ... 
* Move minor api/ux enhancements ...le_previous_wrap)

#### Feature Requests / Bug reports

* https://github.com/Sifchain/sifnode/issues/new/choose


#### References

   * https://medium.com/thorchain/thorchains-liquidity-breakthrough-85a0fdbcd396
   * https://blog.cosmos.network/the-internet-of-blockchains-how-cosmos-does-interoperability-starting-with-the-ethereum-peg-zone-8744d4d2bc3f