import tests from "../../../../test/test-tables/reverse_single_swap_result.json";
import { calculateReverseSwapResult } from "./formulae";
import Big from "big.js";

tests.SingleSwap.forEach(({ S, X, Y, expected }: any) => {
  const bigS = Big(S);
  const bigX = Big(X);
  const bigY = Big(Y);
  test(`Swapping ${S}, expecting ${expected}`, () => {
    const output = calculateReverseSwapResult(bigS, bigX, bigY);
    expect(output.toFixed(0)).toBe(expected);
  });
});
