import { calculatePriceImpact } from "./formulae";
import {Fraction, TEN} from "./fraction/Fraction";
import tests from "../../../../test/test-tables/sample_price_impact.json";
import B from "./utils/B";
import JSBI from "jsbi";

tests.Slip.forEach(({ x, X, expected } : any) => {
  // Need to convert inputs to JSBI to be able to test decimal input from tables.
  // In the actual logic, user input is converted before calculations are made.
  const bigx = B(x);
  const bigX = B(X);
  test(`Calc Price Impact for swapping ${x}, expecting ${expected}`, () => {
    const output = calculatePriceImpact(
      new Fraction(bigx, JSBI.exponentiate(TEN, JSBI.BigInt(18))), // Swap Amount
      new Fraction(bigX, JSBI.exponentiate(TEN, JSBI.BigInt(18))) // In Asset Pool Balance
    );
    expect(output.toFixed(18)).toBe(expected);
  });
});
