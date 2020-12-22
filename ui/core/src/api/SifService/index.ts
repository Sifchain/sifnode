import {
  coins,
  isBroadcastTxFailure,
  makeCosmoshubPath,
  Msg,
  Secp256k1HdWallet,
} from "@cosmjs/launchpad";
import { reactive } from "@vue/reactivity";
import {
  Address,
  Asset,
  AssetAmount,
  Coin,
  Network,
  TxParams,
} from "../../entities";
import { Mnemonic } from "../../entities/Wallet";
import { IWalletService } from "../IWalletService";
import { SifClient } from "../utils/SifClient";
import { ensureSifAddress } from "./utils";

export type SifServiceContext = {
  sifAddrPrefix: string;
  sifApiUrl: string;
  assets: Asset[];
};

type ISifService = IWalletService & { getSupportedTokens: () => Asset[] };

/**
 * Constructor for SifService
 *
 * SifService handles communication between our ui core Domain and the SifNode blockchain
 */
export default function createSifService({
  sifAddrPrefix,
  sifApiUrl,
  assets,
}: SifServiceContext): ISifService {
  const {} = sifAddrPrefix;

  // Reactive state for communicating state changes
  const state: {
    connected: boolean;
    address: Address;
    accounts: Address[];
    balances: AssetAmount[];
    log: string; // latest transaction hash
  } = reactive({
    connected: false,
    accounts: [],
    address: "",
    balances: [],
    log: "unset",
  });

  let client: SifClient | null = null;

  const supportedTokens = assets.filter(
    (asset) => asset.network === Network.SIFCHAIN
  );

  return {
    /**
     * getState returns the service's reactive state to be listened to by consuming clients.
     */
    getState() {
      return state;
    },

    getSupportedTokens() {
      return supportedTokens;
    },

    async connect() {},

    async disconnect() {},

    isConnected() {
      return state.connected;
    },

    async setPhrase(mnemonic: Mnemonic): Promise<Address> {
      try {
        if (!mnemonic) {
          throw "No mnemonic. Can't generate wallet.";
        }

        const wallet = await Secp256k1HdWallet.fromMnemonic(
          mnemonic,
          makeCosmoshubPath(0),
          sifAddrPrefix
        );

        state.accounts = (await wallet.getAccounts()).map(
          ({ address }) => address
        );

        [state.address] = state.accounts;

        client = new SifClient(sifApiUrl, state.address, wallet);

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

    async getBalance(address?: Address): Promise<AssetAmount[]> {
      if (!client) throw "No client. Please sign in.";
      if (!address) throw "Address undefined. Fail";

      ensureSifAddress(address);

      try {
        const account = await client.getAccount(address);

        if (!account) throw "No Address found on chain";

        state.balances = account.balance.map(({ amount, denom }) => {
          // HACK: Following should be a lookup of tokens loaded from genesis somehow
          const asset = Coin({
            symbol: denom,
            decimals: 0,
            name: denom,
            network: Network.SIFCHAIN,
          });
          return AssetAmount(asset, amount);
        });
        return state.balances;
      } catch (error) {
        throw error;
      }
    },

    async transfer(params: TxParams): Promise<any> {
      if (!client) throw "No client. Please sign in.";
      if (!params.asset) throw "No asset.";
      try {
        // https://github.com/tendermint/vue/blob/develop/src/store/cosmos.js#L91
        const msg = {
          type: "cosmos-sdk/MsgSend",
          value: {
            amount: [
              {
                amount: params.amount.toString(),
                denom: params.asset.symbol,
              },
            ],
            from_address: client.senderAddress,
            to_address: params.recipient,
          },
        };

        const fee = {
          amount: coins(0, params.asset.symbol),
          gas: "200000", // need gas fee for tx to work - see genesis file
        };

        const txHash = await client.signAndBroadcast([msg], fee, params.memo);

        this.getBalance(state.address);

        return txHash;
      } catch (err) {
        console.error(err);
      }
    },

    async signAndBroadcast(msg: Msg | Msg[], memo?: string) {
      if (!client) throw "No client. Please sign in.";
      try {
        const fee = {
          amount: coins(0, "rowan"),
          gas: "200000", // need gas fee for tx to work - see genesis file
        };

        const msgArr = Array.isArray(msg) ? msg : [msg];

        const txHash = await client.signAndBroadcast(msgArr, fee, memo);

        if (isBroadcastTxFailure(txHash)) {
          console.log(txHash.rawLog);
          throw new Error(txHash.rawLog);
        }
        this.getBalance(state.address);

        return txHash;
      } catch (err) {
        console.error(err);
      }
    },
  };
}
