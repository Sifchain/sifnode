import { Amount } from "./Amount";

import { calculateReverseSwapResult } from "./formulae";

const ReverseSwapAmounts = [
  {
    expected: "99999999999999999999999999",
    X: "1000000000000000000000000000",
    Y: "1000000000000000000000000000",
    S: "82644628099173553719008264",
  },
  {
    expected: "1",
    X: "1000000",
    Y: "100000000000000000000",
    S: "99999800000299",
  },
  {
    expected: "99999900000400800499",
    X: "100000000000000000000000",
    Y: "500000000",
    S: "499001",
  },
  {
    expected: "0",
    X: "100000000000000000000000",
    Y: "500000000",
    S: "0",
  },
  {
    expected: "0",
    X: "0",
    Y: "5000000",
    S: "0",
  },
  {
    expected: "10",
    X: "100",
    Y: "1000",
    S: "82",
  },
  {
    expected: "0",
    X: "100000000000000000000000000",
    Y: "0",
    S: "0",
  },
];

ReverseSwapAmounts.forEach(({ S, X, Y, expected }: any) => {
  const x = calculateReverseSwapResult(Amount(S), Amount(X), Amount(Y));
  test("", () => {
    expect(x.toBigInt().toString()).toBe(expected);
  });
});
