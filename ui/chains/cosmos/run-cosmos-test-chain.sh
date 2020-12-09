#!/usr/bin/env bash

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

echo "Generating deterministic account - shadowfiend"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | sifnodecli keys add shadowfiend --recover

echo "Generating deterministic account - akasha"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | sifnodecli keys add akasha --recover

sifnoded add-genesis-account $(sifnodecli keys show shadowfiend -a) 1000000000rwn,1000000000catk,1000000000cbtk,1000000000ceth,100000000stake
sifnoded add-genesis-account $(sifnodecli keys show akasha -a) 1000000000rwn,1000000000catk,1000000000cbtk,1000000000ceth,100000000stake

sifnoded gentx --name shadowfiend --keyring-backend test

echo "Collecting genesis txs..."
sifnoded collect-gentxs

echo "Validating genesis file..."
sifnoded validate-genesis

echo "Starting test chain"

parallelizr "sifnoded start" "sifnodecli rest-server  --unsafe-cors --trace"


#sifnoded start --log_level="main:info,state:error,statesync:info,*:error"

