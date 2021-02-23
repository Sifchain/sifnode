#!/bin/bash 

# Sets up a bare Ubuntu environment with all the tools we use
# for integration tests

set -e

scriptdir=$(dirname $0)

sudo bash $scriptdir/setup-linux-environment-root.sh $(id -u -n)
bash $scriptdir/setup-linux-environment-user.sh
