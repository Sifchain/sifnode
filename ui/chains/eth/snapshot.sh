#!/bin/bash
. ./config.sh

mkdir -p ../snapshots
mkdir -p $db_loc

# archive data folder
here=$(pwd) 
cd $db_loc && tar -zcvf $here/$snapshot_loc .  
cd $here