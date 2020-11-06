import { reactive } from "@vue/reactivity";
import Web3 from "web3";
import { provider, WebsocketProvider } from "web3-core";
import { IWalletService } from "../IWalletService";
import { TxHash, TxParams, Asset, Balance, Token } from "../../entities";
import {
  getEtheriumBalance,
  getTokenBalance,
  isEventEmittingProvider,
  isToken,
  transferAsset,
} from "./utils/ethereumUtils";

type Address = string;
type Balances = Balance[];

export type EthereumServiceContext = {
  getWeb3Provider: () => Promise<provider>;
  getSupportedTokens: () => Promise<Token[]>;
};

type MetaMaskProvider = WebsocketProvider & {
  request?: (a: any) => Promise<void>;
  isConnected(): boolean;
};

function isMetaMaskProvider(provider?: provider): provider is MetaMaskProvider {
  return typeof (provider as any).request === "function";
}

export class EthereumService implements IWalletService {
  private web3: Web3 | null = null;
  private supportedTokens: Token[] = [];
  private blockSubscription: any;
  private provider: provider | undefined;

  // This is shared reactive state
  private state: {
    connected: boolean;
    address: Address;
    accounts: Address[];
    balances: Balance[];
    log: string;
  } = reactive({
    connected: false,
    accounts: [],
    balances: [],
    address: "",
    log: "unset",
  });

  constructor(
    getWeb3Provider: () => Promise<provider>,
    private getSupportedTokens: () => Promise<Token[]>
  ) {
    getWeb3Provider().then((provider) => {
      if (isEventEmittingProvider(provider)) {
        provider.on("connect", () => {
          this.state.connected = true;
        });
        provider.on("disconnect", () => {
          this.state.connected = false;
        });
      }
      this.provider = provider;
    });
  }

  getState() {
    return this.state;
  }

  private async updateAccounts() {
    if (!this.web3) return;
    this.state.accounts = (await this.web3?.eth.getAccounts()) ?? [];
  }

  getAddress(): Address {
    return this.state.address;
  }

  isConnected() {
    return this.state.connected;
  }

  async connect() {
    try {
      this.supportedTokens = await this.getSupportedTokens();

      if (!this.provider)
        throw new Error("Cannot connect because provider is not yet loaded!");

      this.web3 = new Web3(this.provider);

      await this.updateAccounts();

      [this.state.address] = this.state.accounts;

      // Let's test for Metamask
      if (isMetaMaskProvider(this.provider)) {
        if (this.provider.request) {
          // If metamask lets try and connect
          await this.provider.request({ method: "eth_requestAccounts" });
        }
      }
      this.state.connected = true;

      this.addWeb3Subscription();
    } catch (err) {
      this.web3 = null;
    }
  }

  addWeb3Subscription() {
    this.blockSubscription = this.web3?.eth.subscribe(
      "newBlockHeaders",
      (error, result) => {
        this.state.log = result?.hash ?? "null";
      }
    );
  }

  removeWeb3Subscription() {
    this.blockSubscription?.unsubscribe();
  }

  async disconnect() {
    this.removeWeb3Subscription();
    this.state.connected = false;
    this.web3 = null;
  }

  async getBalance(
    address?: Address,
    asset?: Asset | Token
  ): Promise<Balances> {
    const supportedTokens = this.supportedTokens;
    const addr = address || this.getAddress();

    if (!this.web3 || !addr) return [];
    const web3 = this.web3;

    if (asset) {
      if (!isToken(asset)) {
        // Asset must be eth
        const ethBalance = await getEtheriumBalance(web3, addr);
        return [ethBalance];
      }

      // Asset must be ERC-20
      const tokenBalance = await getTokenBalance(web3, addr, asset);
      return [tokenBalance];
    }

    // No address no asset get everything
    const balances = await Promise.all([
      getEtheriumBalance(web3, addr),
      ...supportedTokens.map((token: Token) => {
        return getTokenBalance(web3, addr, token);
      }),
    ]);

    this.state.balances = balances;
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

  async setPhrase() {
    // We currently delegate auth to metamask so this is irrelavent
    return "";
  }

  purgeClient() {
    // We currently delegate auth to metamask so this is irrelavent
  }

  static create({
    getWeb3Provider,
    getSupportedTokens,
  }: EthereumServiceContext): IWalletService {
    return new EthereumService(getWeb3Provider, getSupportedTokens);
  }
}

export default EthereumService.create;
