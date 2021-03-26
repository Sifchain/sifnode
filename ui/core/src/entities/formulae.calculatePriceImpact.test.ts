import { calculatePriceImpact } from "./formulae";
import { SingleSwapStandardSlip } from "../../../../test/test-tables/singleswap_standardslip.json";
import { Amount } from "./Amount";

SingleSwapStandardSlip.forEach(({ x, X, expected }: any) => {
  // Need to convert inputs to JSBI to be able to test decimal input from tables.
  // In the actual logic, user input is converted before calculations are made.
  test(`Calc Price Impact for swapping ${x}, expecting ${expected}`, () => {
    const output = calculatePriceImpact(
      Amount(x), // Swap Amount
      Amount(X), // In Asset Pool Balance
    );
    expect(output.toString()).toBe(expected);
  });
});
