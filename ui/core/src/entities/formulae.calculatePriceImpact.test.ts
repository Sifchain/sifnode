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
      new Fraction(bigx, JSBI.exponentiate(TEN, JSBI.BigInt(18))), // Swap Amount
      new Fraction(bigX, JSBI.exponentiate(TEN, JSBI.BigInt(18))) // In Asset Pool Balance
    );
    expect(output.toFixed(18)).toBe(expected);
  });
});
