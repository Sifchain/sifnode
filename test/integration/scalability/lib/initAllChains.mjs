import { $ } from "zx";
import { getChainProps } from "../utils/getChainProps.mjs";
import { initChain } from "./initChain.mjs";

import { createRequire } from "module";
const require = createRequire(import.meta.url);
const chains = require("../config/chains.json");

export async function initAllChains({ network, home = `/tmp/localnet` }) {
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
