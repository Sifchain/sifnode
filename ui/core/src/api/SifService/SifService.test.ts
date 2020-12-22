import JSBI from "jsbi";

import { AssetAmount, Coin, Network } from "../../entities";
import createSifService, { SifServiceContext } from ".";

const TOKENS = {
  rowan: Coin({
    symbol: "rowan",
    decimals: 0,
    name: "Rowan",
    network: Network.SIFCHAIN,
  }),

  atk: Coin({
    symbol: "catk",
    decimals: 0,
    name: "catk",
    network: Network.SIFCHAIN,
  }),

  btk: Coin({
    symbol: "cbtk",
    decimals: 0,
    name: "cbtk",
    network: Network.SIFCHAIN,
  }),
};

// This is required because we need to wait for the blockchain to process transactions
jest.setTimeout(20000);

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
  assets: [],
};

function getBalance(balances: AssetAmount[], symbol: string): AssetAmount {
  const bal = balances.find((bal) => bal.asset.symbol === symbol);
  if (!bal) throw new Error("Asset not found in balances");
  return bal;
}

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
    expect(balance?.toFixed()).toEqual("1000000000");
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
    expect(balance?.toFixed()).toEqual("999999950");
  });
});
