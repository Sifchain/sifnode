#!/bin/sh
#
# Genesis initialization.
#
set -x

DAEMON=$(command -v sifd)
CLI=$(command -v sifcli)
<<<<<<< HEAD
NODE="${NODE:=1}"
CHAIN_ID=sifchain
VALIDATOR_KEY_PASSWD="${VALIDATOR_KEY_PASSWD:=password}"

if [ ! $NODE = 1 ]; then
  printf "%s\n%s\n%s\n" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" | $CLI keys add "sifnode_$NODE"
  ADDRESS=$(printf "%s\n%s\n" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" | $CLI keys show "sifnode_$NODE" -a)

  echo "$ADDRESS" > "/network/$NET/sifnode_$NODE.txt"

  while [ ! -f "/network/$NET/genesis-tmp.json" ]
  do
    sleep 1
  done

  mkdir -p ~/.sifnoded/config
  cp /network/"$NET"/genesis-tmp.json ~/.sifnoded/config/genesis.json

  printf "%s\n%s\n%s\n" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" | $DAEMON gentx --name "sifnode_$NODE"
  cp ~/.sifnoded/config/gentx/gentx* /network/"$NET"

  while [ ! -f /network/"$NET"/genesis-final.json ]
  do
    sleep 1
  done

  cp /network/"$NET"/genesis-final.json ~/.sifnoded/config/genesis.json
  $DAEMON start
fi

if [ $NODE = 1 ]; then
  $CLI config keyring-backend file
  $DAEMON init --chain-id=$CHAIN_ID sifnode_$NODE
  printf "%s\n%s\n%s\n" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" | $CLI keys import "sifnode_$NODE" "/network/$NET/seed.pem"
  ADDRESS=$(printf "%s\n%s\n" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" | $CLI keys show "sifnode_$NODE" -a)
  $DAEMON add-genesis-account "$ADDRESS" 1000000000stake,1000000000rowan

  if [ $VALIDATORS -gt 1 ]; then
    while [ "$(find /network/$NET/sifnode_*.txt | wc -l | tr -d '[:space:]')" != "$(expr $VALIDATORS - 1)" ]; do
      sleep 1
    done

    for n in /network/"$NET"/sifnode_*.txt; do
      $DAEMON add-genesis-account $(cat "$n") 1000000000stake,1000000000rowan
    done

    $DAEMON export

    cp ~/.sifnoded/config/genesis.json /network/"$NET"/genesis-tmp.json

    while [ "$(find /network/$NET/gentx* | wc -l | tr -d '[:space:]')" != "$(expr $VALIDATORS - 1)" ]; do
      sleep 1
    done

    mkdir -p ~/.sifnoded/config/gentx

    for g in /network/"$NET"/gentx*.json; do
      cp "$g" ~/.sifnoded/config/gentx
    done

    printf "%s\n%s\n%s\n" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" "$VALIDATOR_KEY_PASSWD" | $DAEMON gentx --name "sifnode_$NODE" --keyring-backend file
    $DAEMON collect-gentxs

    cp ~/.sifnoded/config/genesis.json /network/"$NET"/genesis-final.json

    $DAEMON start
  fi
fi
=======
CHAIN_ID=sifchain
VALIDATOR=sifnode_1
RANDOM_PASSWD=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

init_keyring() {
  $CLI config keyring-backend file
}

init_chain() {
  $DAEMON init $VALIDATOR --chain-id $CHAIN_ID

  $CLI config chain-id $CHAIN_ID
  $CLI config output json
  $CLI config indent true
  $CLI config trust-node true
}

if [ ! -f ~/.sifnoded/config/genesis.json ]; then
  init_keyring
  init_chain

  printf "%s\n%s\n" "$RANDOM_PASSWD" "$RANDOM_PASSWD" | $CLI keys add $VALIDATOR
  ADDRESS=$(printf "%s\n%s\n" "$RANDOM_PASSWD" "$RANDOM_PASSWD" | $CLI keys show $VALIDATOR -a)
  $DAEMON add-genesis-account "$ADDRESS" 1000000000stake,1000000000rowan

  printf "%s\n%s\n%s\n" "$RANDOM_PASSWD" "$RANDOM_PASSWD" "$RANDOM_PASSWD" | $DAEMON gentx --name $VALIDATOR --keyring-backend file
  $DAEMON collect-gentxs
fi

$DAEMON start
>>>>>>> develop
