import { Network } from "../entities";
import { Amount } from "../entities/Amount";
import { AssetAmount } from "../entities/AssetAmount_";
import { Asset } from "../entities/Asset_";
import { exponentiateString, format, IFormatOptions } from "./format";

type Test = {
  only?: boolean;
  skip?: boolean;
  input: string;
  decimals?: number;
  options: IFormatOptions;
  expected: string;
};

describe("format", () => {
  const tests: Test[] = [
    {
      input: "100000000000",
      decimals: undefined,
      options: { mantissa: 2, separator: true },
      expected: `100,000,000,000.00`,
    },
    {
      input: "100000000000",
      decimals: undefined,
      options: { shorthand: true },
      expected: `100b`,
    },
    {
      input: "100000000000",
      decimals: undefined,
      options: { shorthand: true, mantissa: 6 },
      expected: `100.000000b`,
    },
    {
      input: "990000000000000000",
      decimals: 18,
      options: { mantissa: 6 },
      expected: `0.990000`,
    },
    {
      input: "990000000000000000",
      decimals: 18,
      options: { mantissa: 6, trimMantissa: true },
      expected: `0.99`,
    },
    {
      input: "999999800000000000",
      decimals: 18,
      options: { mantissa: 7 },
      expected: `0.9999998`,
    },
    {
      input: "100",
      decimals: undefined,
      options: { mode: "percent", mantissa: 1 },
      expected: `1.0%`,
    },
    {
      input: "1000",
      decimals: undefined,
      options: { mode: "percent", mantissa: 2 },
      expected: `10.00%`,
    },
    {
      input: "12345",
      decimals: undefined,
      options: { mode: "percent", mantissa: 3, exponent: 3 },
      expected: `12.345%`,
    },
    {
      input: "12345",
      decimals: undefined,
      options: { mode: "percent", mantissa: 3, exponent: 3, space: true },
      expected: `12.345 %`,
    },
    {
      input: "-990000000000000000",
      decimals: 18,
      options: { mantissa: 6, trimMantissa: true },
      expected: `-0.99`,
    },
    {
      input: "999999800000000000",
      decimals: 18,
      options: { mantissa: 7, shorthand: true, forceSign: true },
      expected: `+0.9999998`,
    },
  ];

  tests.forEach(({ only, skip, decimals, options, expected, input }) => {
    const tester = only ? test.only : skip ? test.skip : test;
    tester(expected, () => {
      const amount =
        typeof decimals !== "undefined"
          ? AssetAmount(
              Asset({
                symbol: "foo",
                decimals,
                address: "12345678",
                label: "cFOO",
                name: "Foo",
                network: Network.ETHEREUM,
              }),
              input,
            )
          : Amount(input);

      expect(format(amount, options)).toBe(expected);
    });
  });
});

describe("exponentiateString", () => {
  test("1", () => {
    expect(exponentiateString("12345678", -4)).toBe("1234.5678");
  });
});
