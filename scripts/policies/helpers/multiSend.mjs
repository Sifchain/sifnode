#!/usr/bin/env zx

import _ from "lodash";
import { spinner } from "zx/experimental";
import { getFee } from "./getFee.mjs";

export async function multiSend(
  signingClient,
  senderAccount,
  accounts,
  amount = "10000"
) {
  await spinner("multi sending tokens", () =>
    within(() =>
      signingClient.signAndBroadcast(
        senderAccount.address,
        _.flatten(
          accounts.map(({ account, pools }) => [
            ...pools.map(({ symbol: denom, decimals }) => ({
              typeUrl: "/cosmos.bank.v1beta1.MsgSend",
              value: {
                fromAddress: senderAccount.address,
                toAddress: account.address,
                amount: [{ denom, amount: `${amount}${"0".repeat(decimals)}` }],
              },
            })),
            {
              typeUrl: "/cosmos.bank.v1beta1.MsgSend",
              value: {
                fromAddress: senderAccount.address,
                toAddress: account.address,
                amount: [
                  { denom: "rowan", amount: `${amount}${"0".repeat(18)}` },
                ],
              },
            },
          ])
        ),
        getFee(accounts)
      )
    )
  );
}
