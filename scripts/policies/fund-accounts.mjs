#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import {
  binary,
  flags,
  feesFlags,
  keyringFlags,
} from "./helpers/constants.mjs";
import { getAccounts } from "./helpers/getAccounts.mjs";
import { getAccountNumber } from "./helpers/getAccountNumber.mjs";

const { ADMIN_KEY } = process.env;

const accounts = await getAccounts();

const { accountNumber, sequence } = await getAccountNumber(ADMIN_KEY);

await spinner("fund accounts                              ", () =>
  within(async () => {
    // $.verbose = true;
    let seq = sequence;
    for (let { key, pools } of accounts) {
      for (let { symbol, decimals } of pools) {
        await $`\
${binary} tx bank send \
  ${ADMIN_KEY} \
  $(${binary} keys show ${key} -a ${keyringFlags}) \
  10000${"0".repeat(decimals)}${symbol} \
  ${flags} \
  ${feesFlags} \
  --broadcast-mode=async \
  --account-number=${accountNumber} \
  --sequence=${seq}`;
        seq++;
      }
      await $`\
${binary} tx bank send \
  ${ADMIN_KEY} \
  $(${binary} keys show ${key} -a ${keyringFlags}) \
  1000000${"0".repeat(18)}rowan \
  ${flags} \
  ${feesFlags} \
  --broadcast-mode=async \
  --account-number=${accountNumber} \
  --sequence=${seq}`;
      seq++;
    }
  })
);
