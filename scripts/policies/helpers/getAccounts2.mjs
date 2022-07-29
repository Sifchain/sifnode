import { spinner } from "zx/experimental";

import { pickRandomPools } from "./pickRandomPools.mjs";

const { DirectSecp256k1HdWallet } = require("@cosmjs/proto-signing");
const { SifSigningStargateClient } = require("@sifchain/stargate");

const { SIFNODE_NODE } = process.env;

export async function getAccounts2(
  pools,
  entries,
  nAccounts = 100,
  nPools = 10
) {
  const accounts = await spinner(
    "loading accounts                              ",
    () =>
      within(() =>
        Promise.all(
          [...Array(nAccounts).keys()].map(async () => {
            const wallet = await DirectSecp256k1HdWallet.generate(12, {
              prefix: "sif",
            });
            const [account] = await wallet.getAccounts();
            const signingClient =
              await SifSigningStargateClient.connectWithSigner(
                SIFNODE_NODE,
                wallet
              );
            return {
              wallet,
              account,
              signingClient,
              pools: pickRandomPools(pools, entries, nPools),
            };
          })
        )
      )
  );
  return accounts;
}
