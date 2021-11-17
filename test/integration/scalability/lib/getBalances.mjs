import { $ } from "zx";
import { getAddress } from "./getAddress.mjs";

export async function getBalances({ binary, name, node, chainId }) {
  if (!node) throw new Error("missing requirement argument: --node");
  if (!chainId) throw new Error("missing requirement argument: --chain-id");
  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!name) throw new Error("missing requirement argument: --name");

  const addr = await getAddress({ binary, name });

  const result = await $`
${binary} \
  q \
  bank \
  balances \
  ${addr} \
  --node ${node} \
  --chain-id ${chainId} \
  --output json`;

  const balances = JSON.parse(result).balances;

  console.log(`balances:`);
  console.log(balances);

  return balances;
}
