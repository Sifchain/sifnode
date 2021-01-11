#!/bin/bash

set -xe

basedir=$(dirname $0)
. $basedir/vagrantenv.sh

if [ -z $EBRELAYER_LOG ]
then
  echo $0: must specify EBRELAYER_LOG
  exit 1
fi

output=$datadir/sifchaintxs
rm -f $output*

hashes=$(cat $EBRELAYER_LOG | grep "^txhash: " | sed -e "s/txhash: //")
for i in $hashes
do
  sifnodecli q tx $i >> $output.txt
  sifnodecli q tx $i -o json | jq -c . >> $output.json
done
