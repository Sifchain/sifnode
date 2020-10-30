#!/usr/bin/env bash
brew install ethereum
cd testnet-contracts
yarn
# ebrelayer generate
cp .env.example .env