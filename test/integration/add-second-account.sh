#!/bin/bash
set -x
BASEDIR=$(pwd)
PASSWORD=$(yq r network-definition.yml "(*==$MONIKER).password")
ADDR=$(yq r network-definition.yml "(*==$MONIKER).address")

yes $PASSWORD | sifnodecli keys add user1

yes $PASSWORD | sifnodecli keys show user1 >> /network-definition.yml

