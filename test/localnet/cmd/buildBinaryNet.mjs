import { buildBinaryNet } from "../lib/buildBinaryNet.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--chain": String,
      "--network": String,
      "--configPath": String,
    },
    `
Usage:

  yarn buildBinaryNet [options]

Initiate two IBC chains locally based on an existing remote chain and take a snapshot.

Options:

--chain           Select a predifined chain in chains.json
--network         Select a predifined network in chains.json
--configPath      Global directory for config and data of initiated chains
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const configPath = args["--configPath"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
  });

  await buildBinaryNet({ chainProps, network, configPath });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
