import { coins, makeCosmoshubPath, Secp256k1HdWallet } from "@cosmjs/launchpad";
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
import { SifClient } from "./SifClient";
import { ensureSifAddress } from "./utils";

export type SifServiceContext = {
  sifAddrPrefix: string;
  sifApiUrl: string;
};

type IClpService = {
  swap: (params: { receivedAsset: Asset; sentAmount: AssetAmount }) => any;
  addLiquidity: (params: {
    nativeAssetAmount: AssetAmount;
    externalAssetAmount: AssetAmount;
  }) => any;
  createPool: (params: {
    nativeAssetAmount: AssetAmount;
    externalAssetAmount: AssetAmount;
  }) => any;
};

/**
 * Constructor for SifService
 *
 * SifService handles communication between our ui core Domain and the SifNode blockchain
 */
export default function createSifService({
  sifAddrPrefix,
  sifApiUrl,
}: SifServiceContext): IWalletService & IClpService {
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

  return {
    /**
     * getState returns the service's reactive state to be listened to by consuming clients.
     */
    getState() {
      return state;
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
    },

    async addLiquidity(params: {
      nativeAssetAmount: AssetAmount;
      externalAssetAmount: AssetAmount;
    }) {
      if (!client) throw "No client. Please sign in.";
      const response = await client.addLiquidity({
        base_req: { chain_id: "sifnode", from: state.address },
        external_asset: {
          source_chain: params.externalAssetAmount.asset.network as string,
          symbol: params.externalAssetAmount.asset.symbol,
          ticker: params.externalAssetAmount.asset.symbol,
        },
        external_asset_amount: params.externalAssetAmount.toFixed(0),
        native_asset_amount: params.nativeAssetAmount.toFixed(0),
        signer: state.address,
      });
      const fee = {
        amount: coins(0, params.externalAssetAmount.asset.symbol),
        gas: "200000", // need gas fee for tx to work - see genesis file
      };
      // alert(JSON.stringify(response));
      const txHash = await client.signAndBroadcast(response.value.msg, fee);

      this.getBalance(state.address);

      return txHash;
    },

    async createPool(params: {
      nativeAssetAmount: AssetAmount;
      externalAssetAmount: AssetAmount;
    }) {
      if (!client) throw "No client. Please sign in.";
      const response = await client.createPool({
        base_req: { chain_id: "sifnode", from: state.address },
        external_asset: {
          source_chain: params.externalAssetAmount.asset.network as string,
          symbol: params.externalAssetAmount.asset.symbol,
          ticker: params.externalAssetAmount.asset.symbol,
        },
        external_asset_amount: params.externalAssetAmount.toFixed(0),
        native_asset_amount: params.nativeAssetAmount.toFixed(0),
        signer: state.address,
      });
      const fee = {
        amount: coins(0, params.externalAssetAmount.asset.symbol),
        gas: "200000", // need gas fee for tx to work - see genesis file
      };
      // alert(JSON.stringify(response));
      const txHash = await client.signAndBroadcast(response.value.msg, fee);

      this.getBalance(state.address);

      return txHash;
    },

    async swap(params: { sentAmount: AssetAmount; receivedAsset: Asset }) {
      if (!client) throw "No client. Please sign in.";

      // Validate params

      const response = await client.swap({
        base_req: { chain_id: "sifchain", from: state.address },
        received_asset: {
          source_chain: params.receivedAsset.network as string,
          symbol: params.receivedAsset.symbol,
          ticker: params.receivedAsset.symbol,
        },
        sent_amount: params.sentAmount.numerator.toString(),
        sent_asset: {
          source_chain: params.sentAmount.asset.network as string,
          symbol: params.sentAmount.asset.symbol,
          ticker: params.sentAmount.asset.symbol,
        },
        signer: state.address,
      });
      const fee = {
        amount: coins(0, params.sentAmount.asset.symbol),
        gas: "200000", // need gas fee for tx to work - see genesis file
      };

      const txHash = await client.signAndBroadcast(response.value.msg, fee);

      this.getBalance(state.address);

      return txHash;
    },
  };
}
