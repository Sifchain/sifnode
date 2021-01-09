import { computed, Ref } from "@vue/reactivity";
import { Asset, Token, AssetAmount, IPool, Pool } from "../entities";

export function assetPriceMessage(
  amount: AssetAmount | null,
  pair: IPool | null,
  decimals: number = -1
) {
  if (!pair || !amount || amount.equalTo("0")) return "";
  const swapResult = pair.calcSwapResult(amount);

  const assetPriceStr = [
    swapResult
      .divide(amount)
      .toFixed(decimals > -1 ? decimals : amount.asset.decimals),
    HACK_labelDecorator(swapResult.asset.symbol),
  ].join(" ");

  const formattedPerSymbol = HACK_labelDecorator(amount.asset.symbol);

  return `${assetPriceStr} per ${formattedPerSymbol}`;
}

export function trimZeros(amount: string) {
  if (amount.indexOf(".") === -1) return `${amount}.0`;
  const tenDecimalsMax = parseFloat(amount).toFixed(10);
  return tenDecimalsMax.replace(/0+$/, "").replace(/\.$/, ".0");
}

export function useBalances(balances: Ref<AssetAmount[]>) {
  return computed(() => {
    const map = new Map<string, AssetAmount>();

    for (const item of balances.value) {
      map.set(item.asset.symbol, item);
    }
    return map;
  });
}

export function buildAsset(val: string | null) {
  return val === null ? val : Asset.get(val);
}

export function buildAssetAmount(asset: Asset | null, amount: string) {
  return asset ? AssetAmount(asset, amount) : asset;
}

// Major HACK alert to get the demo working for a demo instead of changing a tonne of tests
// we are simply renaming tokens in the view to look like real ERC-20 tokens

export function HACK_labelDecorator(symbol: string) {
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

export function HACK_assetDecorator(asset: Asset) {
  return Token({
    address: "",
    ...asset,
    symbol: HACK_labelDecorator(asset.symbol),
  });
}

export function HACK_assetAmountDecorator(assetAmount: AssetAmount) {
  return AssetAmount(
    HACK_assetDecorator(assetAmount.asset),
    assetAmount.amount
  );
}

// end HACK
