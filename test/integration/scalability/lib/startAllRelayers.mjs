import { getChains } from "../utils/getChains.mjs";
import { getChainsProps } from "../utils/getChainsProps.mjs";

export async function startAllRelayers({
  network,
  home = `/tmp/localnet`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  // 0) retrieve chains + metadata
  const chains = getChains({
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
    home,
  });
  const chainsProps = getChainsProps({ chains, network });
  const { sifchain: sifChainProps, ...otherChainsProps } = chainsProps;

  return Promise.all(
    Object.values(otherChainsProps).map(async ({ home }) => {
      const relayerHome = `${home}/relayer`;

      const proc = await nothrow(
        $`ibc-relayer start -v --poll 10 --home ${relayerHome}`
      );

      return {
        proc,
      };
    })
  );
}
