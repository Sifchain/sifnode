import { initAllRelayers } from "../lib/initAllRelayers.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--network": String,
      "--configPath": String,
      "--registryFrom": Number,
      "--rpcInitialPort": Number,
      "--p2pInitialPort": Number,
      "--pprofInitialPort": Number,
    },
    `
Usage:

  yarn initAllRelayers [options]

Initiate all the IBC relayers connected to sifchain.

Options:

--network               Select a predifined network in chains.json
--configPath            Global directory for config and data of initiated chains
--registryFrom          Directory for storing global relayer registry data
--rpcInitialPort        Initial RPC port number
--p2pInitialPort        Initial P2P port number
--pprofInitialPort      Initial pprof port number
`
  );

  const network = args["--network"] || undefined;
  const configPath = args["--configPath"] || undefined;
  const registryFrom = args["--registryFrom"] || undefined;
  const rpcInitialPort = args["--rpcInitialPort"] || undefined;
  const p2pInitialPort = args["--p2pInitialPort"] || undefined;
  const pprofInitialPort = args["--pprofInitialPort"] || undefined;

  const chainProps = getChainProps({
    network,
    configPath,
    registryFrom,
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
  });

  await initAllRelayers({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
