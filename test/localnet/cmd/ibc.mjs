import { ibc } from "../lib/ibc.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--times": Number,
      "--broadcast": String,
      "--chain": String,
      "--network": String,
      "--type": String,
      "--channelId": String,
      "--counterpartyChannelId": String,
      "--node": String,
      "--chain-id": String,
      "--binary": String,
      "--name": String,
      "--denom": String,
      "--amount": Number,
      "--fees": Number,
      "--gas": Number,
      "--timeout": Number,
      "--dry-run": Boolean,
    },
    `
Usage:

  yarn send [options]

Transfer any given amount between two addresses of a same chain.

Options:

--times                         Number of times to run the test
--broadcast [async|sync|block]  Select a broadcast mode
--chain                         Select a predifined chain in chains.json
--network                       Select a predifined network in chains.json
--type [issuer|receiver]        Selected chain is an issuer or receiver
--channelId                     Channel ID
--counterpartyChannelId         Counterparty channel ID
--node                          Node address
--chain-id                      Chain ID
--binary                        Binary name of the chain
--name                          Issuer and receiver account name
--denom                         Chain denom
--amount                        Amount to send to receiver account
--fees                          Minimum required fees amount to pay
--gas                           Minimum required gas amount
--timeout                       Packet timeout timestamp
--dry-run                       Dry run
`
  );

  const times = args["--times"] || undefined;
  const broadcast = args["--broadcast"] || undefined;
  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const type = args["--type"] || undefined;
  const channelId = args["--channelId"] || undefined;
  const counterpartyChannelId = args["--counterpartyChannelId"] || undefined;
  const node = args["--node"] || undefined;
  const chainId = args["--chain-id"] || undefined;
  const binary = args["--binary"] || undefined;
  const name = args["--name"] || undefined;
  const denom = args["--denom"] || undefined;
  const amount = args["--amount"] || undefined;
  const fees = args["--fees"] || undefined;
  const gas = args["--gas"] || undefined;
  const timeout = args["--timeout"] || undefined;
  const dryRun = args["--dry-run"] || undefined;

  const chainProps = getChainProps({
    chain,
    network,
    type,
    times,
    broadcast,
    node,
    chainId,
    binary,
    channelId,
    counterpartyChannelId,
    name,
    denom,
    amount,
    fees,
    gas,
    timeout,
    dryRun,
  });
  await ibc({
    ...chainProps,
  });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
