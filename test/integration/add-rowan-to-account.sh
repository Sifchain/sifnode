#!/bin/bash

# sends $1 rowan to account $2

amount=$1
destination=$2
set -x

BASEDIR=$(pwd)
PASSWORD=$(yq r network-definition.yml "(*==$MONIKER).password")
ADDR=$(yq r network-definition.yml "(*==$MONIKER).address")

yes $PASSWORD | sifnodecli tx send $ADDR $(yes $PASSWORD | sifnodecli keys show $destination -a) "${amount}rowan" -y

