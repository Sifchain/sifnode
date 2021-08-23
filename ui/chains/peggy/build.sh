#!/bin/bash

. ../credentials.sh

cd $PWD/../../../smart-contracts

yarn && yarn add truffle && yarn truffle compile