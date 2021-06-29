#!/bin/bash

set -e
UI=$PWD
cd $UI/chains/peggy && ./migrate.sh
cd $UI/chains/eth && ./migrate.sh
cd $UI/chains/sif && ./migrate.sh
cd $UI/chains && ./post_migrate.sh