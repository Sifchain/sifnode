import { reactive } from "@vue/reactivity";
import Web3 from "web3";
import { provider } from "web3-core";
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
  transferAsset,
} from "./utils/ethereumUtils";
import { isToken } from "../../entities/utils/isToken";
import notify from "../utils/Notifications";

type Address = string;
type Balances = AssetAmount[];

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

export class EthereumService implements IWalletService {
  private web3: Web3 | null = null;
  private supportedTokens: Asset[] = [];

  // This is shared reactive state
  private state: {
    connected: boolean;
    address: Address;
    accounts: Address[];
    balances: AssetAmount[];
    log: string;
  };

  constructor(getWeb3Provider: () => Promise<provider>, assets: Asset[]) {
    // init state
    this.state = reactive({ ...initState });
    this.supportedTokens = assets.filter(t => t.network === Network.ETHEREUM);
    this.web3 = new Web3(Web3.givenProvider || "ws://localhost:7545");

    if (this.web3 && this.web3.givenProvider) {
      this.web3.givenProvider.on('accountsChanged', () => {
        this.updateData();
      });

      this.web3.givenProvider.on('chainChanged', () => {
        window.location.reload();
      });

      this.web3.givenProvider.on('message', (message: any) => {
        console.log("onMessage", message);
        this.updateData();
      });
    }
  }

  getState() {
    return this.state;
  }

  private updateData = debounce(
    async () => {
      if (!this.web3) {
        this.state = reactive({ ...initState });
        return;
      }
      this.state.accounts = (await this.web3.eth.getAccounts()) ?? [];
      this.state.connected = this.state.accounts.length > 0;
      if (this.state.connected) {
        this.state.address = this.state.accounts[0];
        this.state.balances = await this.getBalance();
      } else {
        this.state.address = "";
        this.state.balances = [];
        this.state.log = "unset";
      }
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
    try {
      if (!this.web3) {
        throw new Error("There is not yet a connection to Ethereum.");
      }
      await this.web3.givenProvider.request({ method: "eth_requestAccounts" });
      notify({ type: "success", message: "Connected to Metamask" });
      await this.updateData();
    } catch (err) {
      console.log(err);
      this.state = reactive({ ...initState });
    }
  }

  async disconnect() {
    await this.web3?.givenProvider.request({
      method: "wallet_requestPermissions",
      params: [
        {
          eth_accounts: {}
        }
      ]
    });
    this.state = reactive({ ...initState });
  }

  async getBalance(
    address?: Address,
    asset?: Asset | Token
  ): Promise<Balances> {
    let balances: any[] = [];
    const addr = address || this.state.address;
    const web3 = this.web3;
    if (!web3 || !addr) {
      return balances;
    }

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
      const supportedTokens = this.getSupportedTokens();
      // No address no asset get everything
      balances = await Promise.all([
        getEtheriumBalance(web3, addr),
        ...supportedTokens
          .slice(0, 10)
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
