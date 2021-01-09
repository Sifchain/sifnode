// Major HACK alert to get the demo working instead of changing a tonne of tests
// we are simply renaming tokens in the view to look like real ERC-20 tokens

import { Asset, AssetAmount, Token } from "ui-core";

export function labelDecorator(symbol: string) {
  // Return fake labels
  if (symbol.toUpperCase() === "ATK") {
    return "USDC";
  }

  // Return fake labels
  if (symbol.toUpperCase() === "BTK") {
    return "LINK";
  }

  // Return fake labels
  if (symbol.toUpperCase() === "CATK") {
    return "cUSDC";
  }

  // Return fake labels
  if (symbol.toUpperCase() === "CBTK") {
    return "cLINK";
  }
  return symbol;
}

export function assetDecorator(asset: Asset) {
  return Token({
    address: "",
    ...asset,
    symbol: labelDecorator(asset.symbol),
  });
}

export function assetAmountDecorator(assetAmount: AssetAmount) {
  return AssetAmount(assetDecorator(assetAmount.asset), assetAmount.amount);
}
