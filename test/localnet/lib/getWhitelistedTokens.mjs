import { $ } from "zx";

export async function getWhitelistedTokens({ node, chainId, binary }) {
  if (!node) throw new Error("missing requirement argument: --node");
  if (!chainId) throw new Error("missing requirement argument: --chain-id");
  if (!binary) throw new Error("missing requirement argument: --binary");

  const result = await $`
  ${binary} \
    q \
    tokenregistry \
    entries \
    --node ${node} \
    --chain-id ${chainId} \
    --output json`;

  const tokens = JSON.parse(result);

  return tokens.entries
    .filter(
      ({
        is_whitelisted,
        permissions,
        ibc_channel_id,
        ibc_counterparty_channel_id,
      }) =>
        is_whitelisted === true &&
        permissions.includes("IBCIMPORT") &&
        permissions.includes("IBCEXPORT") &&
        ibc_channel_id !== "" &&
        ibc_counterparty_channel_id !== ""
    )
    .map(
      ({
        decimals,
        denom,
        base_denom,
        ibc_channel_id,
        ibc_counterparty_channel_id,
        display_name,
        display_symbol,
        external_symbol,
        permissions,
      }) => ({
        decimals,
        denom,
        base_denom,
        ibc_channel_id,
        ibc_counterparty_channel_id,
        display_name,
        display_symbol,
        external_symbol,
        permissions,
      })
    );
}
