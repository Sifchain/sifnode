#!/bin/bash
. ./config.sh

mkdir -p ../snapshots
mkdir -p $db_loc
mkdir -p $db2_loc

# archive data folder
here=$(pwd) 
cd $db_loc && tar -zcvf $here/$snapshot_loc .  

# peggy is special in that it relies on data within the build folder of smart-contracts
# Lets cache the build folder
cd $here/$db2_loc && tar -zcvf $here/$snapshot2_loc .  

cd $here
