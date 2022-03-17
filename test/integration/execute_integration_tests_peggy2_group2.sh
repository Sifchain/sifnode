#!/bin/bash

function _autodetect_project_root {
    PRG="$0"
    while [ -h "$PRG" ]; do
        ls=$(ls -ld "$PRG")
        link=$(expr "$ls" : '.*-> \(.*\)$')
        if expr "$link" : '/.*' > /dev/null; then
            PRG="$link"
        else
            PRG="$(dirname "$PRG")/$link"
        fi
    done
    realpath "$(dirname "$PRG")/../.."
}

TEST_INTEGRATION_PY_DIR="$(_autodetect_project_root)/test/integration/src/peggy2"

python3 -m pytest -olog_level=DEBUG -v -olog_file=/tmp/pytest-peggy2.log \
  ${TEST_INTEGRATION_PY_DIR}/test_sifnode_transfers.py \
