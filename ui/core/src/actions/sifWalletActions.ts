import { SifTransaction } from "../entities";
import { validateMnemonic } from "bip39";
import { Mnemonic, SifAddress } from "../entities/Wallet";
import { ActionContext } from ".";

export default ({ api }: ActionContext<"SifService">) => {
  const actions = {
    async getCosmosAction(address: SifAddress) {
      // check if sif prefix
      return await api.SifService.getBalance(address);
    },

    async signInCosmosWallet(mnemonic: Mnemonic): Promise<string> {
      if (!mnemonic) throw "Mnemonic must be defined";
      if (!validateMnemonic(mnemonic)) throw "Invalid Mnemonic. Not sent.";
      return await api.SifService.setPhrase(mnemonic);
    },

    async sendTransaction(sifTransaction: SifTransaction) {
      return await api.SifService.transfer(sifTransaction);
    },
  };

  return actions;
};
