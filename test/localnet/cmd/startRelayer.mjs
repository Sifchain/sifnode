import { startRelayer } from "../lib/startRelayer.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--chain": String,
      "--network": String,
      "--home": String,
    },
    `
Usage:

  yarn startRelayer [options]

Start a relayer locally to pair a foreign IBC chain to sifchain.

Options:

--chain         Select a predifined chain in chains.json
--network       Select a predifined network in chains.json
--home          Directory for config and data
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const home = args["--home"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
    home,
  });
  await startRelayer({
    ...chainProps,
  });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
