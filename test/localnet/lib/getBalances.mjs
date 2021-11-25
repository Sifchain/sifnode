import { $ } from "zx";
import { getAddress } from "./getAddress.mjs";

export async function getBalances({
  binary,
  name,
  node,
  chainId,
  binPath = "/tmp/localnet/bin",
  debug = false,
}) {
  if (!node) throw new Error("missing requirement argument: --node");
  if (!chainId) throw new Error("missing requirement argument: --chain-id");
  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!name) throw new Error("missing requirement argument: --name");
  if (!binPath) throw new Error("missing requirement argument: --binPath");

  const addr = await getAddress({ binary, name });

  const result = await $`
${binPath}/${binary} \
  q \
  bank \
  balances \
  ${addr} \
  --node ${node} \
  --chain-id ${chainId} \
  --output json`;

  const balances = JSON.parse(result).balances;

  if (debug) {
    console.log(`balances:`);
    console.log(balances);
  }

  return balances;
}
