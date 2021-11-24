import { getChains } from "../utils/getChains.mjs";
import { getChainsProps } from "../utils/getChainsProps.mjs";
import { runRelayer } from "../utils/runRelayer.mjs";

export async function startAllRelayers({
  network,
  configPath = `/tmp/localnet/config`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  // 0) retrieve chains + metadata
  const chains = getChains({
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
    configPath,
  });
  const chainsProps = getChainsProps({ chains, network });
  const { sifchain: sifChainProps, ...otherChainsProps } = chainsProps;

  // 1) start relayers
  return Promise.all(
    Object.values(otherChainsProps).map(async ({ home }) => {
      return runRelayer({ home });
    })
  );
}
