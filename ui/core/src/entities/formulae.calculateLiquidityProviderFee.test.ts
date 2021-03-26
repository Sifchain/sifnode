import { calculateProviderFee } from "./formulae";
import tests from "../../../../test/test-tables/singleswap_liquidityfees.json";

import { Amount } from "./Amount";

tests.SingleSwapLiquidityFee.forEach(({ x, X, Y, expected }: any) => {
  test(`Calc LP fee for swapping ${x}, expecting ${expected}`, () => {
    const output = calculateProviderFee(
      Amount(x), // Swap Amount
      Amount(X), // In Asset Pool Balance
      Amount(Y), // Out Asset Pool Balance
    );
    expect(output.toBigInt().toString()).toBe(expected);
  });
});
