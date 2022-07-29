#!/usr/bin/env zx

import _ from "lodash";
import { spinner } from "zx/experimental";
import { getFee } from "./getFee.mjs";

export async function multiAddLiquidity(accounts, amount = "100") {
  await spinner("multi add liquidity", () =>
    within(async () => {
      for (let { account, pools, signingClient } of accounts) {
        process.stdout.write(`${" ".repeat(30)}${account.address}\r`);
        const msgs = pools.map(({ symbol, decimals, swapPriceExternal }) => ({
          typeUrl: "/sifnode.clp.v1.MsgAddLiquidity",
          value: {
            signer: account.address,
            externalAsset: {
              symbol,
            },
            nativeAssetAmount: `${parseInt(
              Number(amount) * Number(swapPriceExternal)
            )}${"0".repeat(18)}`,
            externalAssetAmount: `${amount}${"0".repeat(decimals)}`,
          },
        }));
        console.log(JSON.stringify(msgs, null, 2));
        const checkTx = await signingClient.signAndBroadcast(
          account.address,
          msgs,
          getFee(accounts)
        );
        console.log(checkTx);
      }
    })
  );
}
