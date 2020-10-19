import { Account, SigningCosmosClient } from "@cosmjs/launchpad";
import { SifTransaction } from "../entities";
import { SifWalletStore } from "../store/wallet";
import { validateMnemonic, generateMnemonic } from "bip39";
import {
  cosmosSignin,
  getCosmosBalance,
  signAndBroadcast,
} from "../api/SifService";
import { Mnemonic, SifAddress } from "../entities/Wallet";

export async function getCosmosBalanceAction(address: SifAddress) {
  // check if sif prefix
  return await getCosmosBalance(address);
}
export async function signInCosmosWallet(
  mnemonic: Mnemonic
): Promise<SigningCosmosClient> {
  if (!mnemonic) throw "Mnemonic must be defined";
  if (!mnemonicIsValid(mnemonic)) throw "Invalid Mnemonic. Not sent.";
  return await cosmosSignin(mnemonic);
}

export function mnemonicIsValid(mnemonic: Mnemonic): Boolean {
  return validateMnemonic(mnemonic);
}

export function generateMnemonicAction(): Mnemonic {
  return generateMnemonic();
}

export async function sendTransaction(
  sifWalletClient: SifWalletStore["client"],
  sifTransaction: SifTransaction
) {
  return await signAndBroadcast(sifWalletClient, sifTransaction);
}
