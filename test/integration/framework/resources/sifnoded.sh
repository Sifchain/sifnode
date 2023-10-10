#!/bin/bash

set -euo pipefail

BASEDIR="$0"
while [ -h "$BASEDIR" ]; do
    ls=$(ls -ld "$BASEDIR")
    link=$(expr "$ls" : '.*-> \(.*\)$')
    if expr "$link" : '/.*' > /dev/null; then
        BASEDIR="$link"
    else
        BASEDIR="$(dirname "$BASEDIR")/$link"
    fi
done
BASEDIR="$(dirname "$BASEDIR")"
BASEDIR="$(cd "$BASEDIR"; pwd)"

SERVICE_NAME="${1-}"

if [ -z "$SERVICE_NAME" ]; then
    echo "Mising parameter: SERVICE_NAME"
    exit 255
fi

echo $BASEDIR
echo $SERVICE_NAME

export DAEMON_NAME="sifnoded"
export DAEMON_HOME="$BASEDIR/$SERVICE_NAME"
export DAEMON_ALLOW_DOWNLOAD_BINARIES="true"
export DAEMON_RESTART_AFTER_UPGRADE="true"
export UNSAFE_SKIP_BACKUP="true"

echo "Running cosmovisor for ${SERVICE_NAME} in ${DAEMON_HOME}..."

"$BASEDIR/cosmovisor" start --home "${DAEMON_HOME}"
