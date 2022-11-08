import { startChain } from "../lib/startChain.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--chain": String,
      "--network": String,
      "--binary": String,
      "--rpcPort": Number,
      "--p2pPort": Number,
      "--pprofPort": Number,
      "--home": String,
      "--binPath": String,
    },
    `
Usage:

  yarn startChain [options]

Start a new chain locally.

Options:

--chain     Select a predifined chain in chains.json
--network   Select a predifined network in chains.json
--binary    Binary name of the chain
--rpcPort   RPC port number
--p2pPort   P2P port number
--pprofPort pprof port number
--home      Directory for config and data
--binPath   Directory for binaries location
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const binary = args["--binary"] || undefined;
  const rpcPort = args["--rpcPort"] || undefined;
  const p2pPort = args["--p2pPort"] || undefined;
  const pprofPort = args["--pprofPort"] || undefined;
  const home = args["--home"] || undefined;
  const binPath = args["--binPath"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
    binary,
    rpcPort,
    p2pPort,
    pprofPort,
    home,
    binPath,
  });
  await startChain({
    ...chainProps,
  });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
