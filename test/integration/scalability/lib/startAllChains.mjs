import { getChains } from "../utils/getChains.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";
import { startChain } from "./startChain.mjs";

export async function startAllChains({
  network,
  home = `/tmp/localnet`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  const chains = getChains({
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
    home,
  });

  return Promise.all(
    Object.entries(chains)
      .filter(([_, { disabled = false }]) => disabled === false)
      .map(async ([chain, chainProps]) => {
        const newChainProps = getChainProps({
          chain,
          network,
          ...chainProps,
        });
        return startChain({
          ...newChainProps,
        });
      })
  );
}
