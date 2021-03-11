# only once, add a sifchain account that has rowan
# sifnodecli keys import rowansource ~/sandpitkey.txt --keyring-backend test
# move rowan to erowan in the ETHEREUM_ADDRESS
# sifnodecli tx ethbridge lock rowansource $ETHEREUM_ADDRESS 1000000000000000000000 rowan 18332015000000000 --keyring-backend test --fees 100000rowan --ethereum-chain-id=5777 --chain-id=ropsten --home /home/james/.sifnodecli --from rowansource --yes
# Set an ethereum private key for an acount that has eth on the testnet
export SMART_CONTRACTS_DIR=/Users/kevindegraaf/sifnode/smart-contracts
export SOLIDITY_JSON_PATH=/Users/kevindegraaf/sifnode/smart-contracts/deployments/sandpit
export ETHEREUM_ADDRESS=0x1e0220B251eE648C7F3B6Fc31E6d309141f2e464
export OPERATOR_ADDRESS=0x1e0220B251eE648C7F3B6Fc31E6d309141f2e464
export OPERATOR_ACCOUNT=0x1e0220B251eE648C7F3B6Fc31E6d309141f2e464
#export ROWAN_SOURCE=sif1pvnu2kh826vn8r0ttlgt82hsmfknvcnf7qmpvk
export ROWAN_SOURCE=sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5
export ROWAN_SOURCE_KEY=kevin_test
#export ETHEREUM_NETWORK=http://localhost:8545
export ETHEREUM_NETWORK=ropsten
export SIFNODE=http://54.218.170.168:26657
#export INFURA_PROJECT_ID=c413023ff7944d21b694664b31a52faf
export INFURA_PROJECT_ID=fafcdfa80cc14ba2a7fb414ba86e8b24
#cd ~/sifnode/test/integration
# source vagrantenv.sh
export CHAINNET="sandpit"
source <(./smart_contract_env.sh ~/sifnode/smart-contracts/deployments/sandpit)
#export ETHEREUM_PRIVATE_KEY=e58d27b35688a5764de589eb9624e02c0bdb9be633db27cc746d3b5e3773e33d  # test account, not owner
#export ETHEREUM_PRIVATE_KEY=c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3
export ETHEREUM_PRIVATE_KEY=30dd94b42b731aa5fe738353d897fb938cf7e2a1dbce629dead1b3294ede4f3c
#export ETHEREUM_PRIVATE_KEY=54d0655c3cdf5b1f5c46468829ad1d7ed95ace59cafe1613187249b3b20d5b65
echo python3 -m pytest --color=yes -x -olog_cli=true -olog_level=DEBUG -v -olog_file=sandpit.log -v src/py/test_liquidity_pools.py::test_add_faucet_coins