import { startAllRelayers } from "../lib/startAllRelayers.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--network": String,
      "--rpcInitialPort": Number,
      "--p2pInitialPort": Number,
      "--pprofInitialPort": Number,
      "--configPath": String,
    },
    `
Usage:

  yarn startAllRelayers [options]

Start all the IBC realyers locally.

Options:

--network               Select a predifined network in chains.json
--rpcInitialPort        Initial RPC port number
--p2pInitialPort        Initial P2P port number
--pprofInitialPort      Initial pprof port number
--configPath                  Global directory for config and data of initiated chains
`
  );

  const network = args["--network"] || undefined;
  const rpcInitialPort = args["--rpcInitialPort"] || undefined;
  const p2pInitialPort = args["--p2pInitialPort"] || undefined;
  const pprofInitialPort = args["--pprofInitialPort"] || undefined;
  const configPath = args["--configPath"] || undefined;

  const chainProps = getChainProps({
    network,
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
    configPath,
  });
  await startAllRelayers({
    ...chainProps,
  });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
