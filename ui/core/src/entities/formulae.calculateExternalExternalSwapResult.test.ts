import { calculateExternalExternalSwapResult } from "./formulae";
import { Fraction } from "./fraction/Fraction";

const tests = [
  {
    skip: false,
    only: false,
    name: "even",
    input: {
      ax: "1",
      aX: "8300000", // eth
      aY: "10000000000",
      bX: "10000000000", // cusdc
      bY: "10000000000",
    },
    expected: "1204.818696472882384427",
  },

  {
    skip: false,
    only: false,
    name: "even",
    input: {
      ax: "1",
      aX: "588235000", // link
      aY: "10000000000",
      bX: "10000000000", // cusdc
      bY: "10000000000",
    },
    expected: "17.000008384404135090",
  },
];

tests.forEach(({ name, only, skip, input, expected }) => {
  const tester = only ? test.only : skip ? test.skip : test;

  tester(name, () => {
    const output = calculateExternalExternalSwapResult(
      // External -> Native pool
      new Fraction(input.ax), // Swap Amount
      new Fraction(input.aX), // External Balance
      new Fraction(input.aY), // Native Balance
      // Native -> External pool
      new Fraction(input.bX), // External Balance
      new Fraction(input.bY) // Native Balance
    );

    expect(output.toFixed(18)).toBe(expected);
  });
});
