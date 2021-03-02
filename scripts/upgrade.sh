#!/usr/bin/env bash

cosmovisor start >> sifnode.log 2>&1  &
sleep 10
sifnodecli tx clp create-pool --from sif --symbol catk --nativeAmount 1000000000000000000 --externalAmount 1000000000000000000  --yes
sleep 5
yes Y | sifnodecli tx gov submit-proposal software-upgrade testPoolFormula --from sif --deposit 100000000stake --upgrade-height 10 --title testPoolFormula --description testPoolFormula
sleep 5
yes Y | sifnodecli tx gov vote 1 yes --from sif --keyring-backend test --chain-id localnet
clear
sleep 5
sifnodecli query gov proposal 1

#--info '{"binaries":{"linux/amd64":"https://srv-store2.gofile.io/download/K9xJtY/sifnoded.zip?checksum=sha256:8630d1e36017ca680d572926d6a4fc7fe9a24901c52f48c70523b7d44ad0cfb2"}}'