#!/usr/bin/env bash

cosmovisor start >> sifnode.log 2>&1  &
sleep 10
yes Y | sifnodecli tx gov submit-proposal software-upgrade testupgrade --from shadowfiend --deposit 100000000stake --upgrade-height 5 --info '{"binaries":{"linux/amd64":"https://srv-store2.gofile.io/download/K9xJtY/sifnoded.zip?checksum=sha256:8630d1e36017ca680d572926d6a4fc7fe9a24901c52f48c70523b7d44ad0cfb2"}}' --title testupgrade --description testupgrade
sleep 5
yes Y | sifnodecli tx gov vote 1 yes --from shadowfiend --keyring-backend test --chain-id sifchain
clear
sleep 5
sifnodecli query gov proposal 1

