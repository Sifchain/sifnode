#!/bin/bash

sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
  '{"remove_liquidity":{"w_basis_points": "5000", "asymmetry": "0"}}' \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y