import { coins, isBroadcastTxFailure, Msg } from "@cosmjs/launchpad";
import { reactive } from "@vue/reactivity";
import { debounce } from "lodash";
import {
  Address,
  Asset,
  AssetAmount,
  Network,
  TransactionStatus,
  TxParams,
} from "../../entities";

import { Mnemonic } from "../../entities";

import { SifClient, SifUnSignedClient } from "../utils/SifClient";
import { ensureSifAddress } from "./utils";
import getKeplrProvider from "./getKeplrProvider";
import { KeplrChainConfig } from "../../utils/parseConfig";
import { parseTxFailure } from "./parseTxFailure";

export type SifServiceContext = {
  sifAddrPrefix: string;
  sifApiUrl: string;
  sifWsUrl: string;
  sifRpcUrl: string;
  keplrChainConfig: KeplrChainConfig;
  assets: Asset[];
};
type HandlerFn<T> = (a: T) => void;

/**
 * Constructor for SifService
 *
 * SifService handles communication between our ui core Domain and the SifNode blockchain
 */
export default function createSifService({
  sifAddrPrefix,
  sifApiUrl,
  sifWsUrl,
  sifRpcUrl,
  keplrChainConfig,
  assets,
}: SifServiceContext) {
  const {} = sifAddrPrefix;

  const initState = {
    connected: false,
    accounts: [],
    address: "",
    balances: [],
    log: "unset",
  };

  const state: {
    connected: boolean;
    address: Address;
    accounts: Address[];
    balances: AssetAmount[];
    log: string; // latest transaction hash
  } = reactive(initState);

  const keplrProviderPromise = getKeplrProvider();
  let keplrProvider: any;
  let offlineSigner: any;
  let client: SifClient | null = null;
  let polling: any;

  const unSignedClient = new SifUnSignedClient(sifApiUrl, sifWsUrl, sifRpcUrl);

  const supportedTokens = assets.filter(
    (asset) => asset.network === Network.SIFCHAIN
  );

  const triggerUpdate = debounce(
    async () => {
      try {
        if (!polling) {
          polling = setInterval(() => {
            triggerUpdate();
          }, 2000);
        }
        await instance.setClient();
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
      } catch (e) {
        state.connected = false;
        state.address = "";
        state.balances = [];
        state.accounts = [];
        state.log = "";
        if (polling) {
          clearInterval(polling);
          polling = null;
        }
      }
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

    async setClient() {
      if (!offlineSigner) {
        client = null;
        return;
      }
      const accounts = await offlineSigner.getAccounts();
      const address = accounts.length > 0 ? accounts[0].address : "";
      if (!address) {
        throw "No address on sif account";
      }
      client = new SifClient(
        sifApiUrl,
        address,
        offlineSigner,
        sifWsUrl,
        sifRpcUrl
      );
    },

    async initProvider() {
      try {
        keplrProvider = await keplrProviderPromise;
        if (!keplrProvider) {
          return;
        }
        offlineSigner = keplrProvider.getOfflineSigner(
          keplrChainConfig.chainId
        );
        await instance.setClient();
        triggerUpdate();
      } catch (e) {
        console.log("initProvider", e);
      }
    },

    async connect() {
      if (!keplrProvider) {
        keplrProvider = await keplrProviderPromise;
      }
      // open extension
      if (keplrProvider.experimentalSuggestChain) {
        try {
          await keplrProvider.experimentalSuggestChain(keplrChainConfig);
          await keplrProvider.enable(keplrChainConfig.chainId);
          triggerUpdate();
        } catch (error) {
          console.log(error);
          throw { message: "Failed to Suggest Chain" };
        }
      } else {
        throw {
          message: "Keplr Outdated",
          detail: { type: "info", message: "Need at least 0.6.4" },
        };
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
      // We currently delegate auth to Keplr so this is irrelevant
      return "";
    },

    async purgeClient() {
      // We currently delegate auth to Keplr so this is irrelevant
    },

    async getBalance(address?: Address, asset?: Asset): Promise<AssetAmount[]> {
      if (!client) {
        throw "No client. Please sign in.";
      }
      if (!address) {
        throw "Address undefined. Fail";
      }

      ensureSifAddress(address);

      try {
        const account = await client.getAccount(address);
        if (!account) {
          throw "No Address found on chain";
        } // todo handle this better
        const supportedTokenSymbols = supportedTokens.map((s) => s.symbol);
        return account.balance
          .filter((balance) => supportedTokenSymbols.includes(balance.denom))
          .map(({ amount, denom }) => {
            const asset = supportedTokens.find(
              (token) => token.symbol === denom
            )!; // will be found because of filter above
            return AssetAmount(asset, amount, { inBaseUnit: true });
          })
          .filter((balance) => {
            // If an aseet is supplied filter for it
            if (!asset) {
              return true;
            }
            return balance.asset.symbol === asset.symbol;
          });
      } catch (error) {
        throw error;
      }
    },

    async transfer(params: TxParams): Promise<any> {
      if (!client) {
        throw "No client. Please sign in.";
      }
      if (!params.asset) {
        throw "No asset.";
      }
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
          gas: "300000", // TODO - see if "auto" setting
        };

        const txHash = await client.signAndBroadcast([msg], fee, params.memo);

        triggerUpdate();

        return txHash;
      } catch (err) {
        console.error(err);
      }
    },

    async signAndBroadcast(
      msg: Msg | Msg[],
      memo?: string
    ): Promise<TransactionStatus> {
      if (!client) {
        throw "No client. Please sign in.";
      }
      try {
        const fee = {
          amount: coins(0, "rowan"),
          gas: "300000", // TODO - see if "auto" setting
        };

        const msgArr = Array.isArray(msg) ? msg : [msg];

        const result = await client.signAndBroadcast(msgArr, fee, memo);

        if (isBroadcastTxFailure(result)) {
          /* istanbul ignore next */ // TODO: fix coverage
          return parseTxFailure(result);
        }

        triggerUpdate();

        return {
          hash: result.transactionHash,
          memo,
          state: "accepted",
        };
      } catch (err) {
        console.log("signAndBroadcast ERROR", err);
        return parseTxFailure({ transactionHash: "", rawLog: err.message });
      }
    },
  };

  instance.initProvider();

  return instance;
}
