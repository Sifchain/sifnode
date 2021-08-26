#!/bin/bash

# this script can be run only against a test chain.  It relies on tight control over
# ganache and knows where ganache stores its data.

set -e

. $(dirname $0)/vagrantenv.sh
. $(dirname $0)/shell_utilities.sh

loglevel=${LOG_LEVEL:-DEBUG}

logecho $0 starting

# To run these tests you must first create a snapshot with "make.py create_snapshot s1"
python3 -m pytest -olog_level=$loglevel -v -olog_file=/tmp/log.txt -v \
    ${TEST_INTEGRATION_PY_DIR}/test_random_currency_roundtrip_with_snapshots.py
