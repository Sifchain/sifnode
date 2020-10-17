import { EventEmitter2 } from "eventemitter2";
import Web3 from "web3";
import { AbiItem } from "web3-utils";
import { ETH } from "../../constants";
import { Asset, Balance, Token } from "../../entities";
import { provider } from "web3-core";
import { IWalletService } from "..";
import { Web3ProviderLoader } from "./Web3ProviderLoader";

type Address = string;
type Balances = Balance[];

export type EtheriumServiceContext = {
  getWeb3Provider: () => Promise<provider>;
  getSupportedTokens: () => Promise<Token[]>;
};

// Hmm maybe we need to load each token from compiled json? Or is every ERC-20 token the same?
const generalTokenAbi: AbiItem[] = [
  // balanceOf
  {
    constant: true,
    inputs: [{ name: "_owner", type: "address" }],
    name: "balanceOf",
    outputs: [{ name: "balance", type: "uint256" }],
    type: "function",
  },
  // decimals
  {
    constant: true,
    inputs: [],
    name: "decimals",
    outputs: [{ name: "", type: "uint8" }],
    type: "function",
  },
];

function isToken(value?: Asset | Token): value is Token {
  return !value || Object.keys(value).includes("address");
}

async function getTokenBalance(web3: Web3, address: Address, asset: Token) {
  const contract = new web3.eth.Contract(generalTokenAbi, asset.address);
  const tokenBalance = await contract.methods.balanceOf(address).call();
  return Balance.create(asset, tokenBalance);
}

async function getEtheriumBalance(web3: Web3, address: Address) {
  const ethBalance = await web3.eth.getBalance(address);
  return Balance.create(ETH, ethBalance);
}

type ListenerFn = (...a: any[]) => void;

export class EtheriumService implements IWalletService {
  private web3: Web3 | null = null;
  private supportedTokens: Token[] = [];
  private address: Address | null = null;
  // private handleDisconnect: ListenerFn = () => {};
  private providerLoader: Web3ProviderLoader;

  constructor(
    getWeb3Provider: () => Promise<provider>,
    private getSupportedTokens: () => Promise<Token[]>
  ) {
    this.providerLoader = new Web3ProviderLoader(getWeb3Provider);
    this.providerLoader.load();
  }

  onDisconnected(handler: ListenerFn) {
    this.providerLoader.on("disconnect", handler);
  }

  onConnected(handler: ListenerFn) {
    this.providerLoader.on("connect", handler);
  }

  async getAddress(): Promise<Address | null> {
    if (!this.address) {
      [this.address] = (await this.web3?.eth.getAccounts()) ?? [];
    }

    return this.address;
  }

  isConnected() {
    return !!this.web3 || this.providerLoader.isConnected();
  }

  async connect() {
    try {
      this.supportedTokens = await this.getSupportedTokens();
      this.web3 = new Web3(this.providerLoader.getProvider());
      const accounts = await this.web3.eth.getAccounts();
      [this.address] = accounts;

      await this.providerLoader.connect();
    } catch (err) {
      this.web3 = null;
    }
  }

  async disconnect() {
    this.providerLoader.disconnect();
    this.web3 = null;
  }

  async getBalance(
    address?: Address,
    asset?: Asset | Token
  ): Promise<Balances> {
    const supportedTokens = this.supportedTokens;
    const addr = address || (await this.getAddress());

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

    return balances;
  }

  static create({
    getWeb3Provider,
    getSupportedTokens,
  }: EtheriumServiceContext): IWalletService {
    return new EtheriumService(getWeb3Provider, getSupportedTokens);
  }

  // FOLLOWING ARE YTI:

  // setPhrase(phrase: string): Address
  // setNetwork(net: Network): void
  // getNetwork(): Network

  // getExplorerAddressUrl(address: Address): string
  // getExplorerTxUrl(txID: string): string
  // getTransactions(params?: TxHistoryParams): Promise<TxsPage>

  // getFees(): Promise<Fees>

  // transfer(params: TxParams): Promise<TxHash>
  // deposit(params: TxParams): Promise<TxHash>

  // purgeClient(): void
}

export default EtheriumService.create;
