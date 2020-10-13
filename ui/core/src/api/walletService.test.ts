// This test must be run in an environment that supports ganace

import { getSupportedTokens } from "./utils/getSupportedTokens";
import { createWalletService } from "./walletService";
import { getWeb3 } from "../../test/getWeb3";

test("it should connect without error", async () => {
  const web3 = await getWeb3();
  const walletService = createWalletService(
    getWeb3,
    await getSupportedTokens(web3)
  );

  let causedError = false;
  try {
    await walletService.getAssetBalances();
  } catch (err) {
    causedError = true;
  }
  expect(causedError).toBeFalsy();
});

test("that it returns the correct wallet amounts", async () => {
  const web3 = await getWeb3();
  const walletService = createWalletService(
    getWeb3,
    await getSupportedTokens(web3)
  );
  expect(await walletService.getAssetBalances()).toMatchObject([
    {
      amount: "99950481140000000000",
      asset: { decimals: 18, name: "Etherium", symbol: "ETH" },
    },
    {
      amount: "10000000000",
      asset: { decimals: 6, name: "AliceToken", symbol: "ATK" },
    },
    {
      amount: "10000000000",
      asset: { decimals: 6, name: "BobToken", symbol: "BTK" },
    },
  ]);
});
