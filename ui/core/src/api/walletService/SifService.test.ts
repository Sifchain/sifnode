import * as sifService from "./SifService"

const badMnemonic = "Ever have that feeling where you’re not sure if you’re awake or dreaming?"

const mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

const account =  {
  address: 'sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd',
  balance: [ { denom: 'nametoken', amount: '1000' } ],
  pubkey: {
    type: 'tendermint/PubKeySecp256k1',
    value: 'AvUEsFHbsr40nTSmWh7CWYRZHGwf4cpRLtJlaRO4VAoq'
  },
  accountNumber: 2,
  sequence: 1
}

// This is redundant. CWalletActions and SifService should be combined imo
describe("sifService", () => {
    it("should use mnemeonic to sign into cosmos wallet", async () => {
      // catch Error["Bad mnemonic"]??
      test.todo("more tests on bad mnemonic ")

      try {
        await sifService.cosmosSignin("")
      } catch (error) {
        expect(error).toEqual("No mnemonic. Can't generate wallet.")
      }

      expect(await sifService.cosmosSignin(mnemonic))
        .toMatchObject({senderAddress: account.address})
    })

    it("should get cosmos balance", async () => {
      try {
        await sifService.getCosmosBalance("")
      } catch (error) {
        expect(error).toEqual("Address undefined. Fail")
      }
      try {
        await sifService.getCosmosBalance("asdfasdsdf")
      } catch (error) {
        expect(error).toEqual("Address not valid (length). Fail")
      }

      expect(await sifService.getCosmosBalance(account.address))
        .toMatchObject(account)
    })
});