# README

Scripts for local development and testing.

## How to set up the e2e test environment in local host
1. run init.sh
```
./scripts/init.sh
```
2. run setup-e2e.sh
```
./scripts/setup-e2e.sh
```
3. start the sifnoded, truffle ethereum and relayer.
the instruction can be found in peggy-e2e-test.md

## Test environment for testnet
For Testnet, the sifnoded, relayer should be automatically deployed by devops tool.
Ralayer connected to the Ethereum network (for exmaple ropsten) endpoint.

## Use the python scripts to test
1. run basic test for sifchain, just need sifnode running
```
python3 peggy-basic-test.py
```
2. run e2e test, need both sifnode and relayer running. also the all js scripts in testnet-contracts works well.
```
python3 peggy-e2e-test.py
```
