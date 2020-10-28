#!/usr/bin/env bash
mac install ethereum
cd testnet-contracts
yarn
ebrelayer generate
cp .env.example .env