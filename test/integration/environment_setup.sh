required="${BASEDIR:?"Must be set to where you check out sifnode"}"

cd $BASEDIR/test/integration

required="${DEPLOYMENT_NAME:?"Must be set to a deployment name like sandpit"}"
required="${ROWAN_SOURCE:?"Must be set to a sif address that contains rowan, and the key must be in your keyring"}"
required="${ETHEREUM_NETWORK:?"Must be set to an etherereum endpoint.  Can be ropsten or a url."}"
required="${ETHEREUM_PRIVATE_KEY:?"Must be set to the private key of the address specified in ETHEREUM_ADDRESS or OPERATOR_ADDRESS"}"
required="${SIFNODE:?"Must be set to sifnode endpoint."}"
required="${INFURA_PROJECT_ID:?"Must be set."}"

if [ -z "$ETHEREUM_ADDRESS$OPERATOR_ADDRESS" ]; then echo must set one of ETHEREUM_ADDRESS or OPERATOR_ADDRESS; exit 1; fi

export SMART_CONTRACTS_DIR=$BASEDIR/smart-contracts
export SOLIDITY_JSON_PATH=$BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME
export SMART_CONTRACT_ARTIFACT_DIR=$SOLIDITY_JSON_PATH

export BRIDGE_REGISTRY_ADDRESS=$(cat $SOLIDITY_JSON_PATH/BridgeRegistry.json | jq -r ".networks[\"$ETHEREUM_NETWORK_ID\"].address")
export BRIDGE_TOKEN_ADDRESS=$(cat $SOLIDITY_JSON_PATH/BridgeToken.json | jq -r ".networks[\"$ETHEREUM_NETWORK_ID\"].address")
export BRIDGE_BANK_ADDRESS=$(cat $SOLIDITY_JSON_PATH/BridgeBank.json | jq -r ".networks[\"$ETHEREUM_NETWORK_ID\"].address")

cp $BASEDIR/smart-contracts/build/contracts/SifchainTestToken.json $SOLIDITY_JSON_PATH

echo ========== Sample commands ==========

echo; echo == erowan balance
echo yarn -s --cwd $BASEDIR/smart-contracts integrationtest:getTokenBalance \
  --symbol $BRIDGE_TOKEN_ADDRESS \
  --ethereum_private_key_env_var "ETHEREUM_PRIVATE_KEY" \
  --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME \
  --gas estimate \
  --ethereum_network $ETHEREUM_NETWORK \
  --ethereum_address $ETHEREUM_ADDRESS \

echo; echo == eth balance
echo yarn -s --cwd $BASEDIR/smart-contracts integrationtest:getTokenBalance \
  --symbol eth \
  --ethereum_private_key_env_var "ETHEREUM_PRIVATE_KEY" \
  --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME \
  --gas estimate \
  --ethereum_network $ETHEREUM_NETWORK \
  --ethereum_address $ETHEREUM_ADDRESS \

echo; echo == mint erowan
echo yarn -s --cwd /home/james/workspace/sifnode/smart-contracts integrationtest:mintTestnetTokens  \
  --symbol $BRIDGE_TOKEN_ADDRESS \
  --ethereum_private_key_env_var "OPERATOR_PRIVATE_KEY" \
  --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME \
  --gas estimate \
  --ethereum_network $ETHEREUM_NETWORK \
  --ethereum_address $ETHEREUM_ADDRESS \
  --operator_address $OPERATOR_ADDRESS \
  --amount 100000000000000000000000000

echo; echo == lock eth
echo yarn -s --cwd $BASEDIR/smart-contracts integrationtest:sendLockTx --sifchain_address $ROWAN_SOURCE \
  --symbol eth \
  --ethereum_private_key_env_var "ETHEREUM_PRIVATE_KEY" \
  --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME \
  --gas estimate \
  --ethereum_network $ETHEREUM_NETWORK \
  --bridgebank_address $BRIDGE_BANK_ADDRESS \
  --ethereum_address $ETHEREUM_ADDRESS \
  --amount 1700000000000000000

echo; echo == burn erowan
echo yarn -s --cwd $BASEDIR/smart-contracts integrationtest:sendBurnTx \
  --symbol $BRIDGE_TOKEN_ADDRESS \
  --ethereum_private_key_env_var "ETHEREUM_PRIVATE_KEY" \
  --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME \
  --gas estimate \
  --ethereum_network $ETHEREUM_NETWORK \
  --bridgebank_address $BRIDGE_BANK_ADDRESS \
  --ethereum_address $ETHEREUM_ADDRESS \
  --sifchain_address $ROWAN_SOURCE \
  --amount 17

echo; echo == burn erowan from operator account
echo yarn -s --cwd /home/james/workspace/sifnode/smart-contracts integrationtest:sendBurnTx \
  --symbol $BRIDGE_TOKEN_ADDRESS \
  --ethereum_private_key_env_var "OPERATOR_PRIVATE_KEY" \
  --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME \
  --gas estimate \
  --ethereum_network $ETHEREUM_NETWORK \
  --bridgebank_address $BRIDGE_BANK_ADDRESS \
  --ethereum_address $OPERATOR_ADDRESS \
  --sifchain_address $ROWAN_SOURCE \
  --amount 100000000000000000000000000

echo; echo == whitelisted tokens
echo yarn -s --cwd $BASEDIR/smart-contracts \
  integrationtest:whitelistedTokens \
  --bridgebank_address $BRIDGE_BANK_ADDRESS \
  --json_path $BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME \
  --ethereum_network $ETHEREUM_NETWORK \

sifnodecmd=sifnoded

echo; echo == sifchain balance
echo $sifnodecmd q auth account --node $SIFNODE $ROWAN_SOURCE

echo; echo == sifchain transaction
echo $sifnodecmd q tx --node $SIFNODE --chain-id $DEPLOYMENT_NAME 193EFB4A5D20BEC58ADE8BACEB38264870ADD8BAFEA9D6DAABE554B0ACBC0C93

echo; echo == all account balances
echo "$sifnodecmd keys list --keyring-backend test -o json | jq -r '.[].address' | parallel $sifnodecmd q auth account --node $SIFNODE -o json {} | grep coins"

echo; echo == burn ceth
echo $sifnodecmd tx ethbridge burn \
  $ROWAN_SOURCE $ETHEREUM_ADDRESS 100 ceth 58560000000000000 \
  --node $SIFNODE \
  --keyring-backend test \
  --fees 100000rowan \
  --ethereum-chain-id=$ETHEREUM_NETWORK_ID \
  --chain-id=$DEPLOYMENT_NAME  \
  --yes \
  --from $ROWAN_SOURCE \

echo; echo == send ceth
echo $sifnodecmd tx send $ROWAN_SOURCE sifsomedestination 100rowan \
  --node $SIFNODE \
  --keyring-backend test \
  --fees 100000rowan \
  --chain-id=$DEPLOYMENT_NAME  \
  --yes \

echo; echo == Simple test run against $DEPLOYMENT_NAME:
echo python3 -m pytest --color=yes -x -olog_cli=true -olog_level=DEBUG -v -olog_file=vagrant/data/pytest.log -v src/py/test_eth_transfers.py

echo; echo == Load test run against $DEPLOYMENT_NAME - change NTRANSFERS to a large number:
echo NTRANSFERS=2 python3 -m pytest -olog_level=DEBUG -olog_file=vagrant/data/pytest.log -v src/py/test_bulk_transfers_to_ethereum.py::test_bulk_transfers_from_sifchain

echo; echo
