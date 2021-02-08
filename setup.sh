#!/usr/bin/env bash
./test/integration/start-integration-env.sh
source test/integration/vagrantenv.sh
./test/integration/ganache_start.sh 10
