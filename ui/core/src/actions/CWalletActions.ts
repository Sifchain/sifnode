import { ActionContext } from "..";


import {Ref, ComputedRef} from "@vue/reactivity"
import { CWalletStore, ICWalletStore } from "../store/wallet"
import { cosmosSignin } from "../api/walletService/Cosmos"


// howto not duplicate mnemonic: ICWalletStore["mnemonic"]
export async function connectToWallet(mnemonic: ICWalletStore["mnemonic"]) {
  // set mnemonic in store (definitely won't do IRL)
  CWalletStore.mnemonic = mnemonic

  if(!mnemonicIsValid(mnemonic)) throw "Invalid Mnemonic. Not sent."
  // "sign in" on chain
  await cosmosSignin(mnemonic)
}

import { validateMnemonic } from "bip39"

function mnemonicIsValid(mnemonic:ICWalletStore["mnemonic"]): Boolean {
  if (!mnemonic) {throw "Mnemonic must be defined"}
  return validateMnemonic(mnemonic)
}
