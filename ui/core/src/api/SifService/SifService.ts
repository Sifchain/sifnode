import {
  coins,
  isBroadcastTxFailure,
  makeCosmoshubPath,
  Msg,
  Secp256k1HdWallet,
} from "@cosmjs/launchpad";
import { reactive } from "@vue/reactivity";
import { debounce } from "lodash";
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
  sifWsUrl: string;
  assets: Asset[];
};

export type ISifService = IWalletService & {
  getSupportedTokens: () => Asset[];
};

/**
 * Constructor for SifService
 *
 * SifService handles communication between our ui core Domain and the SifNode blockchain
 */
export default function createSifService({
  sifAddrPrefix,
  sifApiUrl,
  sifWsUrl,
  assets,
}: SifServiceContext): ISifService {
  const {} = sifAddrPrefix;

  // Reactive state for communicating state changes
  // TODO this should be replaced with event handlers
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

  async function createSifClientFromMnemonic(mnemonic: string) {
    const wallet = await Secp256k1HdWallet.fromMnemonic(
      mnemonic,
      makeCosmoshubPath(0),
      sifAddrPrefix
    );
    const accounts = await wallet.getAccounts();

    const address = accounts.length > 0 ? accounts[0].address : "";

    if (!address) {
      throw new Error("No address on sif account");
    }

    return new SifClient(sifApiUrl, address, wallet, sifWsUrl);
  }

  const triggerUpdate = debounce(
    async () => {
      if (!client) {
        state.connected = false;
        state.address = "";
        state.balances = [];
        state.accounts = [];
        state.log = "";
        return;
      }

      state.connected = !!client;
      state.address = client.senderAddress;
      state.accounts = await client.getAccounts();
      state.balances = await instance.getBalance(client.senderAddress);
    },
    100,
    { leading: true }
  );

  const instance = {
    /**
     * getState returns the service's reactive state to be listened to by consuming clients.
     */
    getState() {
      return state;
    },

    getSupportedTokens() {
      return supportedTokens;
    },

    async connect() {
      // connect to Keplr
    },

    async disconnect() {
      // disconnect from Keplr
    },

    isConnected() {
      return state.connected;
    },

    async setPhrase(mnemonic: Mnemonic): Promise<Address> {
      try {
        if (!mnemonic) {
          throw "No mnemonic. Can't generate wallet.";
        }

        client = await createSifClientFromMnemonic(mnemonic);

        client.getUnsignedClient().onNewBlock(() => {
          triggerUpdate();
        });

        triggerUpdate();

        return client.senderAddress;
      } catch (error) {
        throw error;
      }
    },

    purgeClient() {
      client = null;
      triggerUpdate();
    },

    async getBalance(address?: Address): Promise<AssetAmount[]> {
      if (!client) throw "No client. Please sign in.";
      if (!address) throw "Address undefined. Fail";

      ensureSifAddress(address);

      try {
        const account = await client.getAccount(address);

        if (!account) throw "No Address found on chain";

        const balances = account.balance.map(({ amount, denom }) => {
          // HACK: Following should be a lookup of tokens loaded from genesis somehow
          const asset = Coin({
            symbol: denom,
            decimals: 0,
            name: denom,
            network: Network.SIFCHAIN,
          });
          return AssetAmount(asset, amount);
        });
        return balances;
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

        triggerUpdate();

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

        triggerUpdate();

        return txHash;
      } catch (err) {
        console.error(err);
      }
    },
  };
  return instance;
}
