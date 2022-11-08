import { getTx } from "../lib/getTx.mjs";
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
      "--hash": String,
    },
    `
Usage:

  yarn checkTx [options]

Check the transaction logs of any given transaction hash.

Options:

--chain     Select a predifined chain in chains.json
--network   Select a predifined network in chains.json
--node      Node address
--chain-id  Chain ID
--binary    Binary name of the chain
--hash      Transaction hash
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const node = args["--node"] || undefined;
  const chainId = args["--chain-id"] || undefined;
  const binary = args["--binary"] || undefined;
  const hash = args["--hash"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
    node,
    chainId,
    binary,
    hash,
  });
  const tx = await getTx({
    ...chainProps,
  });

  console.log(JSON.stringify(tx, null, 2));
}

if (process.env.NODE_ENV !== "test") {
  start();
}
