#!/bin/bash

sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
  '{"pool":{"external_asset": "ceth"}}' \
  --chain-id localnet
