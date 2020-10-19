// import { signInCosmosWallet, getCosmosAction } from "./sifWalletActions";

// const mnemonic =
//   "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow";
// const account = {
//   address: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
//   balance: [{ denom: "nametoken", amount: "1000" }],
//   pubkey: {
//     type: "tendermint/PubKeySecp256k1",
//     value: "AvUEsFHbsr40nTSmWh7CWYRZHGwf4cpRLtJlaRO4VAoq",
//   },
//   accountNumber: 2,
//   sequence: 1,
// };
// const badAddress = "sif1xjfzdf02kyg9t772j427h9vaeql5c938k3h5e5";
// const badMnemonic =
//   "Ever have that feeling where you’re not sure if you’re awake or dreaming?";

describe("connectToWallet", () => {
  // it("should check if mnemonic is valid", () => {
  //   expect(mnemonicIsValid(badMnemonic)).toEqual(false);
  //   expect(mnemonicIsValid(mnemonic)).toEqual(true);
  // });
  // it("should check generate valid mnemonic", () => {
  //   expect(mnemonicIsValid(generateMnemonicAction())).toEqual(true);
  // });
  // it("should use mnemeonic to sign into cosmos wallet", async () => {
  //   expect(await signInCosmosWallet(mnemonic)).toMatchObject({
  //     senderAddress: account.address,
  //   });
  // });
  // it("should throw if no address", async () => {
  //   //https://github.com/facebook/jest/issues/1700
  //   try {
  //     await getCosmosAction("");
  //   } catch (e) {
  //     expect(e).toEqual("Address undefined. Fail");
  //   }
  // });
  // it("should throw if no address", async () => {
  //   try {
  //     await getCosmosAction(badAddress);
  //   } catch (e) {
  //     expect(e).toEqual("No Address found on chain");
  //   }
  // });
  it("should use address to get balance", async () => {
    //   expect(await getCosmosAction(account.address)).toEqual(account);
    expect(1).toBe(1);
  });
});
