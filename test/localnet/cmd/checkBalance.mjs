import { getBalances } from "../lib/getBalances.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";
import { pickChains } from "../utils/pickChains.mjs";

export async function start() {
  const args = arg(
    {
      "--chain": String,
      "--network": String,
      "--node": String,
      "--chain-id": String,
      "--binary": String,
      "--name": String,
    },
    `
Usage:

  yarn checkBalance [options]

Check the balance of any given chain account name or address.

Options:

--chain     Select a predifined chain in chains.json
--network   Select a predifined network in chains.json
--node      Node address
--chain-id  Chain ID
--binary    Binary name of the chain
--name      Account name or address
--binPath   Location of binaries
`
  );

  const chain = args["--chain"] || undefined;
  const network = args["--network"] || undefined;
  const node = args["--node"] || undefined;
  const chainId = args["--chain-id"] || undefined;
  const binary = args["--binary"] || undefined;
  const name = args["--name"] || undefined;
  const binPath = args["--binPath"] || undefined;

  const chains = pickChains({ chain });

  for (let currentChain of chains) {
    console.log(`Chain "${currentChain}"`);
    const chainProps = getChainProps({
      chain: currentChain,
      network,
      node,
      chainId,
      binary,
      name,
      binPath,
    });
    const balances = await getBalances({
      ...chainProps,
    });

    balances.forEach(({ denom, amount }) => {
      console.log(`
  denom   ${denom}
  amount  ${amount}
  `);
    });
  }
}

if (process.env.NODE_ENV !== "test") {
  start();
}
