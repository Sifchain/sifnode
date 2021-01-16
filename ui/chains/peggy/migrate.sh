#!/bin/bash

. ../credentials.sh

cd ../../../smart-contracts

yarn && yarn test:setup && yarn migrate 

