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

echo "starting..." 
# restart sifnode
./start.sh