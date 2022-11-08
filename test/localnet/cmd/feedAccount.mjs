import { feedAccount } from "../lib/feedAccount.mjs";
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
      "--name": String,
      "--faucet": String,
    },
    `
Usage:

  yarn checkTx [options]

Fund the chain account.

Options:

--chain     Select a predifined chain in chains.json
--network   Select a predifined network in chains.json
--node      Node address
--chain-id  Chain ID
--binary    Binary name of the chain
--name      Account name or address
--facet     Faucet URL
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const node = args["--node"] || undefined;
  const chainId = args["--chain-id"] || undefined;
  const binary = args["--binary"] || undefined;
  const name = args["--name"] || undefined;
  const faucet = args["--faucet"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
    node,
    chainId,
    binary,
    name,
    faucet,
  });
  const result = await feedAccount({
    ...chainProps,
  });

  console.log(JSON.stringify(result, null, 2));
}

if (process.env.NODE_ENV !== "test") {
  start();
}
