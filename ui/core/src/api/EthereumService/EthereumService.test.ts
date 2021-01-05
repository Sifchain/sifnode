// This test must be run in an environment that supports ganace

import createEthereumService from "./EthereumService";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { Asset, AssetAmount } from "../../entities";
import JSBI from "jsbi";
import B from "../../entities/utils/B";
import { getTestingTokens } from "../../test/utils/getTestingToken";

// ^ NOTE: we have had issues where truffle deploys contracts that cost a different amount of gas in CI versus locally.
// These test have been altered to be less deterministic as a consequence

function getBalance(balances: AssetAmount[], symbol: string) {
  const bal = balances.find(
    ({ asset }) => asset.symbol.toUpperCase() === symbol.toUpperCase()
  );
  if (!bal) throw new Error("Symbol not found in balances");
  return bal;
}

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

    const ethAmount = getBalance(balances, "eth").amount;
    const atkAmount = getBalance(balances, "atk").amount;
    const btkAmount = getBalance(balances, "btk").amount;

    // 98xxxxxxxxxxxxxxxxxx = 98ish eth ^ (see above)
    expect(/^98\d{18}/.test(ethAmount.toString())).toBe(true);
    expect(atkAmount.toString()).toEqual("10000000000");
    expect(btkAmount.toString()).toEqual("10000000000");
  });

  test("isConnected", async () => {
    expect(EthereumService.isConnected()).toBe(false);
    await EthereumService.connect();
    expect(EthereumService.isConnected()).toBe(true);
  });

  test("transfer ERC-20 to smart contract", async () => {
    await EthereumService.connect();
    const state = EthereumService.getState();

    const balances = state.balances;
    const account0AtkAmount = getBalance(balances, "atk").amount;

    expect(account0AtkAmount.toString()).toEqual("10000000000");

    await EthereumService.transfer({
      amount: B("10.000000", ATK.decimals),
      recipient: state.accounts[1],
      asset: ATK,
    });

    const account0NewAtkAmount = getBalance(
      await EthereumService.getBalance(),
      "atk"
    ).amount;
    const account1NewAtkAmount = getBalance(
      await EthereumService.getBalance(state.accounts[1]),
      "atk"
    ).amount;

    expect(account0NewAtkAmount.toString()).toEqual("9990000000");
    expect(account1NewAtkAmount.toString()).toEqual("10000000");
  });

  test("transfer ETH", async () => {
    await EthereumService.connect();
    const state = EthereumService.getState();

    const TEN_ETH = JSBI.BigInt(10 * 10 ** 18);

    await EthereumService.transfer({
      amount: TEN_ETH,
      recipient: state.accounts[1],
    });

    const balanceAccount1 = await EthereumService.getBalance(state.accounts[1]);
    const account1NewBalance = getBalance(balanceAccount1, "ETH").amount;

    expect(account1NewBalance.toString()).toEqual("110000000000000000000");
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
