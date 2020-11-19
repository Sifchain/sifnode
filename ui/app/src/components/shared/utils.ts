import { Asset, Network } from "ui-core";

export function getAssetLabel(t: Asset) {
  if (t.network === Network.SIFCHAIN && t.symbol.indexOf("c") === 0) {
    return ["c", t.symbol.slice(1).toUpperCase()].join("");
  }
  return t.symbol.toUpperCase();
}
