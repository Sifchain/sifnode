#!/bin/bash

$1 keys show $2 --keyring-backend test -a 2> /dev/null || echo $2