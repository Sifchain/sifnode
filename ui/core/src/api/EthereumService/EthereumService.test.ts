// This test must be run in an environment that supports ganace

import createEthereumService from "./EthereumService";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { Asset } from "../../entities";
import JSBI from "jsbi";
import { getBalance, getTestingTokens } from "../../test/utils/getTestingToken";
import { useStack } from "../../../../test/stack";

// ^ NOTE: we have had issues where truffle deploys contracts that cost a different amount of gas in CI versus locally.
// These test have been altered to be less deterministic as a consequence

useStack("every-test");

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

    const ethAmount = getBalance(balances, "eth");
    const atkAmount = getBalance(balances, "atk");
    const btkAmount = getBalance(balances, "btk");

    // We dont know what the amount is going to be as it changes
    // depending on a bunch of factors so just checking for a string of digits
    expect(/^\d+/.test(ethAmount.toString())).toBeTruthy();
    expect(atkAmount.toBigInt().toString()).toEqual("10000000000000000000000");
    expect(btkAmount.toBigInt().toString()).toEqual("10000000000000000000000");
  });

  test("isConnected", async () => {
    expect(EthereumService.isConnected()).toBe(true);
  });

  test("transfer ERC-20 to smart contract", async () => {
    await EthereumService.connect();
    const state = EthereumService.getState();

    const balances = state.balances;
    const account0AtkAmount = getBalance(balances, "atk");

    expect(account0AtkAmount.toBigInt().toString()).toEqual(
      "10000000000000000000000",
    );

    await EthereumService.transfer({
      amount: JSBI.BigInt("10000000"),
      recipient: state.accounts[1],
      asset: ATK,
    });

    const account0NewAtkAmount = getBalance(
      await EthereumService.getBalance(),
      "atk",
    );
    const account1NewAtkAmount = getBalance(
      await EthereumService.getBalance(state.accounts[1]),
      "atk",
    );

    expect(account0NewAtkAmount.toBigInt().toString()).toEqual(
      "9999999999999990000000",
    );
    expect(account1NewAtkAmount.toBigInt().toString()).toEqual("10000000");
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
    const account1NewBalance = getBalance(balanceAccount1, "ETH");

    expect(account1NewBalance.toBigInt().toString()).toEqual(
      "110000000000000000000",
    );
  });

  it("should not do anything with phase and purgingClient", async () => {
    // TODO: We probably don't need this right now because we delegate to metamask
    expect(await EthereumService.setPhrase("testing one two three")).toEqual(
      "",
    );
    expect(EthereumService.purgeClient()).toBe(undefined);
  });
});
