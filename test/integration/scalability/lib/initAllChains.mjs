import { getChainProps } from "../utils/getChainProps.mjs";
import { initChain } from "./initChain.mjs";

const chains = require("../config/chains.json");

export async function initAllChains(props) {
  const { network, home = `/tmp/localnet` } = props;

  await $`rm -rf ${home}`;

  await Promise.all(
    Object.entries(chains)
      .filter(([_, { disabled = false }]) => disabled === false)
      .map(async ([chain]) => {
        const chainProps = getChainProps({ chain, network });

        return initChain({
          ...chainProps,
          home: `${home}/${chainProps.chain}/${chainProps.chainId}`,
        });
      })
  );
}
