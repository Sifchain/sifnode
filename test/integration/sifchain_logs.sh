#!/bin/bash

output=$(dirname $0)/vagrant/data/sifchaintxs.txt
rm -f $output

for i in $(docker logs integration_sifnode1_1 2>&1 | grep "^txhash: " | sed -e "s/txhash: //")
do
  docker exec -ti integration_sifnode1_1 bash -c "sifnodecli q tx $i" >> $output
done
