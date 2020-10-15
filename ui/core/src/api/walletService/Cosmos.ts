import {  ICWalletStore, CWalletStore } from "../../store/wallet"

import axios from "axios";
import {
  Secp256k1HdWallet,
  SigningCosmosClient,
  makeCosmoshubPath,
  coins,
  Account
} from "@cosmjs/launchpad";

const API = "http://localhost:1317";
const ADDR_PREFIX = "sif";

export async function cosmosSignin( mnemonic: ICWalletStore["mnemonic"] ) {

  try {
    if (!mnemonic) { throw "No mnemonic. Can't generate wallet."}
    const wallet = await Secp256k1HdWallet.fromMnemonic(
      mnemonic,
      makeCosmoshubPath(0),
      ADDR_PREFIX
    );
    // localStorage.setItem("mnemonic", mnemonic);
    const [{ address }] = await wallet.getAccounts();
    const url = `${API}/auth/accounts/${address}`;
    const acc = (await axios.get(url)).data;
    const account: Account = acc.result.value;
    console.log(account)
    CWalletStore.account = account
    // commit("set", { key: "account", value: account });
    const client = new SigningCosmosClient(API, address, wallet);
    console.log(client)
    CWalletStore.client = client

    // commit("set", { key: "client", value: client });
    // // dispatch("delegationsFetch");
    // // dispatch("transfersIncomingFetch");
    // // dispatch("transfersOutgoingFetch");
    // try {
    //   await dispatch("bankBalancesGet");
    // } catch {
    //   console.log("Error in getting a bank balance.");
    // }
    // resolve(account);
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
