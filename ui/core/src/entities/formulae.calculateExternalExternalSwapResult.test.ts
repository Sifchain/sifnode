import { calculateExternalExternalSwapResult } from "./formulae";
import { Fraction, TEN } from "./fraction/Fraction";
import tests from "../../../../test/test-tables/doubleswap_result.json";
import B from "./utils/B";
import JSBI from "jsbi";

tests.DoubleSwap.forEach(({ ax, aX, aY, bX, bY, expected }: any) => {
  // Need to convert inputs to JSBI to be able to test decimal input from tables.
  // In the actual logic, user input is converted before calculations are made.
  const bigax = B(ax);
  const bigaX = B(aX);
  const bigaY = B(aY);
  const bigbX = B(bX);
  const bigbY = B(bY);
  test(`Swapping ${ax}, expecting ${expected}`, () => {
    const output = calculateExternalExternalSwapResult(
      // External -> Native pool
      new Fraction(bigax, JSBI.exponentiate(TEN, JSBI.BigInt(18))), // Swap Amount
      new Fraction(bigaX, JSBI.exponentiate(TEN, JSBI.BigInt(18))), // External Balance
      new Fraction(bigaY, JSBI.exponentiate(TEN, JSBI.BigInt(18))), // Native Balance
      // Native -> External pool
      new Fraction(bigbX, JSBI.exponentiate(TEN, JSBI.BigInt(18))), // External Balance
      new Fraction(bigbY, JSBI.exponentiate(TEN, JSBI.BigInt(18))), // Native Balance
    );
    expect(output.toFixed(0)).toBe(expected);
  });
});
