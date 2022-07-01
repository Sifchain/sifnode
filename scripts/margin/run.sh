#!/usr/bin/env bash

export FEATURE_TOGGLE_SDK_045=1
export FEATURE_TOGGLE_MARGIN_CLI_ALPHA=1
export GOFLAGS="-modfile=go_FEATURE_TOGGLE_SDK_045.mod"
export GOTAGS="FEATURE_TOGGLE_SDK_045,FEATURE_TOGGLE_MARGIN_CLI_ALPHA"

killall sifnoded

cd ../..
make install
sifnoded start --trace
