import {
  coins,
  isBroadcastTxFailure,
  makeCosmoshubPath,
  Msg,
  Secp256k1HdWallet,
} from "@cosmjs/launchpad";
import { reactive } from "@vue/reactivity";
import { debounce } from "lodash";
import { Address, Asset, AssetAmount, Network, TxParams } from "../../entities";

import { Mnemonic } from "../../entities/Wallet";
import { IWalletService } from "../IWalletService";
import { SifClient, SifUnSignedClient } from "../utils/SifClient";
import { ensureSifAddress } from "./utils";
import getKeplrProvider from "./getKeplrProvider";
import { KeplrChainConfig } from "../../utils/parseConfig";

export type SifServiceContext = {
  sifAddrPrefix: string;
  sifApiUrl: string;
  sifWsUrl: string;
  // todo fix
  keplrChainConfig: KeplrChainConfig;
  assets: Asset[];
};
type HandlerFn<T> = (a: T) => void;
export type ISifService = IWalletService & {
  getSupportedTokens: () => Asset[];
  onSocketError: (handler: HandlerFn<any>) => void;
  onTx: (handler: HandlerFn<any>) => void;
  onNewBlock: (handler: HandlerFn<any>) => void;
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
  keplrChainConfig,
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

  const keplrProvider = getKeplrProvider();
  let client: SifClient | null = null;
  let closeUpdateListener = () => {};

  const unSignedClient = new SifUnSignedClient(sifApiUrl, sifWsUrl);

  const supportedTokens = assets.filter(
    asset => asset.network === Network.SIFCHAIN
  );

  // TODO: deletion ?
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
      console.log("triggerUpdate" + client?.senderAddress);
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
      console.log("connect service", keplrChainConfig, keplrProvider);
      if (!keplrProvider) {
        throw {
          message: "Keplr Not Found",
          detail: "Check if extension enabled for this URL",
        };
      }
      // open extension
      if (keplrProvider.experimentalSuggestChain) {
        try {
          await keplrProvider.experimentalSuggestChain(keplrChainConfig);
          await keplrProvider.enable(keplrChainConfig.chainId);

          const offlineSigner = keplrProvider.getOfflineSigner(
            keplrChainConfig.chainId
          );
          // https://github.com/chainapsis/keplr-extension/blob/960e50f1d9360d21d6935b974a0cb8b57c27d9d9/src/content-scripts/inject/cosmjs-offline-signer.ts
          const accounts = await offlineSigner.getAccounts();

          // get balances
          const address = accounts.length > 0 ? accounts[0].address : "";

          if (!address) {
            throw "No address on sif account";
          }

          client = new SifClient(sifApiUrl, address, offlineSigner, sifWsUrl);
          triggerUpdate();
          closeUpdateListener = client.getUnsignedClient().onNewBlock(() => {
            triggerUpdate();
          });
        } catch (error) {
          console.log(error);
          throw { message: "Failed to Suggest Chain" };
        }
      } else {
        throw { message: "Keplr Outdated", detail: "Need at least 0.6.4" };
      }
    },

    async disconnect() {
      // disconnect from Keplr
      await this.purgeClient();
    },

    isConnected() {
      return state.connected;
    },

    onSocketError(handler: HandlerFn<any>) {
      unSignedClient.onSocketError(handler);
    },

    onTx(handler: HandlerFn<any>) {
      unSignedClient.onTx(handler);
    },

    onNewBlock(handler: HandlerFn<any>) {
      unSignedClient.onNewBlock(handler);
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

    async purgeClient() {
      client = null;
      await triggerUpdate();
      closeUpdateListener();
    },

    async getBalance(address?: Address): Promise<AssetAmount[]> {
      if (!client) throw "No client. Please sign in.";
      if (!address) throw "Address undefined. Fail";

      ensureSifAddress(address);

      try {
        const account = await client.getAccount(address);
        console.log(account);
        if (!account) throw "No Address found on chain";

        const supportedTokenSymbols = supportedTokens.map(s => s.symbol);
        const balances = account.balance
          .filter(balance => supportedTokenSymbols.includes(balance.denom))
          .map(({ amount, denom }) => {
            const asset = supportedTokens.find(
              token => token.symbol === denom
            )!; // will be found because of filter above

            return AssetAmount(asset, amount, { inBaseUnit: true });
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
        console.log("msgArr", msgArr);

        const txHash = await client.signAndBroadcast(msgArr, fee, memo);

        if (isBroadcastTxFailure(txHash)) {
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
