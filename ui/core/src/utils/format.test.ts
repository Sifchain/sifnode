import { AssetAmount, Network } from "../entities";
import { Amount } from "../entities/Amount";
import { Asset } from "../entities/Asset";
import { round, format, IFormatOptions } from "./format";

type Test = {
  only?: boolean;
  skip?: boolean;
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

describe("round", () => {
  test("rounding", () => {
    expect(round("1.23456789", 4)).toBe("1.2346");
  });
});

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
      input: format(Amount("1.000000"), { mantissa: 6, trimMantissa: true }),
      expected: `1.0`,
    },
    {
      input: format(Amount("1.100000"), { mantissa: 6, trimMantissa: true }),
      expected: `1.1`,
    },
    {
      input: format(Amount("1.12300000"), { mantissa: 6, trimMantissa: true }),
      expected: `1.123`,
    },

    {
      input: format(Amount("1.1234567800000"), {
        mantissa: 6,
        trimMantissa: true,
      }),
      expected: `1.123457`,
    },
    {
      input: format(Amount("0"), {
        mantissa: 6,
        zeroFormat: "N/A",
        trimMantissa: true,
      }),
      expected: `N/A`,
    },
    {
      input: format(Amount("100000000000000000000"), mockAsset(18), {
        mantissa: 18,
      }),
      expected: `100.000000000000000000`,
    },
    // Dynamic mantissa
    // Will adjust based on given hash map
    // We should tokenize these as reusable options
    ...(() => {
      const dynamicMantissa = {
        1: 6,
        1000: 4,
        100000: 2,
        infinity: 0,
      };
      return [
        {
          input: format(Amount("0.12345678"), {
            mantissa: dynamicMantissa,
          }),
          expected: "0.123457",
        },
        {
          input: format(Amount("100.12345678"), {
            mantissa: dynamicMantissa,
          }),
          expected: "100.1235",
        },
        {
          input: format(Amount("5000.12345678"), {
            mantissa: dynamicMantissa,
          }),
          expected: "5000.12",
        },
        {
          input: format(Amount("500000.12345678"), {
            mantissa: dynamicMantissa,
          }),
          expected: "500000",
        },
      ];
    })(),
  ];

  tests.forEach(({ only, skip, input, expected }) => {
    const tester = only ? test.only : skip ? test.skip : test;
    tester(expected, () => expect(input).toBe(expected));
  });

  test("float mode", () => {
    expect(format(Amount("100").divide(Amount("3")), { mantissa: 18 })).toBe(
      "33.333333333333333333",
    );
  });

  test("does not throw on undefined and null inputs", () => {
    // because we are not using JSX there is a chance that we accidentally send undefined or null to format

    expect(() => {
      format(undefined as any, mockAsset(18));
    }).not.toThrow();

    expect(() => {
      format(null as any, mockAsset(18));
    }).not.toThrow();

    expect(() => {
      format(Amount("10"), undefined as any);
    }).not.toThrow();

    expect(() => {
      format(null as any, undefined as any);
    }).not.toThrow();
  });
});
