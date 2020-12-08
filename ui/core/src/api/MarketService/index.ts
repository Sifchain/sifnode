import ReconnectingWebSocket from "reconnecting-websocket";
import { Asset, Pool } from "../../entities";
import { SifUnSignedClient } from "../utils/SifClient";
import { toPool } from "../utils/toPool";

export type MarketServiceContext = {
  loadAssets: () => Promise<Asset[]>;
  sifApiUrl: string;
  nativeAsset: Asset;
};

type PoolHandlerFn = (pools: Pool[]) => void;

export default function createMarketService({
  loadAssets,
  sifApiUrl,
}: MarketServiceContext) {
  let ws: ReconnectingWebSocket;
  const sifClient = new SifUnSignedClient(sifApiUrl);

  let poolHandler: PoolHandlerFn = () => {};

  async function setupPoolWatcher() {
    await new Promise((res, rej) => {
      ws = new ReconnectingWebSocket("ws://localhost:26657/websocket");
      ws.onopen = () => {
        ws.send(
          JSON.stringify({
            jsonrpc: "2.0",
            method: "subscribe",
            id: "1",
            params: {
              query: `tm.event='Tx'`,
            },
          })
        );
        // This assumes every transaction means an update to pool values
        // Subscribing to all pool addresses could mean having a tone of
        // open connections to our node because there is no "OR" query
        // syntax so have chosen to go with debouncing getPools for now.
        ws.onmessage = async (...argoids) => {
          console.log({ argoids });
          poolHandler(await instance.getPools());
        };
        res(ws);
      };
      ws.onerror = (err) => rej(err);
    });
  }

  async function initialize() {
    await loadAssets();
    await setupPoolWatcher();
  }

  initialize();

  const instance = {
    async getPools() {
      const rawPools = await sifClient.getPools();
      return rawPools.map(toPool);
    },
    onPoolsUpdated(handler: PoolHandlerFn) {
      poolHandler = handler;
    },
  };

  return instance;
}
