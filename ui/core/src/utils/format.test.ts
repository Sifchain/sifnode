import { Network } from "../entities";
import { Amount } from "../entities/Amount";
import { Asset } from "../entities/Asset";
import { format, IFormatOptions } from "./format";

type Test = {
  only?: boolean;
  skip?: boolean;
  // input: string;
  // decimals?: number;
  // options: IFormatOptions;
  input: string;
  expected: string;
};

function mockAsset(decimals: number) {
  return Asset({
    symbol: "foo",
    decimals,
    address: "12345678",
    label: "cFOO",
    name: "Foo",
    network: Network.ETHEREUM,
  });
}

describe("format", () => {
  const tests: Test[] = [
    {
      input: format(Amount("100000000000"), { mantissa: 2, separator: true }),
      expected: `100,000,000,000.00`,
    },
    {
      input: format(Amount("100000000000"), { shorthand: true }),
      expected: `100b`,
    },
    {
      input: format(Amount("100000000000"), { shorthand: true, mantissa: 6 }),
      expected: `100.000000b`,
    },
    {
      input: format(Amount("990000000000000000"), mockAsset(18), {
        mantissa: 6,
      }),
      expected: `0.990000`,
    },
    {
      input: format(Amount("990000000000000000"), mockAsset(18), {
        mantissa: 6,
        trimMantissa: true,
      }),
      expected: `0.99`,
    },
    {
      input: format(Amount("999999800000000000"), mockAsset(18), {
        mantissa: 7,
      }),

      expected: `0.9999998`,
    },
    {
      input: format(Amount("0.01"), { mode: "percent", mantissa: 1 }),
      expected: `1.0%`,
    },
    {
      input: format(Amount("0.1"), { mode: "percent", mantissa: 2 }),
      expected: `10.00%`,
    },
    {
      input: format(Amount("0.12345"), { mode: "percent", mantissa: 3 }),
      expected: `12.345%`,
    },
    {
      input: format(Amount(".12345"), {
        mode: "percent",
        mantissa: 3,
        space: true,
      }),
      expected: `12.345 %`,
    },
    {
      input: format(Amount("-990000000000000000"), mockAsset(18), {
        mantissa: 6,
        trimMantissa: true,
      }),
      expected: `-0.99`,
    },
    {
      input: format(
        Amount("999999800000000000"),
        mockAsset(18),

        { mantissa: 7, shorthand: true, forceSign: true },
      ),
      expected: `+0.9999998`,
    },
    {
      input: format(Amount("100000000000000000000"), mockAsset(18), {
        mantissa: 18,
      }),
      expected: `100.000000000000000000`,
    },
  ];

  tests.forEach(({ only, skip, input, expected }) => {
    const tester = only ? test.only : skip ? test.skip : test;
    tester(expected, () => expect(input).toBe(expected));
  });

  test("float mode", () => {
    expect(
      format(Amount("100").divide(Amount("3")), { float: true, mantissa: 18 }),
    ).toBe("33.333333333333333333"); // Precision loss because of numbro - should we remove it?
  });
});
