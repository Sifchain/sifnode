import { sleep } from "zx";
import { getChainProps } from "../utils/getChainProps.mjs";
import { initRelayer } from "./initRelayer.mjs";
import { startChain } from "./startChain.mjs";

const chains = require("../config/chains.json");

export async function startAllChains(props) {
  const {
    network,
    home = `/tmp/localnet`,
    rpcPort = 11000,
    p2pPort = 12000,
    pprofPort = 13000,
  } = props;

  const chainsProps = (
    await Promise.all(
      Object.entries(chains)
        .filter(([_, { disabled = false }]) => disabled === false)
        .map(async ([chain], index) => {
          const chainProps = getChainProps({ chain, network });

          return startChain({
            ...chainProps,
            rpcPort: rpcPort + index,
            p2pPort: p2pPort + index,
            pprofPort: pprofPort + index,
            home: `${home}/${chainProps.chain}/${chainProps.chainId}`,
          });
        })
    )
  ).reduce((acc, cur) => ({ ...acc, [cur.chain]: cur }), {});

  await initRelayer({ chainsProps });

  // await Object.values(chainsProps).map(async ({ proc }) => {
  //   proc.kill("SIGINT");
  // });
}
