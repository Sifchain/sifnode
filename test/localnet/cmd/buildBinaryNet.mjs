import { buildBinaryNet } from "../lib/buildBinaryNet.mjs";
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

  yarn buildBinaryNet [options]

Initiate two IBC chains locally based on an existing remote chain and take a snapshot.

Options:

--chain     Select a predifined chain in chains.json
--network   Select a predifined network in chains.json
--home      Global directory for config and data of initiated chains
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const home = args["--home"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
  });

  await buildBinaryNet({ chainProps, network, home });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
