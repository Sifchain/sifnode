#!/bin/bash

# Takes a snapshot of the ganache data directory and prints
# the snapshot directory

. $(dirname $0)/vagrantenv.sh
GANACHE_DB_DIR=${1:-$GANACHE_DB_DIR}
shift

newganachedir=$(mktemp -d /tmp/ganachedbSnapshot.XXXX)
# Warning, this sudo will require a password on Mac...
#sudo rsync -a $GANACHE_DB_DIR/ $newganachedir/
rsync -a $GANACHE_DB_DIR/ $newganachedir/
echo $newganachedir
