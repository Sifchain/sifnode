# How to: Add IBC ERC20 tokens

## Overview

There are two distinct steps to this script.

One is creating and deploying the new bridge tokens.

Two is registering the bridge tokens in the bridgebank.

Usually, each script will be run by a different person. The script to register bridge tokens will be run by a user with priviledged access to the bridgebank with the owner role.

If you are deploying new tokens, please consult the runbook 'DeployIbcTokens.md'.

If you are registering tokens, please consult the runbook 'RegisterIbcTokens.md'.

If you are an engineer and want to TEST the scripts, please keep reading this doc.

## Testing with forked mainnnet:

Since you're running two scripts, you'll need a hardhat node running (otherwise the first script will run, execute transactions, then throw them away).

Start a hardhat node in a shell:

    npx hardhat node

Run the two scripts sequentially in a different shell:

    deployIbcTokens:test
    registerIbcTokens:test

## Update symbol_translator.json

Modify https://github.com/Sifchain/chainOps/blob/main/.github/workflows/variableMapping/ebrelayer.yaml
with the new symbol_translation.json entries. It's a yaml file, so you need to escape json - you could
use something like `jq -aRs` to do the quoting.

To push that data out and restart the relayers, use https://github.com/Sifchain/chainOps/actions/workflows/peggy-ebrelayer-deployment.yml
