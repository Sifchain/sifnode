import Web3 from "web3";
import { AbiItem } from "web3-utils";
import { ETH } from "../../constants";
import { Asset, Balance, Token } from "../../entities";
import { provider } from "web3-core";
import { IWalletService, TxParams, TxHash } from "..";
import { Web3ProviderLoader } from "./Web3ProviderLoader";
import JSBI from "jsbi";
import { EventEmitter2 } from "eventemitter2";
import { CHANGE, CONNECTED, DISCONNECTED } from "./events";

// import { asset } from "src/store/asset";

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
  // transfer
  {
    constant: false,
    inputs: [
      {
        name: "_to",
        type: "address",
      },
      {
        name: "_value",
        type: "uint256",
      },
    ],
    name: "transfer",
    outputs: [
      {
        name: "",
        type: "bool",
      },
    ],
    type: "function",
  },
];

function isToken(value?: Asset | Token): value is Token {
  return value ? Object.keys(value).includes("address") : false;
}

function getTokenContract(web3: Web3, asset: Token) {
  return new web3.eth.Contract(generalTokenAbi, asset.address);
}

async function getTokenBalance(web3: Web3, address: Address, asset: Token) {
  const contract = getTokenContract(web3, asset);
  const tokenBalance = await contract.methods.balanceOf(address).call();
  return Balance.create(asset, tokenBalance);
}

async function transferToken(
  web3: Web3,
  fromAddress: Address,
  toAddress: Address,
  amount: JSBI,
  asset: Token
) {
  const contract = getTokenContract(web3, asset);
  return new Promise<string>((resolve, reject) => {
    let hash: string;
    let receipt: boolean;

    function resolvePromise() {
      if (receipt && hash) resolve(hash);
    }

    contract.methods
      .transfer(toAddress, JSBI.toNumber(amount))
      .send({ from: fromAddress })
      .on("transactionHash", (_hash: string) => {
        hash = _hash;
        resolvePromise();
      })
      .on("receipt", (_receipt: boolean) => {
        receipt = _receipt;
        resolvePromise();
      })
      .on("error", (err: any) => {
        reject(err);
      });
  });
}

async function transferEther(
  web3: Web3,
  fromAddress: Address,
  toAddress: Address,
  amount: JSBI
) {
  return new Promise<string>((resolve, reject) => {
    web3.eth
      .sendTransaction({
        from: fromAddress,
        to: toAddress,
        value: amount.toString(),
      })
      .on("transactionHash", (hash) => {
        resolve(hash);
      })
      .on("error", (err) => reject(err));
  });
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
  private providerLoader: Web3ProviderLoader; // loader for the provider
  private emitter: EventEmitter2; // event emitter
  private connected: boolean = false;

  constructor(
    getWeb3Provider: () => Promise<provider>,
    private getSupportedTokens: () => Promise<Token[]>
  ) {
    this.emitter = new EventEmitter2();
    this.providerLoader = new Web3ProviderLoader(getWeb3Provider);
    this.providerLoader.load();
    this.providerLoader.on(DISCONNECTED, () => {
      this.emitter.emit(CHANGE);
    });
    this.providerLoader.on(CONNECTED, () => {
      this.emitter.emit(CHANGE);
    });
  }

  onDisconnected(handler: ListenerFn) {
    if (!this.providerLoader.listeners(DISCONNECTED).includes(handler)) {
      this.providerLoader.on(DISCONNECTED, (...args: any[]) => {
        this.connected = false;
        handler(...args);
      });
    }
  }

  onConnected(handler: ListenerFn) {
    if (!this.providerLoader.listeners(CONNECTED).includes(handler)) {
      this.providerLoader.on(CONNECTED, (...args: any[]) => {
        this.connected = true;
        handler(...args);
      });
    }
  }

  onChange(handler: ListenerFn) {
    if (!this.emitter.listeners(CHANGE).includes(handler)) {
      this.emitter.on(CHANGE, handler);
    }
  }

  private reportChange() {
    const { address, connected } = this;
    this.emitter.emit(CHANGE, { address, connected });
  }

  async getAddress(): Promise<Address | null> {
    if (!this.address) {
      [this.address] = (await this.web3?.eth.getAccounts()) ?? [];
    }

    return this.address;
  }

  isConnected() {
    return this.connected;
  }

  async connect() {
    try {
      this.supportedTokens = await this.getSupportedTokens();
      this.web3 = new Web3(this.providerLoader.getProvider());
      const accounts = await this.web3.eth.getAccounts();
      [this.address] = accounts;

      await this.providerLoader.connect();
      this.connected = true;
      this.reportChange();
    } catch (err) {
      this.web3 = null;
    }
  }

  async disconnect() {
    this.connected = false;
    this.web3 = null;
    this.providerLoader.disconnect();
    this.reportChange();
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

  async transfer(params: TxParams): Promise<TxHash> {
    // TODO: validate params!!
    if (!this.web3) {
      throw new Error(
        "Cannot do transfer because there is not yet a connection to Ethereum."
      );
    }

    const { amount, recipient, asset } = params;
    const from = await this.getAddress();
    if (!from) {
      throw new Error(
        "Transaction attempted but 'from' address cannot be determined!"
      );
    }
    console.log({ asset });
    const hash = isToken(asset)
      ? await transferToken(this.web3, from, recipient, amount, asset)
      : await transferEther(this.web3, from, recipient, amount);
    this.reportChange();
    return hash;
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
