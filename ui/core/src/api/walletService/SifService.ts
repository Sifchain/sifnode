import {
  Secp256k1HdWallet,
  SigningCosmosClient,
  makeCosmoshubPath,
  CosmosClient,
  Account
} from "@cosmjs/launchpad"
import { SifWalletStore } from "src/store/wallet";

import { ADDR_PREFIX, API } from "../../constants"
import { Mnemonic, SifAddress } from "../../entities/Wallet" 
import { SifTransaction } from "../../entities/Transaction" 

// Warning This creates a client object used to *sign* TX
export async function cosmosSignin( mnemonic: Mnemonic ): 
  Promise<SigningCosmosClient> {
  try {
    if (!mnemonic) { throw "No mnemonic. Can't generate wallet."}
    const wallet = await Secp256k1HdWallet.fromMnemonic(
      mnemonic,
      makeCosmoshubPath(0),
      ADDR_PREFIX
    );
    const [{ address }] = await wallet.getAccounts();
    return new SigningCosmosClient(API, address, wallet);
  } catch (error) {
    throw error
  }
}

export async function getCosmosBalance( address: SifAddress ): 
  Promise<Account>  { 
    if (!address) throw "Address undefined. Fail"
    if (address.length !== 42) throw "Address not valid (length). Fail" // this is simple check, limited to default address type (check bech32)
    const client = new CosmosClient(API)
    try {
      const account = await client.getAccount(address)
      if (!account) throw "No Address found on chain"
      return account
    } catch (error) {
      throw error
    }
  }

  export async function sendSifToken(
    sifWallet: SifWalletStore, 
    sifTransaction: any
  ): Promise<any> {
    if (!sifWallet.client) throw "No signed in client. Sign in with mnemonic."
    if (!sifTransaction) throw "No user input data. Define who, what, and for how much."

    const from_address = sifWallet.client.senderAddress;
    const msg = {
      type: "cosmos-sdk/MsgSend",
      value: {
        amount: [
          {
            sifTransaction.amount,
            sifTransaction.denom,
          },
        ],
        from_address,
        sifTransaction.to_address,
      },
    };
    const fee = {
      amount: coins(0, denom),
      gas: "200000",
    };
    return await state.client.signAndPost([msg], fee, memo);
  }