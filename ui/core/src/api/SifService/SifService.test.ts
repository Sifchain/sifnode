import JSBI from "jsbi";
import { NCN } from "../../constants";
import { Balance } from "../../entities";
import createSifService, { SifServiceContext } from ".";

const badMnemonic =
  "Ever have that feeling where you’re not sure if you’re awake or dreaming?";

const mnemonic =
  "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow";

const account = {
  address: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
  balance: [Balance.n(NCN, "1000")],
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
    const [balance] = await sifService.getBalance(account.address);
    expect(balance).toEqual(account.balance[0]);
  });

  // Skipping until we can get deterministic ordering
  it("should transfer transaction", async () => {
    const sifService = createSifService(testConfig);

    const address = await sifService.setPhrase(mnemonic);
    await sifService.transfer({
      amount: JSBI.BigInt("50"),
      asset: NCN,
      recipient: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
      memo: "",
    });
    const [balance] = await sifService.getBalance(address);
    expect(balance).toEqual(Balance.n(NCN, "950"));
  });
});
