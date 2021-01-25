#!/bin/bash

. $(dirname $0)/shell_utilities.sh

json_contracts_dir=${1}
shift
envexportfile=${1}
shift

set_persistant_env_var BRIDGE_REGISTRY_ADDRESS $(cat $json_contracts_dir/BridgeRegistry.json | jq -r '.networks["5777"].address') $envexportfile required
set_persistant_env_var BRIDGE_TOKEN_ADDRESS $(cat $json_contracts_dir/BridgeToken.json | jq -r '.networks["5777"].address') $envexportfile required
set_persistant_env_var BRIDGE_BANK_ADDRESS $(cat $json_contracts_dir/BridgeBank.json | jq -r '.networks["5777"].address') $envexportfile required

cat $envexportfile