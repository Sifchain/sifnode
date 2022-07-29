#!/usr/bin/env zx

$.verbose = false;

import { getAccounts2 } from "./helpers/getAccounts2.mjs";
import { getSigningClient } from "./helpers/getSigningClient.mjs";
import { multiSend } from "./helpers/multiSend.mjs";
import { getEntries } from "./helpers/getEntries.mjs";
import { getPools } from "./helpers/getPools.mjs";
import { multiAddLiquidity } from "./helpers/multiAddLiquidity.mjs";

const { ADMIN_MNEMONIC } = process.env;

const entries = await getEntries();
const pools = await getPools();

const { account: adminAccount, signingClient: adminSigningClient } =
  await getSigningClient(ADMIN_MNEMONIC);

const accounts = await getAccounts2(pools, entries, 1, 1);

console.log(accounts.map(({ account: { address } }) => address));

await multiSend(adminSigningClient, adminAccount, accounts);

await multiAddLiquidity(accounts);
