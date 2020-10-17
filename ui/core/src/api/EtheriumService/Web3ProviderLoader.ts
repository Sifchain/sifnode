import { EventEmitter2 } from "eventemitter2";
import { IpcProvider, provider, WebsocketProvider } from "web3-core";

function isEventEmittingProvider(
  provider?: provider
): provider is WebsocketProvider | IpcProvider {
  if (!provider || typeof provider === "string") return false;
  return typeof (provider as any).on === "function";
}
type MetaMaskProvider = WebsocketProvider & {
  request?: (a: any) => Promise<void>;
};
function isMetaMaskProvider(provider?: provider): provider is MetaMaskProvider {
  return typeof (provider as any).request === "function";
}

// Load a provider and delegate the connected and disconnected events to it
export class Web3ProviderLoader extends EventEmitter2 {
  public provider: provider | null = null;

  constructor(private _getProvider: () => Promise<provider>) {
    super();
  }

  async load() {
    this.provider = await this._getProvider();

    if (isEventEmittingProvider(this.provider)) {
      // Forward events to parent
      this.provider.on("connected", () => {
        this.emit("connected");
      });

      this.provider.on("disconnected", () => {
        this.emit("disconnected");
      });
    }

    return this.provider;
  }

  getProvider() {
    return this.provider;
  }

  async attemptConnection() {
    // Let's test for Metamask
    if (isMetaMaskProvider(this.provider)) {
      if (this.provider.request) {
        // If metamask lets try and connect
        try {
          await this.provider.request({ method: "eth_requestAccounts" });
        } catch (err) {
          console.error(err);
        }
      }
    }
    this.emit("connected");
  }
}
