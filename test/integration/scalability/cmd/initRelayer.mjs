import { initRelayer } from "../lib/initRelayer.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--chain": String,
      "--network": String,
      "--registryFrom": Number,
      "--rpcPortA": Number,
      "--p2pPortA": Number,
      "--pprofPortA": Number,
      "--homeA": String,
      "--rpcPortB": Number,
      "--p2pPortB": Number,
      "--pprofPortB": Number,
      "--homeA": String,
    },
    `
Usage:

  yarn initRelayer [options]

Initiate a new relayer locally to pair a foreign IBC chain to sifchain.

Options:

--chain         Select a predifined chain in chains.json
--network       Select a predifined network in chains.json
--registryFrom  Directory for storing global relayer registry data
--rpcPortA      RPC port number
--p2pPortA      P2P port number
--pprofPortA    pprof port number
--homeA         Directory for config and data
--rpcPortB      RPC port number
--p2pPortB      P2P port number
--pprofPortB    pprof port number
--homeB         Directory for config and data
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const registryFrom = args["--registryFrom"] || undefined;
  const rpcPortA = args["--rpcPortA"] || undefined;
  const p2pPortA = args["--p2pPortA"] || undefined;
  const pprofPortA = args["--pprofPortA"] || undefined;
  const homeA = args["--homeA"] || undefined;
  const rpcPortB = args["--rpcPortB"] || undefined;
  const p2pPortB = args["--p2pPortB"] || undefined;
  const pprofPortB = args["--pprofPortB"] || undefined;
  const homeB = args["--homeB"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
  });
  await initRelayer({
    chainProps,
    registryFrom,
    rpcPortA,
    p2pPortA,
    pprofPortA,
    homeA,
    rpcPortB,
    p2pPortB,
    pprofPortB,
    homeB,
  });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
