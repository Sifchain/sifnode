import { EventEmitter2 } from "eventemitter2";
import { IpcProvider, provider, WebsocketProvider } from "web3-core";
import { CONNECTED, DISCONNECTED } from "./events";

function isEventEmittingProvider(
  provider?: provider
): provider is WebsocketProvider | IpcProvider {
  if (!provider || typeof provider === "string") return false;
  return typeof (provider as any).on === "function";
}

type MetaMaskProvider = WebsocketProvider & {
  request?: (a: any) => Promise<void>;
  isConnected(): boolean;
};

function isMetaMaskProvider(provider?: provider): provider is MetaMaskProvider {
  return typeof (provider as any).request === "function";
}

// Load a provider and delegate the connect and disconnect events to it
export class Web3ProviderLoader extends EventEmitter2 {
  public provider: provider | null = null;

  constructor(private _getProvider: () => Promise<provider>) {
    super();
  }

  isConnected() {
    if (isMetaMaskProvider(this.provider)) {
      return this.provider.isConnected();
    }
    return false;
  }

  async load() {
    this.provider = await this._getProvider();

    if (isEventEmittingProvider(this.provider)) {
      // Forward events to parent
      this.provider.on(CONNECTED, () => {
        this.emit(CONNECTED);
      });

      this.provider.on(DISCONNECTED, () => {
        this.emit(DISCONNECTED);
      });

      this.provider.on("message", () => {
        const msg = arguments;
        console.log({ msg });
      });
    }

    return this.provider;
  }

  async disconnect() {
    this.emit(DISCONNECTED);
  }

  getProvider() {
    return this.provider;
  }

  async connect() {
    // Let's test for Metamask
    if (isMetaMaskProvider(this.provider)) {
      if (this.provider.request) {
        // If metamask lets try and connect
        await this.provider.request({ method: "eth_requestAccounts" });
      }
    }
    this.emit(CONNECTED);
  }
}
