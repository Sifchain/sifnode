#!/bin/bash
set -x
BASEDIR=$(pwd)
PASSWORD=$(yq r network-definition.yml "(*==$MONIKER).password")
ADDR=$(yq r network-definition.yml "(*==$MONIKER).address")

yes $PASSWORD | sifnodecli keys add user1
#yes $PASSWORD | sifnodecli keys add user2

# change these to tx's that move funds out of validator account
yes $PASSWORD | sifnoded add-genesis-account $(sifnodecli keys show user1 -a) 1000rwn,100000000stake

yes $PASSWORD | sifnoded gentx --name user1 --keyring-backend file

yes $PASSWORD | sifnoded --home="/root/.sifnoded" collect-gentxs
sifnoded validate-genesis

#yes $PASSWORD | sifnoded add-genesis-account $(sifnodecli keys show user2 -a) 1000rwn,100000000stake

#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user1 -a) "1000rown, 2clink, 1000chot,1000cusdt,1000cusda" -y

#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user2 -a) "1000rown, 1000clink, 1000chot,1000cusdt,1000cusda" -y
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user1 -a) "1000rown"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user1 -a) "1000clink"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user1 -a) "1000chot"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user1 -a) "1000cusdt"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user1 -a) "1000cusda"
#
#
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user2 -a) "1000rown"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user2 -a) "1000clink"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user2 -a) "1000chot"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user2 -a) "1000cusdt"
#yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show user2 -a) "1000cusda"
