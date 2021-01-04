// This test must be run in an environment that supports ganace

import createEthereumService from "./EthereumService";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { Asset, AssetAmount } from "../../entities";
import JSBI from "jsbi";
import B from "../../entities/utils/B";
import { getTestingTokens } from "../../test/utils/getTestingToken";

describe("EthereumService", () => {
  let EthereumService: ReturnType<typeof createEthereumService>;
  let ATK: Asset;
  let BTK: Asset;
  let ETH: Asset;

  beforeEach(async () => {
    [ATK, BTK, ETH] = getTestingTokens(["ATK", "BTK", "ETH"]);

    EthereumService = createEthereumService({
      getWeb3Provider,
      assets: [ATK, BTK, ETH],
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
    await EthereumService.connect();

    const balances = EthereumService.getState().balances;

    expect(balances[0].toFixed()).toEqual(
      // TODO: Work out a better way to deal with natural amounts eg 99.95048114 ETH
      // AssetAmount(ETH, "99.950481140000000000").toFixed()
      AssetAmount(ETH, "99.700747008430000000").toFixed()
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
    const origBalanceAccount0 = state.balances;

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
      // ).toEqual(AssetAmount(ETH, "99.95048114").toFixed());
    ).toEqual(AssetAmount(ETH, "99.700747008430000000").toFixed());

    await EthereumService.transfer({
      amount: JSBI.BigInt(10 * 10 ** 18),
      recipient: state.accounts[1],
    });

    const balanceAccount0 = await EthereumService.getBalance();
    const balanceAccount1 = await EthereumService.getBalance(state.accounts[1]);

    expect(
      balanceAccount0
        .find(({ asset: { symbol } }) => symbol.toUpperCase() === "ETH")
        ?.toFixed()
      // ).toEqual("89.950061140000000000"); // Including gas
    ).toEqual("89.700327008430000000"); // Including gas

    expect(
      balanceAccount1
        .find(({ asset: { symbol } }) => symbol.toUpperCase() === "ETH")
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
