#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import { binary } from "./helpers/constants.mjs";
import { getEntries } from "./helpers/getEntries.mjs";
import { getPools } from "./helpers/getPools.mjs";
import { pickRandomPools } from "./helpers/pickRandomPools.mjs";

// const nAccounts = 8000; // testnet
// const nPools = 10; // testnet
const nAccounts = 5; // localnet
const nPools = 2; // localnet

const entries = await getEntries();
const pools = await getPools();

const accounts = await spinner(
  "generate accounts                              ",
  () =>
    within(async () => {
      // $.verbose = true;
      return Promise.all(
        [...Array(nAccounts).keys()].map(async (id) => ({
          key: `account-${id + 1}`,
          mnemonic: (await $`${binary} keys mnemonic`).toString().trim(),
          pools: pickRandomPools(pools, entries, nPools),
        }))
      );
    })
);

await spinner("write accounts file                              ", () =>
  within(async () => {
    await fs.writeJson(`accounts.json`, accounts);
  })
);
