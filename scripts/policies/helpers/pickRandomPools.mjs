import _ from "lodash";

export function pickRandomPools(pools, entries) {
  return _.sampleSize(pools, 10).map(
    ({
      external_asset: { symbol },
      swap_price_external: swapPriceExternal,
    }) => {
      const { decimals } = entries.find(({ denom }) => denom === symbol);
      return { symbol, decimals, swapPriceExternal };
    }
  );
}
