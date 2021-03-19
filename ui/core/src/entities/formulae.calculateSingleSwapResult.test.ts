import { Fraction } from "./fraction/Fraction";
import tests from "../../../../test/test-tables/singleswap_result.json";
import { calculateSwapResult } from "./formulae";

tests.SingleSwapResult.forEach(({ x, X, Y, expected }: any) => {
  test(`Swapping ${x}, expecting ${expected}`, () => {
    const output = calculateSwapResult(
      // External -> Native pool
      new Fraction(x), // Swap Amount
      new Fraction(X), // External Balance
      new Fraction(Y), // Native Balance
    );
    expect(output.toFixed(0)).toBe(expected);
  });
});
