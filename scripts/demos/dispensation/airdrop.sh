#!/usr/bin/env bash

sifnodecli tx dispensation airdrop $(sifnodecli keys show mkey -a) input.json output.json --generate-only >> offlinetx.json
sifnodecli tx sign --multisig $(sifnodecli keys show mkey -a) --from $(sifnodecli keys show sif -a)  offlinetx.json >> sig1.json
sifnodecli tx sign --multisig $(sifnodecli keys show mkey -a) --from $(sifnodecli keys show akasha -a)  offlinetx.json >> sig2.json
sifnodecli tx multisign offlinetx.json mkey sig1.json sig2.json >> signedtx.json
sifnodecli tx broadcast signedtx.json
rm -rf offlinetx.json sig1.json sig2.json signedtx.json

#rm -rf ../../../offlinetx.json ../../../sig1.json ../../../sig2.json ../../../signedtx.json

