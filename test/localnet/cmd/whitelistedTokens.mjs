import { getWhitelistedTokens } from "../lib/getWhitelistedTokens.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--chain": String,
      "--network": String,
      "--node": String,
      "--chain-id": String,
      "--binary": String,
    },
    `
Usage:

  yarn whitelistedTokens [options]

Returns a list of all the whitelisted tokens available in the IBC chain.

Options:

--chain     Select a predifined chain in chains.json
--network   Select a predifined network in chains.json
--node      Node address
--chain-id  Chain ID
--binary    Binary name of the chain
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const node = args["--node"] || undefined;
  const chainId = args["--chain-id"] || undefined;
  const binary = args["--binary"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
    node,
    chainId,
    binary,
  });
  const tokens = await getWhitelistedTokens({
    ...chainProps,
  });

  console.log(JSON.stringify(tokens, null, 2));
}

if (process.env.NODE_ENV !== "test") {
  start();
}
