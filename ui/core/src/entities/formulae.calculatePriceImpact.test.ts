import { calculatePriceImpact } from "./formulae";
import { Fraction, TEN } from "./fraction/Fraction";
import { SingleSwapStandardSlip } from "../../../../test/test-tables/singleswap_standardslip.json";
import B from "./utils/B";
import JSBI from "jsbi";

SingleSwapStandardSlip.forEach(({ x, X, expected }: any) => {
  // Need to convert inputs to JSBI to be able to test decimal input from tables.
  // In the actual logic, user input is converted before calculations are made.
  const bigx = B(x);
  const bigX = B(X);
  test.skip(`Calc Price Impact for swapping ${x}, expecting ${expected}`, () => {
    const output = calculatePriceImpact(
      new Fraction(bigx), // Swap Amount
      new Fraction(bigX) // In Asset Pool Balance
    );
    expect(output.divide("1000000000000000000").toFixed(18)).toBe(expected);
  });
});
