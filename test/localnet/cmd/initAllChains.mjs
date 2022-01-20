import { initAllChains } from "../lib/initAllChains.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--network": String,
      "--configPath": String,
    },
    `
Usage:

  yarn initAllChains [options]

Initiate all the chains locally based on an existing remote chain.

Options:

--network         Select a predifined network in chains.json
--configPath      Global directory for config and data of initiated chains
`
  );

  const network = args["--network"] || undefined;
  const configPath = args["--configPath"] || undefined;

  const chainProps = getChainProps({
    network,
    configPath,
  });

  await initAllChains({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
