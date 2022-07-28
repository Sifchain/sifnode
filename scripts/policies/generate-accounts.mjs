#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import { binary } from "./helpers/constants.mjs";
import { getEntries } from "./helpers/getEntries.mjs";
import { getPools } from "./helpers/getPools.mjs";
import { pickRandomPools } from "./helpers/pickRandomPools.mjs";

const quantity = 8000;

const entries = await getEntries();
const pools = await getPools();

const accounts = await spinner(
  "generate accounts                              ",
  () =>
    within(async () => {
      // $.verbose = true;
      return Promise.all(
        [...Array(quantity).keys()].map(async (id) => ({
          key: `account-${id + 1}`,
          mnemonic: (await $`${binary} keys mnemonic`).toString().trim(),
          pools: pickRandomPools(pools, entries),
        }))
      );
    })
);

await spinner("write accounts file                              ", () =>
  within(async () => {
    await fs.writeJson(`accounts.json`, accounts);
  })
);
