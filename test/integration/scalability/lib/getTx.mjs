import { $ } from "zx";

export async function getTx({ binary, hash, node, chainId }) {
  if (!node) throw new Error("missing requirement argument: --node");
  if (!chainId) throw new Error("missing requirement argument: --chain-id");
  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!hash) throw new Error("missing requirement argument: --hash");

  const result = await $`
${binary} \
  q \
  tx \
  ${hash} \
  --node ${node} \
  --chain-id ${chainId} \
  --output json`;

  const tx = JSON.parse(result);

  return tx;
}
