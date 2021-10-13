#!/bin/bash -x

sifnoded keys show $1 --keyring-backend test -a 2> /dev/null || echo $1