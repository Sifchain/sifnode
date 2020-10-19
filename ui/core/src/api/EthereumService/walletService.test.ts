// This test must be run in an environment that supports ganace

import { getFakeTokens } from "./utils/getFakeTokens";
import createWalletService from ".";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { Balance } from "../../entities";
import { ETH } from "../../constants";
import { TEN } from "src/entities/fraction/Fraction";
import JSBI from "jsbi";
import B from "../../entities/utils/B";

describe("EthereumService", () => {
  let EthereumService: ReturnType<typeof createWalletService>;

  beforeEach(async () => {
    const supportedTokens = await getFakeTokens();
    EthereumService = createWalletService({
      getWeb3Provider,
      getSupportedTokens: async () => supportedTokens,
    });
  });

  test("it should connect without error", async () => {
    let causedError = false;
    try {
      await EthereumService.connect();
    } catch (err) {
      causedError = true;
    }
    expect(causedError).toBeFalsy();
  });

  test("that it returns the correct wallet amounts", async () => {
    const supportedTokens = await getFakeTokens();
    const EthereumService = createWalletService({
      getWeb3Provider,
      getSupportedTokens: async () => supportedTokens,
    });
    await EthereumService.connect();
    const balances = await EthereumService.getBalance();

    const ATK = supportedTokens.find(({ symbol }) => symbol === "ATK");
    const BTK = supportedTokens.find(({ symbol }) => symbol === "BTK");

    expect(balances[0].toFixed()).toEqual(
      // TODO: Work out a better way to deal with natural amounts eg 99.95048114 ETH
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
    expect(EthereumService.isConnected()).toBe(false);
    await EthereumService.connect();
    expect(EthereumService.isConnected()).toBe(true);
  });

  test("transfer ERC-20 to smart contract", async () => {
    const supportedTokens = await getFakeTokens();
    const EthereumService = createWalletService({
      getWeb3Provider,
      getSupportedTokens: async () => supportedTokens,
    });
    await EthereumService.connect();
    const state = EthereumService.getState();
    const origBalanceAccount0 = await EthereumService.getBalance();

    expect(
      origBalanceAccount0
        .find(({ asset: { symbol } }) => symbol === "ATK")
        ?.toFixed()
    ).toEqual("10000.000000");

    const ATK = supportedTokens.find(({ symbol }) => symbol === "ATK");

    await EthereumService.transfer({
      amount: B("10.000000", ATK!.decimals),
      recipient: state.accounts[1],
      asset: ATK,
    });

    const balanceAccount0 = await EthereumService.getBalance();
    const balanceAccount1 = await EthereumService.getBalance(state.accounts[1]);

    expect(
      balanceAccount0
        .find(({ asset: { symbol } }) => symbol === "ATK")
        ?.toFixed(2)
    ).toEqual("9990.00");

    expect(
      balanceAccount1
        .find(({ asset: { symbol } }) => symbol === "ATK")
        ?.toFixed(2)
    ).toEqual("10.00");
  });

  test("transfer ETH", async () => {
    const supportedTokens = await getFakeTokens();
    const EthereumService = createWalletService({
      getWeb3Provider,
      getSupportedTokens: async () => supportedTokens,
    });
    await EthereumService.connect();
    const state = EthereumService.getState();
    const origBalanceAccount0 = await EthereumService.getBalance();

    expect(
      origBalanceAccount0
        .find(({ asset: { symbol } }) => symbol === "ETH")
        ?.toFixed()
    ).toEqual(Balance.n(ETH, "99.95048114").toFixed());

    await EthereumService.transfer({
      amount: JSBI.BigInt(10 * 10 ** 18),
      recipient: state.accounts[1],
    });

    const balanceAccount0 = await EthereumService.getBalance();
    const balanceAccount1 = await EthereumService.getBalance(state.accounts[1]);

    expect(
      balanceAccount0
        .find(({ asset: { symbol } }) => symbol === "ETH")
        ?.toFixed()
    ).toEqual("89.950061140000000000"); // Including gas

    expect(
      balanceAccount1
        .find(({ asset: { symbol } }) => symbol === "ETH")
        ?.toFixed()
    ).toEqual("110.000000000000000000");
  });
});
