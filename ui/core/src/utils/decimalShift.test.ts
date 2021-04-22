import { decimalShift, floorDecimal } from "./decimalShift";

type Tests<InputType, OutputType = string> = {
  input: InputType;
  expected: OutputType;
  only?: boolean;
  skip?: boolean;
}[];

describe("decimalShift", () => {
  const decimalShiftTests: Tests<[string, number]> = [
    {
      input: ["-.1234", 2],
      expected: "-12.34",
    },
    {
      input: ["-0.1234", 2],
      expected: "-12.34",
    },
    {
      input: [".1234", 2],
      expected: "12.34",
    },
    {
      input: [".1234", 0],
      expected: "0.1234",
    },
    {
      input: ["123.4", -5],
      expected: "0.001234",
    },
    {
      input: ["12.34", -4],
      expected: "0.001234",
    },
    {
      input: ["12.34", -3],
      expected: "0.01234",
    },
    {
      input: ["12.34", -2],
      expected: "0.1234",
    },
    {
      input: ["12.34", -1],
      expected: "1.234",
    },
    {
      input: ["12.34", 0],
      expected: "12.34",
    },
    {
      input: ["12.34", 1],
      expected: "123.4",
    },
    {
      input: ["0012.34", 2],
      expected: "1234",
    },
    {
      input: ["12.34", 2],
      expected: "1234",
    },
    {
      input: ["12.34", 3],
      expected: "12340",
    },
    {
      input: ["12.34", 4],
      expected: "123400",
    },

    {
      input: ["123456789", 0],
      expected: "123456789",
    },
    {
      input: ["123456789", -2],
      expected: "1234567.89",
    },
    {
      input: ["123456789", 2],
      expected: "12345678900",
    },
    {
      input: ["12345678910", -10],
      expected: "1.2345678910",
    },
    {
      input: ["12345678910", -18],
      expected: "0.000000012345678910",
    },
  ];

  decimalShiftTests.forEach(
    ({ skip, only, input: [decimal, shift], expected }) => {
      const tester = only ? test.only : skip ? test.skip : test;

      tester(`${[decimal, shift].join("\t")}\t=> ${expected}`, () => {
        expect(decimalShift(decimal, shift)).toBe(expected);
      });
    },
  );
});

describe("floorDecimal", () => {
  const floorDecimalTests: Tests<string> = [
    {
      input: "1234.12341234",
      expected: "1234",
    },
    {
      input: "0.123412341234",
      expected: "0",
    },
    {
      input: "0.99999",
      expected: "0",
    },
  ];

  floorDecimalTests.forEach(({ skip, only, input, expected }) => {
    const tester = only ? test.only : skip ? test.skip : test;

    tester(`${input}\t=> ${expected}`, () => {
      expect(floorDecimal(input)).toBe(expected);
    });
  });
});
