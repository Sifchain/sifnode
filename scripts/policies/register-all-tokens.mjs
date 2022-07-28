#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import { binary, flags, gasFlags } from "./helpers/constants.mjs";
import { getAccountNumber } from "./helpers/getAccountNumber.mjs";
import { getEntries } from "./helpers/getEntries.mjs";

const { ADMIN_KEY } = process.env;

const path = `/tmp/denoms`;

const entries = await getEntries();

await spinner("write denom files                              ", () =>
  within(async () => {
    await $`rm -rf ${path}`;
    await $`mkdir -p ${path}`;
    await Promise.all(
      entries.map((entry) =>
        fs.writeJson(`${path}/${entry.denom.replace("/", "_")}.json`, {
          entries: [entry],
        })
      )
    );
  })
);

const { accountNumber, sequence } = await getAccountNumber(ADMIN_KEY);

await spinner("register tokens                              ", () =>
  within(async () => {
    // $.verbose = true;
    let seq = sequence;
    for (let { denom } of entries) {
      const denomPath = `${path}/${denom.replace("/", "_")}.json`;
      await $`${binary} tx tokenregistry register ${denomPath} ${flags} ${gasFlags} --broadcast-mode=async --account-number=${accountNumber} --sequence=${seq}`;
      seq++;
    }
  })
);
