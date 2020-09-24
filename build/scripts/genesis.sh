#!/bin/sh
#
# Genesis initialization.
#
set -x

DAEMON=$(command -v sifd)
CLI=$(command -v sifcli)
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
