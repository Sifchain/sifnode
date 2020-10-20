#!/usr/bin/env bash
cd testnet-contracts
yarn
ebrelayer generate
cp .env.example .env