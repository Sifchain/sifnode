import _ from "lodash";

export function pickRandomPools(pools, entries, nPools) {
  return _.sampleSize(pools, nPools).map(
    ({
      external_asset: { symbol },
      swap_price_external: swapPriceExternal,
    }) => {
      const { decimals } = entries.find(({ denom }) => denom === symbol);
      return { symbol, decimals, swapPriceExternal };
    }
  );
}
