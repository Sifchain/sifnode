import { Coin } from "./Coin";
import { Token } from "./Token";

export type Asset = Token | Coin;
const assetMap = new Map<string, Asset>();

export const Asset = {
  set(key: string, value: Asset) {
    if (!key) return;

    assetMap.set(key.toLowerCase(), value);
  },
  get(key: string | Asset): Asset {
    key = typeof key == "string" ? key : key.symbol;
    const found = key ? assetMap.get(key.toLowerCase()) : false;
    if (!found) {
      console.log(
        "Available keys: " +
          Array.from(assetMap.keys())
            .sort()
            .join(","),
      );
      throw new Error(
        `Attempt to retrieve the asset with key ${key} before it had been created.`,
      );
    }
    return found;
  },
};
