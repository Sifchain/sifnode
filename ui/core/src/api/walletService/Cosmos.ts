import {  ICWalletStore } from "../../store/wallet"

import axios from "axios";
import {
  Secp256k1HdWallet,
  SigningCosmosClient,
  makeCosmoshubPath,
  coins,
} from "@cosmjs/launchpad";

const API = "http://localhost:1317";
const ADDR_PREFIX = process.env.VUE_APP_ADDRESS_PREFIX || "cosmos";

export async function cosmosSignin( mnemonic: ICWalletStore["mnemonic"] ) {
  if (!mnemonic) { throw "No mnemonic. Can't generate wallet."}
  return new Promise(async (resolve, reject) => {
    const wallet = await Secp256k1HdWallet.fromMnemonic(
      mnemonic,
      makeCosmoshubPath(0),
      ADDR_PREFIX
    );
    // localStorage.setItem("mnemonic", mnemonic);
    const [{ address }] = await wallet.getAccounts();
    const url = `${API}/auth/accounts/${address}`;
    const acc = (await axios.get(url)).data;
    const account = acc.result.value;
    console.log(account)
    // commit("set", { key: "account", value: account });
    const client = new SigningCosmosClient(API, address, wallet);
    console.log(client)
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
  });
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
