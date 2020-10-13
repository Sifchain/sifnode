// This test must be run in an environment that supports ganace

import { getSupportedTokens } from "./utils/getSupportedTokens";
import createWalletService from "./walletService";
import { getWeb3 } from "../test/getWeb3";
import { AssetAmount } from "../entities";
import { ETH } from "../constants";

test("it should connect without error", async () => {
  const web3 = await getWeb3();
  const supportedTokens = await getSupportedTokens(web3);
  const walletService = createWalletService({
    getWeb3,
    getSupportedTokens: async () => supportedTokens,
  });

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
  const supportedTokens = await getSupportedTokens(web3);
  const walletService = createWalletService({
    getWeb3,
    getSupportedTokens: async () => supportedTokens,
  });

  const balances = await walletService.getAssetBalances();

  const ATK = supportedTokens.get("ATK");
  const BTK = supportedTokens.get("BTK");

  expect(balances[0].toFixed()).toEqual(
    AssetAmount.create(ETH, "99950481140000000000").toFixed()
  );
  expect(balances[1].toFixed()).toEqual(
    AssetAmount.create(ATK, "10000000000").toFixed()
  );
  expect(balances[2].toFixed()).toEqual(
    AssetAmount.create(BTK, "10000000000").toFixed()
  );
});
