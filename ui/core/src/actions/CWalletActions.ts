import { validateMnemonic, generateMnemonic } from "bip39"
// import {Ref, ComputedRef} from "@vue/reactivity"
// import { ActionContext } from "..";
import { Mnemonic } from "../entities/Wallet"
// import { CWalletStore, ICWalletStore } from "../store/wallet"
import { cosmosSignin } from "../api/walletService/SifService"

import {
  Account,
  SigningCosmosClient
} from "@cosmjs/launchpad";

export function generateMnemonicAction(): Mnemonic {
  return generateMnemonic()
}

export async function signInCosmosWallet(mnemonic: Mnemonic): 
  Promise<{account: Account, client: SigningCosmosClient}> {
  if(!mnemonicIsValid(mnemonic)) throw "Invalid Mnemonic. Not sent."
  return await cosmosSignin(mnemonic)
}

export function mnemonicIsValid(mnemonic: Mnemonic): Boolean {
  if (!mnemonic) throw "Mnemonic must be defined"
  return validateMnemonic(mnemonic)
}
