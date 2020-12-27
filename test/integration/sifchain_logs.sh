#!/bin/bash

set -e

basedir=$(dirname $0)
. $basedir/vagrantenv.sh

output=$datadir/sifchaintxs
rm -f $output*

hashes=$(cat $EBRELAYER_LOG | grep "^txhash: " | sed -e "s/txhash: //")
for i in $hashes
do
  sifnodecli q tx $i >> $output.txt
  sifnodecli q tx $i -o json | jq -c >> $output.json
done
