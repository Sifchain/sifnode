import { send } from "../lib/send.mjs";
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
      "--src": String,
      "--dst": String,
      "--denom": String,
      "--amount": Number,
      "--fees": Number,
      "--dry-run": Boolean,
      "--binPath": String,
    },
    `
Usage:

  yarn send [options]

Transfer any given amount between two addresses of a same chain.

Options:

--chain     Select a predifined chain in chains.json
--network   Select a predifined network in chains.json
--node      Node address
--chain-id  Chain ID
--binary    Binary name of the chain
--src       Issuer address
--dst       Receiver address
--denom     Chain denom
--amount    Amount to send to receiver account
--fees      Minimum required fees amount to pay
--dry-run   Dry run
--binPath   Directory for binaries location
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const node = args["--node"] || undefined;
  const chainId = args["--chain-id"] || undefined;
  const binary = args["--binary"] || undefined;
  const src = args["--src"] || undefined;
  const dst = args["--dst"] || undefined;
  const denom = args["--denom"] || undefined;
  const amount = args["--amount"] || undefined;
  const fees = args["--fees"] || undefined;
  const dryRun = args["--dry-run"] || undefined;
  const binPath = args["--binPath"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
    node,
    chainId,
    binary,
    src,
    dst,
    denom,
    amount,
    fees,
    dryRun,
    binPath,
  });
  await send({
    ...chainProps,
  });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
