import { fetch } from "zx";
import { getAddress } from "./getAddress.mjs";

export async function feedAccount({ binary, node, chainId, name, faucet }) {
  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!node) throw new Error("missing requirement argument: --node");
  if (!chainId) throw new Error("missing requirement argument: --chain-id");
  if (!name) throw new Error("missing requirement argument: --name");
  if (!faucet) throw new Error("missing requirement argument: --faucet");

  const address = await getAddress({ binary, name });

  const response = await fetch(faucet, {
    method: "post",
    body: JSON.stringify({ address: `${address}`.replace("\n", "") }),
  });
  const data = await response.json();

  return data;
}
