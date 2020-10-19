import {
  Secp256k1HdWallet,
  SigningCosmosClient,
  makeCosmoshubPath,
  CosmosClient,
  coins,
  Coin,
} from "@cosmjs/launchpad";

import { ADDR_PREFIX, API } from "../../constants";
import { Mnemonic, SifAddress } from "../../entities/Wallet";
import { SifTransaction } from "../../entities/Transaction";
import { Address } from "../../entities";
import { reactive } from "@vue/reactivity";

export type SifServiceContext = {};

export default function createSifService(_context: SifServiceContext) {
  const state: {
    connected: boolean;
    address: Address;
    accounts: Address[];
    balances: Coin[];
    log: string; // latest transaction hash
  } = reactive({
    connected: false,
    accounts: [],
    address: "",
    balances: [],
    log: "unset",
  });

  let client: SigningCosmosClient | null = null;

  return {
    // Return reactive state
    getState() {
      return state;
    },

    async setPhrase(mnemonic: Mnemonic): Promise<Address> {
      try {
        if (!mnemonic) {
          throw "No mnemonic. Can't generate wallet.";
        }

        const wallet = await Secp256k1HdWallet.fromMnemonic(
          mnemonic,
          makeCosmoshubPath(0),
          ADDR_PREFIX
        );

        state.accounts = (await wallet.getAccounts()).map(
          ({ address }) => address
        );

        [state.address] = state.accounts;

        client = new SigningCosmosClient(API, state.address, wallet);

        state.log = "signed in";
        state.connected = true;
        this.getBalance(state.address);
        return state.address;
      } catch (error) {
        throw error;
      }
    },

    purgeClient() {
      state.address = "";
      state.connected = false;
      state.balances = [];
      state.accounts = [];
      state.log = "";
    },

    async getBalance(address?: SifAddress): Promise<readonly Coin[]> {
      if (!address) throw "Address undefined. Fail";

      if (address.length !== 42) throw "Address not valid (length). Fail"; // this is simple check, limited to default address type (check bech32)
      // TODO: add invariant address starts with "sif" (double check this is correct)

      const client = new CosmosClient(API);

      try {
        const account = await client.getAccount(address);

        if (!account) throw "No Address found on chain";

        state.balances = account.balance as Coin[];
        return account.balance;
      } catch (error) {
        throw error;
      }
    },

    async transfer(params: SifTransaction): Promise<any> {
      if (!client) throw "No signed in client. Sign in with mnemonic.";
      if (!params)
        throw "No user input data. Define who, what, and for how much.";
      // this seems like anti-pattern, with SifWallet.vue, "undefined" as culprit
      // but is alternative to define in vue with empty string?
      if (!params.asset) throw "No asset.";
      // https://github.com/tendermint/vue/blob/develop/src/store/cosmos.js#L91
      const msg = {
        type: "cosmos-sdk/MsgSend",
        value: {
          amount: [
            {
              amount: params.amount,
              denom: params.asset,
            },
          ],
          from_address: client.senderAddress,
          to_address: params.recipient,
        },
      };

      const fee = {
        amount: coins(0, params.asset),
        gas: "200000", // need gas fee for tx to work - see genesis file
      };

      const txHash = await client.signAndBroadcast([msg], fee, "cool");

      this.getBalance(state.address);

      return txHash;
    },
  };
}
