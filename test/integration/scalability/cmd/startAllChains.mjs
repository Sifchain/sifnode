#!/usr/bin/env zx

import { startAllChains } from "../lib/startAllChains.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

const args = arg(
  {
    "--chain": String,
    "--network": String,
    "--binary": String,
    "--rpcPort": Number,
    "--p2pPort": Number,
    "--pprofPort": Number,
    "--home": String,
  },
  `
Usage:

  yarn startAllChains [options]

Start all the chains locally.

Options:

--network   Select a predifined network in chains.json
--rpcPort   Initial RPC port number
--p2pPort   Initial P2P port number
--pprofPort Initial pprof port number
--home      Global directory for config and data of initiated chains
`
);

const network = args["--network"] || undefined;
const rpcPort = args["--rpcPort"] || undefined;
const p2pPort = args["--p2pPort"] || undefined;
const pprofPort = args["--pprofPort"] || undefined;
const home = args["--home"] || undefined;

const chainProps = getChainProps({
  network,
  rpcPort,
  p2pPort,
  pprofPort,
  home,
});
await startAllChains({
  ...chainProps,
});
