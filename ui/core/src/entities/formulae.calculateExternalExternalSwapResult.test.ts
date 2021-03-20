import { calculateExternalExternalSwapResult } from "./formulae";
import { Amount } from "./Amount";
import tests from "../../../../test/test-tables/doubleswap_result.json";

tests.DoubleSwap.forEach(({ ax, aX, aY, bX, bY, expected }: any) => {
  test(`Swapping ${ax}, expecting ${expected}`, () => {
    const output = calculateExternalExternalSwapResult(
      // External -> Native pool
      Amount(ax), // Swap Amount
      Amount(aX), // External Balance
      Amount(aY), // Native Balance
      // Native -> External pool
      Amount(bX), // External Balance
      Amount(bY), // Native Balance
    );
    expect(output.toString()).toBe(expected);
  });
});
