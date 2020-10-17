// This test must be run in an environment that supports ganace

import { getFakeTokens } from "./utils/getFakeTokens";
import createWalletService from ".";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { Balance } from "../../entities";
import { ETH } from "../../constants";

describe("EtheriumService", () => {
  let EtheriumService: ReturnType<typeof createWalletService>;

  beforeEach(async () => {
    const supportedTokens = await getFakeTokens();
    EtheriumService = createWalletService({
      getWeb3Provider,
      getSupportedTokens: async () => supportedTokens,
    });
  });
  test("it should connect without error", async () => {
    let causedError = false;
    try {
      await EtheriumService.connect();
    } catch (err) {
      causedError = true;
    }
    expect(causedError).toBeFalsy();
  });
  test("that it returns the correct wallet amounts", async () => {
    const supportedTokens = await getFakeTokens();
    const EtheriumService = createWalletService({
      getWeb3Provider,
      getSupportedTokens: async () => supportedTokens,
    });
    await EtheriumService.connect();
    const balances = await EtheriumService.getBalance();

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
    expect(EtheriumService.isConnected()).toBe(false);
    await EtheriumService.connect();
  });
});
