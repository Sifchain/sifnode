import JSBI from "jsbi";

import { AssetAmount } from "../../entities";
import createSifService, { SifServiceContext } from ".";
import { getBalance, getTestingTokens } from "../../test/utils/getTestingToken";
import { useStack } from "../../../../test/stack";

useStack("once");

const [ROWAN, CATK, CBTK, CETH] = getTestingTokens([
  "ROWAN",
  "CATK",
  "CBTK",
  "CETH",
]);

const TOKENS = {
  rowan: ROWAN,
  atk: CATK,
  btk: CBTK,
  eth: CETH,
};

// const badMnemonic =
//   "Ever have that feeling where you’re not sure if you’re awake or dreaming?";

const mnemonic =
  "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow";

// To be kept up to date with test state
const account = {
  address: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
  balance: [AssetAmount(TOKENS.rowan, "1000")],
  pubkey: {
    type: "tendermint/PubKeySecp256k1",
    value: "AvUEsFHbsr40nTSmWh7CWYRZHGwf4cpRLtJlaRO4VAoq",
  },
  accountNumber: 2,
  sequence: 1,
};

const testConfig: SifServiceContext = {
  sifAddrPrefix: "sif",
  sifApiUrl: "http://127.0.0.1:1317",
  sifWsUrl: "ws://127.0.0.1:26657/websocket",
  sifRpcUrl: "http://127.0.0.1:26657",
  assets: [ROWAN, CATK, CBTK, CETH],
  keplrChainConfig: {
    rest: "",
    rpc: "",
    chainId: "sifchain",
    chainName: "Sifchain",
    stakeCurrency: {
      coinDenom: "ROWAN",
      coinMinimalDenom: "rowan",
      coinDecimals: 18,
    },
    bip44: {
      coinType: 118,
    },
    bech32Config: {
      bech32PrefixAccAddr: "sif",
      bech32PrefixAccPub: "sifpub",
      bech32PrefixValAddr: "sifvaloper",
      bech32PrefixValPub: "sifvaloperpub",
      bech32PrefixConsAddr: "sifvalcons",
      bech32PrefixConsPub: "sifvalconspub",
    },
    currencies: [],
    feeCurrencies: [
      {
        coinDenom: "ROWAN",
        coinMinimalDenom: "rowan",
        coinDecimals: 18,
      },
    ],
    coinType: 118,
    gasPriceStep: {
      low: 0.01,
      average: 0.025,
      high: 0.04,
    },
  },
};

// This is redundant. CWalletActions and SifService should be combined imo
describe("sifService", () => {
  it("should use mnemeonic to sign into cosmos wallet", async () => {
    const sifService = createSifService(testConfig);
    // catch Error["Bad mnemonic"]??
    test.todo("more tests on bad mnemonic ");

    try {
      await sifService.setPhrase("");
    } catch (error) {
      expect(error).toEqual("No mnemonic. Can't generate wallet.");
    }

    expect(await sifService.setPhrase(mnemonic)).toBe(account.address);
  });

  it("should get cosmos balance", async () => {
    const sifService = createSifService(testConfig);
    await sifService.setPhrase(mnemonic);
    try {
      await sifService.getBalance("");
    } catch (error) {
      expect(error).toEqual("Address undefined. Fail");
    }
    try {
      await sifService.getBalance("asdfasdsdf");
    } catch (error) {
      expect(error).toEqual("Address not valid (length). Fail");
    }

    const balances = await sifService.getBalance(account.address);
    const balance = getBalance(balances, "rowan");
    expect(balance?.toString()).toEqual("100000000000000000000000000000 ROWAN");
  });

  it("should transfer transaction", async () => {
    const sifService = createSifService(testConfig);

    const address = await sifService.setPhrase(mnemonic);
    await sifService.transfer({
      amount: JSBI.BigInt("50"),
      asset: TOKENS.rowan,
      recipient: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
      memo: "",
    });
    const balances = await sifService.getBalance(address);
    const balance = getBalance(balances, "rowan");
    expect(balance?.toString()).toEqual("99999999999999999999999749950 ROWAN");
  });
});
