#!/bin/bash

#
# Sifnode entrypoint.
#

set -x

. /sifnode/test/integration/vagrantenv.sh

ADD_VALIDATOR_TO_WHITELIST=${1:-${ADD_VALIDATOR_TO_WHITELIST}}
shift

if [ -z "${ADD_VALIDATOR_TO_WHITELIST}" ]
then
  # no whitelist validator requested; mostly useful for testing validator whitelisting
  echo $0: no whitelisted validators
else
  whitelisted_validator=$(yes $VALIDATOR1_PASSWORD | sifnodecli keys show -a --bech val $MONIKER --home $CHAINDIR/.sifnodecli)
  echo $0: whitelisted validator $whitelisted_validator
  sifnoded add-genesis-validators $whitelisted_validator --home $CHAINDIR/.sifnoded
fi

sifnoded start --minimum-gas-prices 0.5rowan --rpc.laddr tcp://0.0.0.0:26657 --home $CHAINDIR/.sifnoded
