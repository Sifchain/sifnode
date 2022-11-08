import { $ } from "zx";
import { getChainProps } from "../utils/getChainProps.mjs";
import { initChain } from "./initChain.mjs";

import { createRequire } from "module";
const require = createRequire(import.meta.url);
const chains = require("../config/chains.json");

export async function initAllChains({
  network,
  configPath = `/tmp/localnet/config`,
}) {
  await $`rm -rf ${configPath}`;

  await Promise.all(
    Object.entries(chains)
      .filter(([_, { disabled = false }]) => disabled === false)
      .map(async ([chain]) => {
        const chainProps = getChainProps({ chain, network });

        return initChain({
          ...chainProps,
          configPath: `${configPath}/${chainProps.chain}/${chainProps.chainId}`,
        });
      })
  );
}
