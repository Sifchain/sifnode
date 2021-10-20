#!/usr/bin/env zx

import { initRelayer } from "../lib/initRelayer.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

const args = arg(
  {
    "--chain": String,
    "--network": String,
    "--rpcPort": Number,
    "--p2pPort": Number,
    "--pprofPort": Number,
    "--registryFrom": Number,
    "--home": String,
  },
  `
Usage:

  yarn initRelayer [options]

Initiate a new relayer locally to pair a foreign IBC chain to sifchain.

Options:

--chain         Select a predifined chain in chains.json
--network       Select a predifined network in chains.json
--rpcPort       RPC port number
--p2pPort       P2P port number
--pprofPort     pprof port number
--registryFrom  Directory for storing global relayer registry data
--home          Directory for config and data
`
);

const chain = args["--chain"] || undefined;
const network = args["--network"] || undefined;
const rpcPort = args["--rpcPort"] || undefined;
const p2pPort = args["--p2pPort"] || undefined;
const pprofPort = args["--pprofPort"] || undefined;
const registryFrom = args["--registryFrom"] || undefined;
const home = args["--home"] || undefined;

const chainProps = getChainProps({
  chain,
  network,
  rpcPort,
  p2pPort,
  pprofPort,
  home,
});
await initRelayer({
  chainProps,
  registryFrom,
});
