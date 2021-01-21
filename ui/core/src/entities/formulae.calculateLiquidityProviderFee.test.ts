import { calculateProviderFee } from "./formulae";
import { Fraction } from "./fraction/Fraction";
import tests from "../../../../test/test-tables/sample_liquidity_fee.json";

tests.LiquidityFee.forEach(({ x, X, Y, expected } : any) => {
  test(`Calc LP fee for swapping ${x}, expecting ${expected}`, () => {
    const output = calculateProviderFee(
      new Fraction(x), // Swap Amount
      new Fraction(X), // In Asset Pool Balance
      new Fraction(Y) // Out Asset Pool Balance
    );
    expect(output.toFixed(18)).toBe(expected);
  });
});
