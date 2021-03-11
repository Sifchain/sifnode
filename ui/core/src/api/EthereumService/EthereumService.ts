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
  Token,
  Network,
} from "../../entities";
import {
  getEtheriumBalance,
  getTokenBalance,
  isEventEmittingProvider,
  transferAsset,
} from "./utils/ethereumUtils";
import { isToken } from "../../entities/utils/isToken";

type Address = string;
type Balances = AssetAmount[];

export type EthereumServiceContext = {
  getWeb3Provider: () => Promise<provider>;
  assets: Asset[];
};

type MetaMaskProvider = WebsocketProvider & {
  request?: (a: any) => Promise<void>;
  isConnected(): boolean;
};

function isMetaMaskProvider(provider?: provider): provider is MetaMaskProvider {
  return typeof (provider as any).request === "function";
}

const initState = {
  connected: false,
  accounts: [],
  balances: [],
  address: "",
  log: "unset",
};

export class EthereumService implements IWalletService {
  private web3: Web3 | null = null;
  private supportedTokens: Asset[] = [];
  private blockSubscription: any;
  private provider: provider | undefined;
  private providerPromise: Promise<provider>;

  // This is shared reactive state
  private state: {
    connected: boolean;
    address: Address;
    accounts: Address[];
    balances: AssetAmount[];
    log: string;
  };

  constructor(getWeb3Provider: () => Promise<provider>, assets: Asset[]) {
    this.state = reactive({ ...initState });
    this.supportedTokens = assets.filter(t => t.network === Network.ETHEREUM);
    this.providerPromise = getWeb3Provider();
    this.providerPromise
      .then(provider => {
        if (!provider) {
          return (this.provider = null);
        }
        if (isEventEmittingProvider(provider)) {
          provider.on("chainChanged", () => window.location.reload());
          provider.on("accountsChanged", () => this.updateData());
        }
        this.web3 = new Web3(provider);
        this.provider = provider;
        this.addWeb3Subscription();
        this.updateData();
      })
      .catch(error => {
        console.log("error", error);
      });
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
      this.state.connected = !!this.web3;
      this.state.accounts = (await this.web3.eth.getAccounts()) ?? [];
      this.state.address = this.state.accounts[0];
      this.state.balances = await this.getBalance();
    },
    100,
    { leading: true }
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
      this.web3 = new Web3(provider);
      if (isMetaMaskProvider(provider)) {
        if (provider.request) {
          await provider.request({ method: "eth_requestAccounts" });
        }
      }
      this.addWeb3Subscription();
      await this.updateData();
    } catch (err) {
      console.log(err);
      this.web3 = null;
      this.removeWeb3Subscription();
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
      }
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
    if (isMetaMaskProvider(this.provider)) {
      this.provider.disconnect &&
        this.provider.disconnect(0, "Website disconnected wallet");
    }
    this.removeWeb3Subscription();
    this.web3 = null;
    await this.updateData();
  }

  async getBalance(
    address?: Address,
    asset?: Asset | Token
  ): Promise<Balances> {
    const supportedTokens = this.getSupportedTokens();
    const addr = address || this.state.address;
    if (!this.web3 || !addr) {
      return [];
    }

    const web3 = this.web3;
    let balances = [];

    if (asset) {
      if (!isToken(asset)) {
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
          .filter(t => t.symbol !== "eth")
          .map((token: Asset) => {
            if (isToken(token)) return getTokenBalance(web3, addr, token);
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
        "Cannot do transfer because there is not yet a connection to Ethereum."
      );
    }

    const { amount, recipient, asset } = params;
    const from = this.getAddress();

    if (!from) {
      throw new Error(
        "Transaction attempted but 'from' address cannot be determined!"
      );
    }

    return await transferAsset(this.web3, from, recipient, amount, asset);
  }

  async signAndBroadcast() {}

  async setPhrase() {
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
