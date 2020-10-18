
import { 
  mnemonicIsValid, 
  generateMnemonicAction, 
  signInCosmosWallet, 
  getCosmosBalanceAction  
} from "./CWalletActions"
import { SigningCosmosClient } from "@cosmjs/launchpad";


const badMnemonic = "Ever have that feeling where you’re not sure if you’re awake or dreaming?"

const mnemonicShadowfiend = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

const mnemonicAkasha = "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard"

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
const badAddress = "sif1xjfzdf02kyg9t772j427h9vaeql5c938k3h5e5"

describe("connectToWallet", () => {

    it("should check if mnemonic is valid", () => {
      expect(mnemonicIsValid(badMnemonic)).toEqual(false)
      expect(mnemonicIsValid(mnemonicShadowfiend)).toEqual(true)
    })

    it("should check generate valid mnemonic", () => {
      expect(mnemonicIsValid(generateMnemonicAction())).toEqual(true)
    })

    it("should use mnemeonic to sign into cosmos wallet", async () => {
      // get account mnemonic on 
      expect(await signInCosmosWallet(mnemonicShadowfiend)).toMatchObject({senderAddress: account.address})
    })

    it ("should throw if no address", async () => {
      //https://github.com/facebook/jest/issues/1700
      try {
        await getCosmosBalanceAction("")
      } catch (e) {
        expect(e).toEqual("Address undefined. Fail")
      }
    })

    it ("should throw if no address", async () => {
      try {
        await getCosmosBalanceAction(badAddress)
      } catch (e) {
        expect(e).toEqual("No Address found on chain")
      }
    })

    it("should use address to get balance", async () => {
      expect(await getCosmosBalanceAction(account.address)).toEqual({
        address: account.address,
        accountNumber: account.accountNumber,
        sequence: account.sequence,
        pubkey: account.pubkey,
        balance: [
          { denom: "nametoken", amount: "1000" },
        ],
      });
    });
    
});