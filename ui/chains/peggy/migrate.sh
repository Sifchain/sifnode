#!/bin/bash

. ../credentials.sh

cd ../../../smart-contracts

cp .env.ui.example .env

yarn && yarn migrate 

