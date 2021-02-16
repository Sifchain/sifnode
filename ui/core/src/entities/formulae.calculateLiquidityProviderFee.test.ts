import { calculateProviderFee } from "./formulae";
import { Fraction, TEN } from "./fraction/Fraction";
import tests from "../../../../test/test-tables/singleswap_liquidityfees.json";
import B from "./utils/B";
import JSBI from "jsbi";

tests.SingleSwapLiquidityFee.forEach(({ x, X, Y, expected }: any) => {
  // Need to convert inputs to JSBI to be able to test decimal input from tables.
  // In the actual logic, user input is converted before calculations are made.
  const bigx = B(x);
  const bigX = B(X);
  const bigY = B(Y);
  test(`Calc LP fee for swapping ${x}, expecting ${expected}`, () => {
    const output = calculateProviderFee(
      new Fraction(bigx), // Swap Amount
      new Fraction(bigX), // In Asset Pool Balance
      new Fraction(bigY) // Out Asset Pool Balance
    );
    expect(output.divide("1000000000000000000").toFixed(18)).toBe(expected);
  });
});
