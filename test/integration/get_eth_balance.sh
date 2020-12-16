#!/bin/bash 

addr=$1
shift
token=$1
shift

docker exec -ti ${CONTAINER_NAME} bash -c "cd /smart-contracts; yarn peggy:getTokenBalance ${addr:=${BRIDGE_BANK_ADDRESS}} ${token:=eth}"