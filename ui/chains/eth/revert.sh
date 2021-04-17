#!/bin/bash
. ./config.sh

# replace chain data with archive if exists
if [[ -f "$snapshot_loc" ]]; then
  rm -rf $db_loc
  mkdir -p $db_loc
  tar -zxf $snapshot_loc --directory $db_loc
fi

# restart sifnode
./start.sh