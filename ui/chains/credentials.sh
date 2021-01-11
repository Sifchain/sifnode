# Assumed to run from the ui folder
BASE_DIR=..

export SHADOWFIEND_NAME=shadowfiend
export SHADOWFIEND_MNEMONIC="race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

export AKASHA_NAME=akasha
export AKASHA_MNEMONIC="hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard"

export ETHEREUM_ROOT_MNEMONIC="candy maple cake sugar pudding cream honey rich smooth crumble sweet treat"

# Required to run ebrelayer
export BRIDGE_TOKEN_ADDRESS=$(cat $PWD/../../../smart-contracts/build/contracts/BridgeToken.json | jq -r '.networks["5777"].address')
export BRIDGE_REGISTRY_ADDRESS=$(cat $PWD/../../../smart-contracts/build/contracts/BridgeRegistry.json | jq -r '.networks["5777"].address') 