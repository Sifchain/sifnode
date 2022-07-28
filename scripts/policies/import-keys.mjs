#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import { binary, keyringFlags } from "./helpers/constants.mjs";
import { getAccounts } from "./helpers/getAccounts.mjs";

const accounts = await getAccounts();

await spinner("import keys                              ", () =>
  within(async () => {
    // $.verbose = true;
    return Promise.all(
      accounts.map(async ({ key, mnemonic }) => {
        await $`echo ${mnemonic} | ${binary} keys add ${key} --recover ${keyringFlags}`;
      })
    );
  })
);
