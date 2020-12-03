// This test must be run in an environment that supports ganace

import { getFakeTokens } from "./utils/getFakeTokens";
import createEthereumService from ".";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { Asset, AssetAmount, Network, Token } from "../../entities";

import JSBI from "jsbi";
import B from "../../entities/utils/B";

describe("EthereumService", () => {
  let EthereumService: ReturnType<typeof createEthereumService>;
  let ATK: Asset;
  let BTK: Asset;
  let ETH: Asset;

  beforeEach(async () => {
    const supportedTokens = await getFakeTokens();
    const a = supportedTokens.find(
      ({ symbol }) => symbol.toUpperCase() === "ATK"
    );
    const b = supportedTokens.find(
      ({ symbol }) => symbol.toUpperCase() === "BTK"
    );
    const c = supportedTokens.find(
      ({ symbol }) => symbol.toUpperCase() === "ETH"
    );

    if (!a) throw new Error("ATK not returned");
    if (!b) throw new Error("BTK not returned");
    if (!c) throw new Error("ETH not returned");
    ATK = a;
    BTK = b;
    ETH = c;

    EthereumService = createEthereumService({
      getWeb3Provider,
      loadAssets: async () => supportedTokens,
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
    // const supportedTokens = await getFakeTokens();
    // const EthereumService = createEthereumService({
    //   getWeb3Provider,
    //   loadAssets: async () => supportedTokens,
    // });
    await EthereumService.connect();

    const balances = await EthereumService.getBalance();

    expect(balances[0].toFixed()).toEqual(
      // TODO: Work out a better way to deal with natural amounts eg 99.95048114 ETH
      AssetAmount(ETH, "99.950481140000000000").toFixed()
    );
    expect(balances[1].toFixed()).toEqual(
      AssetAmount(ATK, "10000.000000").toFixed()
    );
    expect(balances[2].toFixed()).toEqual(
      AssetAmount(BTK, "10000.000000").toFixed()
    );
  });

  test("isConnected", async () => {
    expect(EthereumService.isConnected()).toBe(false);
    await EthereumService.connect();
    expect(EthereumService.isConnected()).toBe(true);
  });

  test("transfer ERC-20 to smart contract", async () => {
    await EthereumService.connect();
    const state = EthereumService.getState();
    const origBalanceAccount0 = await EthereumService.getBalance();

    expect(
      origBalanceAccount0
        .find(({ asset: { symbol } }) => symbol.toUpperCase() === "ATK")
        ?.toFixed()
    ).toEqual("10000.000000");

    await EthereumService.transfer({
      amount: B("10.000000", ATK.decimals),
      recipient: state.accounts[1],
      asset: ATK,
    });

    const balanceAccount0 = await EthereumService.getBalance();
    const balanceAccount1 = await EthereumService.getBalance(state.accounts[1]);

    expect(
      balanceAccount0
        .find(({ asset: { symbol } }) => symbol.toUpperCase() === "ATK")
        ?.toFixed(2)
    ).toEqual("9990.00");

    expect(
      balanceAccount1
        .find(({ asset: { symbol } }) => symbol.toUpperCase() === "ATK")
        ?.toFixed(2)
    ).toEqual("10.00");
  });

  test("transfer ETH", async () => {
    await EthereumService.connect();
    const state = EthereumService.getState();
    const origBalanceAccount0 = await EthereumService.getBalance();

    expect(
      origBalanceAccount0
        .find(({ asset: { symbol } }) => symbol.toUpperCase() === "ETH")
        ?.toFixed()
    ).toEqual(AssetAmount(ETH, "99.95048114").toFixed());

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
  it("should disconnect", async () => {
    await EthereumService.disconnect();
    expect(EthereumService.getState().accounts).toEqual([]);
    expect(EthereumService.getState().connected).toBe(false);
    expect(EthereumService.getState().address).toBe("");
    expect(EthereumService.getState().balances).toEqual([]);
  });
  it("should not do anything with phase and purgingClient", async () => {
    // TODO: We probably don't need this right now because we delegate to metamask
    expect(await EthereumService.setPhrase("testing one two three")).toEqual(
      ""
    );
    expect(EthereumService.purgeClient()).toBe(undefined);
  });
});
