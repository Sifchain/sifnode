import connectToWallet from "./connectToWallet";
import { createStore } from "../store";

import { mnemonicIsValid, generateMnemonicAction, signInCosmosWallet  } from "./CWalletActions"

// before all, instantiate test chain
// cd ../
// make install
// sifgen node create 


// starport serve (get mnemonic)
// get mnemonic
// "get balance"

const badMnemonic = "Ever have that feeling where you’re not sure if you’re awake or dreaming?"

const mnemonicShadowfiend = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

const mnemonicAkasha = "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard"


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
      const a = await signInCosmosWallet(mnemonicShadowfiend)
      console.log(a)
      // expect(signInCosmosWallet(mnemonicShadowfiend)).

    });


});