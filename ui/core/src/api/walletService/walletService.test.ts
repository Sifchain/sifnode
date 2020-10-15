// This test must be run in an environment that supports ganace

import { getFakeTokens } from "../utils/getFakeTokens";
import createWalletService from ".";
import { getWeb3 } from "../../test/getWeb3";
import { Balance } from "../../entities";
import { ETH } from "../../constants";
import { wallet } from "src/store/wallet";

describe("walletService", () => {
  let walletService: ReturnType<typeof createWalletService>;

  beforeEach(async () => {
    const supportedTokens = await getFakeTokens();
    walletService = createWalletService({
      getWeb3,
      getSupportedTokens: async () => supportedTokens,
    });
  });
  test("it should connect without error", async () => {
    let causedError = false;
    try {
      await walletService.connect();
    } catch (err) {
      causedError = true;
    }
    expect(causedError).toBeFalsy();
  });
  test("that it returns the correct wallet amounts", async () => {
    const supportedTokens = await getFakeTokens();
    const walletService = createWalletService({
      getWeb3,
      getSupportedTokens: async () => supportedTokens,
    });
    await walletService.connect();
    const balances = await walletService.getBalance();

    const ATK = supportedTokens.find(({ symbol }) => symbol === "ATK");
    const BTK = supportedTokens.find(({ symbol }) => symbol === "BTK");

    expect(balances[0].toFixed()).toEqual(
      Balance.create(ETH, "99950481140000000000").toFixed()
    );
    expect(balances[1].toFixed()).toEqual(
      Balance.create(ATK!, "10000000000").toFixed()
    );
    expect(balances[2].toFixed()).toEqual(
      Balance.create(BTK!, "10000000000").toFixed()
    );
  });
  test("isConnected", async () => {
    expect(walletService.isConnected()).toBe(false);
    await walletService.connect();
  });
});
