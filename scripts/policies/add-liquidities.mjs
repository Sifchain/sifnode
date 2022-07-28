#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import { binary, feesFlags, userFlags } from "./helpers/constants.mjs";
import { getAccounts } from "./helpers/getAccounts.mjs";
import { getAccountNumber } from "./helpers/getAccountNumber.mjs";

const accounts = await getAccounts();

await spinner("accounts add liquidity                              ", () =>
  within(async () => {
    $.verbose = true;
    for (let { key, pools } of accounts) {
      const { accountNumber, sequence } = await getAccountNumber(key);
      let seq = sequence;
      for (let { symbol, decimals, swapPriceExternal } of pools) {
        await $`\
${binary} tx clp add-liquidity \
  --symbol=${symbol} \
  --nativeAmount=${parseInt(100 * Number(swapPriceExternal))}${"0".repeat(18)} \
  --externalAmount=100${"0".repeat(decimals)} \
  --from=${key} \
  ${userFlags} \
  ${feesFlags} \
  --broadcast-mode=async \
  --account-number=${accountNumber} \
  --sequence=${seq}`;
        seq++;
      }
    }
  })
);
