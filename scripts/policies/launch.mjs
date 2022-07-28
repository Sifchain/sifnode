#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import { binary, keyringFlags } from "./helpers/constants.mjs";
import { getEntries } from "./helpers/getEntries.mjs";
import { getFunds } from "./helpers/getFunds.mjs";
import { getTokens } from "./helpers/getTokens.mjs";

const {
  HOME,
  ADMIN_ADDRESS,
  ADMIN_MNEMONIC,
  ADMIN_KEY,
  USER_ADDRESS,
  USER_MNEMONIC,
  USER_KEY,
} = process.env;

await spinner("build binary                              ", () =>
  within(async () => {
    await $`rm -rf ${HOME}/.sifnoded`;
    cd("../../");
    await $`make clean install`;
  })
);

await spinner("init chain                              ", () =>
  within(async () => {
    // $.verbose = true;
    await $`${binary} init test --chain-id=localnet -o`;
  })
);

await spinner(
  "generating deterministic accounts                              ",
  () =>
    within(async () => {
      // $.verbose = true;
      await $`echo ${ADMIN_MNEMONIC} | ${binary} keys add ${ADMIN_KEY} --recover ${keyringFlags}`;
      await $`echo ${USER_MNEMONIC} | ${binary} keys add ${USER_KEY} --recover ${keyringFlags}`;
      await $`${binary} keys add mkey --multisig ${ADMIN_KEY},${USER_KEY} --multisig-threshold 2 ${keyringFlags}`;
    })
);

const entries = await getEntries();
const tokens = await getTokens(entries);
const funds = getFunds(tokens);

await spinner("add accounts to genesis                              ", () =>
  within(async () => {
    // $.verbose = true;
    await $`${binary} add-genesis-account ${ADMIN_ADDRESS} ${funds} ${keyringFlags}`;
    await $`${binary} add-genesis-account ${USER_ADDRESS} ${funds} ${keyringFlags}`;
  })
);

await spinner("set admin account                              ", () =>
  within(async () => {
    // $.verbose = true;
    await $`${binary} add-genesis-clp-admin ${ADMIN_ADDRESS} ${keyringFlags}`;
    await $`${binary} set-genesis-oracle-admin ${ADMIN_KEY} ${keyringFlags}`;
    await $`${binary} add-genesis-validators $(${binary} keys show ${ADMIN_KEY} -a --bech val ${keyringFlags}) ${keyringFlags}`;
    await $`${binary} set-genesis-whitelister-admin ${ADMIN_KEY} ${keyringFlags}`;
  })
);

await spinner("generate tx                              ", () =>
  within(async () => {
    // $.verbose = true;
    await $`${binary} gentx ${ADMIN_KEY} 1000000000000000000000000stake --chain-id=localnet ${keyringFlags}`;
  })
);

await spinner("collecting genesis txs                              ", () =>
  within(async () => {
    // $.verbose = true;
    await $`${binary} collect-gentxs`;
  })
);

await spinner("validating genesis file                              ", () =>
  within(async () => {
    // $.verbose = true;
    await $`${binary} validate-genesis`;
  })
);

await spinner("run chain                              ", () =>
  within(async () => {
    // $.verbose = true;
    await $`killall sifnoded`.nothrow();
    await $`${binary} start --trace`;
  })
);
