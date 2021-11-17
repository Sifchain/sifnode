import { initAllChains } from "../lib/initAllChains.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--network": String,
      "--home": String,
    },
    `
Usage:

  yarn initAllChains [options]

Initiate all the chains locally based on an existing remote chain.

Options:

--network   Select a predifined network in chains.json
--home      Global directory for config and data of initiated chains
`
  );

  const network = args["--network"] || undefined;
  const home = args["--home"] || undefined;

  const chainProps = getChainProps({
    network,
    home,
  });

  await initAllChains({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
