#!/usr/bin/env zx

$.verbose = false;

import { spinner } from "zx/experimental";
import { binary, flags, feesFlags } from "./helpers/constants.mjs";
import { getAccountNumber } from "./helpers/getAccountNumber.mjs";
import { getPools } from "./helpers/getPools.mjs";

const { ADMIN_KEY } = process.env;

const pools = await getPools();

const { accountNumber, sequence } = await getAccountNumber(ADMIN_KEY);

await spinner("create pools                              ", () =>
  within(async () => {
    // $.verbose = true;
    let seq = sequence;
    for (let {
      external_asset: { symbol: externalAsset },
      native_asset_balance: nativeAssetBalance,
      external_asset_balance: externalAssetBalance,
    } of pools) {
      await $`${binary} tx clp create-pool --symbol=${externalAsset} --nativeAmount=${nativeAssetBalance} --externalAmount=${externalAssetBalance} ${flags} ${feesFlags} --broadcast-mode=async --account-number=${accountNumber} --sequence=${seq}`;
      seq++;
    }
  })
);
