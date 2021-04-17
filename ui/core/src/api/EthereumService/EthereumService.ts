import { reactive } from "@vue/reactivity";
import Web3 from "web3";
import { provider, WebsocketProvider } from "web3-core";
import { IWalletService } from "../IWalletService";
import { debounce } from "lodash";
import {
  TxHash,
  TxParams,
  Asset,
  AssetAmount,
  Network,
  IAssetAmount,
} from "../../entities";
import {
  EIPProvider,
  getEtheriumBalance,
  getTokenBalance,
  isEventEmittingProvider,
  isMetaMaskInpageProvider,
  transferAsset,
} from "./utils/ethereumUtils";

import { Msg } from "@cosmjs/launchpad";
import { EventEmitter2 } from "eventemitter2";

type Address = string;
type Balances = IAssetAmount[];
type PossibleProvider = provider | EIPProvider;

export type EthereumServiceContext = {
  getWeb3Provider: () => Promise<provider>;
  assets: Asset[];
};

const initState = {
  connected: false,
  accounts: [],
  balances: [],
  address: "",
  log: "unset",
};

// TODO: Refactor to be Module pattern with constructor function ie. `EthereumService()`
// TODO: Refactor all blockchain services to be RxJS

const PROVIDER_NOT_FOUND_EVENT = "PROVIDER_NOT_FOUND_EVENT";

export class EthereumService implements IWalletService {
  private web3: Web3 | null = null;
  private supportedTokens: Asset[] = [];
  private blockSubscription: any;
  private provider: PossibleProvider | undefined;
  private providerPromise: Promise<PossibleProvider>;
  private emitter: EventEmitter2;

  // This is shared reactive state
  private state: {
    connected: boolean;
    address: Address;
    accounts: Address[];
    balances: IAssetAmount[];
    log: string;
  };

  constructor(
    getWeb3Provider: () => Promise<PossibleProvider>,
    assets: Asset[],
  ) {
    this.state = reactive({ ...initState });
    this.supportedTokens = assets.filter((t) => t.network === Network.ETHEREUM);
    this.providerPromise = getWeb3Provider();
    this.emitter = new EventEmitter2();
    this.providerPromise
      .then((provider) => {
        // Provider not found
        if (!provider) {
          this.provider = null;
          this.emitter.emit(PROVIDER_NOT_FOUND_EVENT);
          return;
        }

        if (isEventEmittingProvider(provider)) {
          provider.on("chainChanged", () => window.location.reload());
          provider.on("accountsChanged", () => this.updateData());
        }

        this.web3 = new Web3(provider as provider);
        this.provider = provider;
        this.addWeb3Subscription();
        this.updateData();
      })
      .catch((error) => {
        console.log("error", error);
      });
  }

  async getChainId() {
    if (isMetaMaskInpageProvider(this.provider)) {
      return (await this.provider?.request({
        method: "eth_chainId",
      })) as string;
    }
    return "";
  }

  onProviderNotFound(handler: () => void) {
    this.emitter.on(PROVIDER_NOT_FOUND_EVENT, handler);
  }

  getState() {
    return this.state;
  }

  private updateData = debounce(
    async () => {
      if (!this.web3) {
        this.state.connected = false;
        this.state.accounts = [];
        this.state.address = "";
        this.state.balances = [];
        return;
      }

      this.state.connected = true;
      this.state.accounts = (await this.web3.eth.getAccounts()) ?? [];
      this.state.address = this.state.accounts[0];
      this.state.balances = await this.getBalance();
    },
    100,
    { leading: true },
  );

  getAddress(): Address {
    return this.state.address;
  }

  isConnected() {
    return this.state.connected;
  }

  getSupportedTokens() {
    return this.supportedTokens;
  }

  async connect() {
    const provider = await this.providerPromise;
    try {
      if (!provider) {
        throw new Error("Cannot connect because provider is not yet loaded!");
      }
      this.web3 = new Web3(provider as provider);
      if (isMetaMaskInpageProvider(provider)) {
        if (provider.request) {
          await provider.request({ method: "eth_requestAccounts" });
        }
      }
      this.addWeb3Subscription();
      await this.updateData();
    } catch (err) {
      this.web3 = null;
      this.removeWeb3Subscription();
      throw err;
    }
  }

  addWeb3Subscription() {
    this.blockSubscription = this.web3?.eth.subscribe(
      "newBlockHeaders",
      (error, blockHeader) => {
        if (blockHeader) {
          this.updateData();
          this.state.log = blockHeader.hash;
        } else {
          this.state.log = error.message;
        }
      },
    );
  }

  removeWeb3Subscription() {
    const success = this.blockSubscription?.unsubscribe();
    if (success) {
      this.blockSubscription = null;
    } else {
      // try again if not success
      this.blockSubscription?.unsubscribe();
    }
  }

  async disconnect() {
    if (isMetaMaskInpageProvider(this.provider)) {
      (this.provider as any).disconnect &&
        (this.provider as any).disconnect(0, "Website disconnected wallet");
    }
    this.removeWeb3Subscription();
    this.web3 = null;
    await this.updateData();
  }

  async getBalance(address?: Address, asset?: Asset): Promise<Balances> {
    const supportedTokens = this.getSupportedTokens();
    const addr = address || this.state.address;
    if (!this.web3 || !addr) {
      return [];
    }

    const web3 = this.web3;
    let balances = [];

    if (asset) {
      if (!asset.address) {
        // Asset must be eth
        const ethBalance = await getEtheriumBalance(web3, addr);
        balances = [ethBalance];
      } else {
        // Asset must be ERC-20
        const tokenBalance = await getTokenBalance(web3, addr, asset);
        balances = [tokenBalance];
      }
    } else {
      // No address no asset get everything
      balances = await Promise.all([
        getEtheriumBalance(web3, addr),
        ...supportedTokens
          .filter((t) => t.symbol !== "eth")
          .map((token: Asset) => {
            if (token.address) return getTokenBalance(web3, addr, token);
            return AssetAmount(token, "0");
          }),
      ]);
    }

    return balances;
  }

  async transfer(params: TxParams): Promise<TxHash> {
    // TODO: validate params!!
    if (!this.web3) {
      throw new Error(
        "Cannot do transfer because there is not yet a connection to Ethereum.",
      );
    }

    const { amount, recipient, asset } = params;
    const from = this.getAddress();

    if (!from) {
      throw new Error(
        "Transaction attempted but 'from' address cannot be determined!",
      );
    }

    return await transferAsset(this.web3, from, recipient, amount, asset);
  }

  async signAndBroadcast(msg: Msg, mmo?: string) {}

  async setPhrase(args: string) {
    // We currently delegate auth to metamask so this is irrelavent
    return "";
  }

  purgeClient() {
    // We currently delegate auth to metamask so this is irrelavent
  }

  static create({
    getWeb3Provider,
    assets,
  }: EthereumServiceContext): IWalletService {
    return new EthereumService(getWeb3Provider, assets);
  }
}

export default EthereumService.create;
