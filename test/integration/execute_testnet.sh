required="${BASEDIR:?"Must be set to where you check out sifnode"}"

cd $BASEDIR/test/integration

required="${DEPLOYMENT_NAME:?"Must be set to a deployment name like sandpit"}"
required="${ROWAN_SOURCE:?"Must be set to a sif address that contains rowan, and the key must be in your keyring"}"
required="${ROWAN_SOURCE_KEY:?"Must be set to the name of a key in your keyring that is the same as ROWAN_SOURCE"}"
required="${ETHEREUM_NETWORK:?"Must be set to an etherereum endpoint.  Can be ropsten or a url."}"
required="${ETHEREUM_PRIVATE_KEY:?"Must be set to the private key of the address specified in ETHEREUM_ADDRESS or OPERATOR_ADDRESS"}"
required="${SIFNODE:?"Must be set to sifnode endpoint."}"
required="${INFURA_PROJECT_ID:?"Must be set."}"

if [ -z "$ETHEREUM_ADDRESS$OPERATOR_ADDRESS" ]; then echo must set one of ETHEREUM_ADDRESS or OPERATOR_ADDRESS; exit 1; fi

export SMART_CONTRACTS_DIR=$BASEDIR/smart-contracts
export SOLIDITY_JSON_PATH=$BASEDIR/smart-contracts/deployments/$DEPLOYMENT_NAME
export SMART_CONTRACT_ARTIFACT_DIR=$SOLIDITY_JSON_PATH

echo you are now set up to run tests like this against $DEPLOYMENT_NAME:
echo python3 -m pytest --color=yes -x -olog_cli=true -olog_level=DEBUG -v -olog_file=vagrant/data/pytest.log -v src/py/test_eth_transfers.py
