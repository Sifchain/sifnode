#!/bin/bash
. ./config.sh

# replace chain data with archive if exists
echo "extracting..." 
if [[ -f "$snapshot_loc" ]]; then
  rm -rf $db_loc
  mkdir -p $db_loc
  tar -zxf $snapshot_loc --directory $db_loc
  echo "extracted '$snapshot_loc' to '$db_loc'"
fi

# peggy is special as we need the contents of the smart-contracts build folder
if [[ -f "$snapshot2_loc" ]]; then
  rm -rf $db2_loc
  mkdir -p $db2_loc
  tar -zxf $snapshot2_loc --directory $db2_loc
  echo "extracted '$snapshot2_loc' to '$db2_loc'"
fi

# restart sifnode
./start.sh