import {
  Secp256k1HdWallet,
  SigningCosmosClient,
  makeCosmoshubPath,
  coins,
  CosmosClient,
  Account
} from "@cosmjs/launchpad"

import { ADDR_PREFIX, API } from "../../constants"
import { Mnemonic, SifAddress } from "../../entities/Wallet" 

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
  Promise<Account | Error>  { 
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


// SEND() -- for TX
// async function validateAddressOnChain() {

//   if (this.valid.to_address && this.valid.amount && !this.inFlight) {
//     const payload = {
//       amount: this.amount,
//       denom: this.denom,
//       to_address: this.to_address,
//       memo: this.memo,
//     };
//     this.txResult = "";
//     this.inFlight = true;
//     this.txResult = await this.$store.dispatch("cosmos/tokenSend", payload);
//     if (!this.txResult.code) {
//       this.amount = "";
//       this.to_address = "";
//       this.memo = "";
//     }
//     this.inFlight = false;
//     await this.$store.dispatch("cosmos/bankBalancesGet");
  
// }
