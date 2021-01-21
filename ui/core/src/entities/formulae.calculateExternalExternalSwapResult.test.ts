import { calculateExternalExternalSwapResult } from "./formulae";
import { Fraction } from "./fraction/Fraction";
import tests from "../../../../test/test-tables/sample_swaps.json";

tests.Swap.forEach(({ ax, aX, aY, bX, bY, expected } : any) => {
  test(`Swapping ${ax}, expecting ${expected}`, () => {
    const output = calculateExternalExternalSwapResult(
      // External -> Native pool
      new Fraction(ax), // Swap Amount
      new Fraction(aX), // External Balance
      new Fraction(aY), // Native Balance
      // Native -> External pool
      new Fraction(bX), // External Balance
      new Fraction(bY) // Native Balance
    );
    expect(output.toFixed(18)).toBe(expected);
  });
});
