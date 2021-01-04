#!/usr/bin/env bash

. ../credentials.sh


parallelizr() {
  for cmd in "$@"; do {
    echo "Process \"$cmd\" started";
    $cmd & pid=$!
    PID_LIST+=" $pid";
  } done

  trap "kill $PID_LIST" SIGINT

  echo "Parallel processes have started";

  wait $PID_LIST

  echo "All processes have completed";
}

rm -rf ~/.sifnoded
rm -rf ~/.sifnodecli

sifnoded init test --chain-id=sifchain

sifnodecli config output json
sifnodecli config indent true
sifnodecli config trust-node true
sifnodecli config chain-id sifchain
sifnodecli config keyring-backend test

echo "Generating deterministic account - ${SHADOWFIEND_NAME}"
echo "${SHADOWFIEND_MNEMONIC}" | sifnodecli keys add ${SHADOWFIEND_NAME} --recover

echo "Generating deterministic account - ${AKASHA_NAME}"
echo "${AKASHA_MNEMONIC}" | sifnodecli keys add ${AKASHA_NAME} --recover

sifnoded add-genesis-account $(sifnodecli keys show ${SHADOWFIEND_NAME} -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000stake
sifnoded add-genesis-account $(sifnodecli keys show ${AKASHA_NAME} -a) 1000000000rowan,1000000000catk,1000000000cbtk,1000000000ceth,100000000stake

sifnoded add-genesis-validators $(sifnodecli keys show ${SHADOWFIEND_NAME} -a --bech val)

sifnoded gentx --name ${SHADOWFIEND_NAME} --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis

echo "Starting test chain"

parallelizr "sifnoded start" "sifnodecli rest-server  --unsafe-cors --trace"


#sifnoded start --log_level="main:info,state:error,statesync:info,*:error"
