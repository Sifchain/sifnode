import { calculateWithdrawal } from "./formulae";
import { Fraction } from "./fraction/Fraction";

const tests = [
  {
    skip: false,
    only: false,
    name: "even",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "0",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "1000000",
    },
    expected: {
      withdrawExternalAssetAmount: "100000",
      withdrawNativeAssetAmount: "100000",
    },
  },
  {
    name: "all external",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "10000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "1000000",
    },
    expected: {
      withdrawExternalAssetAmount: "181000",
      withdrawNativeAssetAmount: "0",
    },
  },
  {
    name: "all native",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "-10000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "1000000",
    },
    expected: {
      withdrawExternalAssetAmount: "0",
      withdrawNativeAssetAmount: "181000",
    },
  },
  {
    name: "25% native",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "-5000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "1000000",
    },
    expected: {
      withdrawExternalAssetAmount: "50000",
      withdrawNativeAssetAmount: "144875",
    },
  },
  {
    name: "25% external",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "5000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "1000000",
    },
    expected: {
      withdrawExternalAssetAmount: "144875",
      withdrawNativeAssetAmount: "50000",
    },
  },
  {
    name: "external worth half",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "0",
      lpUnits: "100000",
      externalAssetBalance: "2000000",
      nativeAssetBalance: "1000000",
    },
    expected: {
      withdrawExternalAssetAmount: "200000",
      withdrawNativeAssetAmount: "100000",
    },
  },
  {
    name: "Calculates withdrawal correctly",
    input: {
      poolUnits: "1000099",
      wBasisPoints: "10000",
      asymmetry: "10000",
      lpUnits: "100000",
      externalAssetBalance: "1010000",
      nativeAssetBalance: "990198",
    },
    expected: {
      withdrawExternalAssetAmount: "179927",
      withdrawNativeAssetAmount: "0",
    },
  },
  {
    name: "Calculates withdrawal correctly",
    input: {
      poolUnits: "1000099",
      wBasisPoints: "10000",
      asymmetry: "6000",
      lpUnits: "100000",
      externalAssetBalance: "1010000",
      nativeAssetBalance: "990198",
    },
    expected: {
      withdrawExternalAssetAmount: "152305",
      withdrawNativeAssetAmount: "39604",
    },
  },
  {
    name: "Calculates withdrawal correctly",
    input: {
      asymmetry: "708",
      externalAssetBalance: "1010000",
      lpUnits: "1000000",
      nativeAssetBalance: "990198",
      poolUnits: "1000099",
      wBasisPoints: "10000",
    },
    expected: {
      withdrawExternalAssetAmount: "1009900",
      withdrawNativeAssetAmount: "920001",
    },
  },
  {
    name: "Calculates withdrawal correctly",
    input: {
      asymmetry: "10000",
      externalAssetBalance: "1010000",
      lpUnits: "1000000",
      nativeAssetBalance: "990198",
      poolUnits: "1000099",
      wBasisPoints: "10000",
    },
    expected: {
      withdrawExternalAssetAmount: "1009900",
      withdrawNativeAssetAmount: "0",
    },
  },
];

tests.forEach(({ name, only, skip, input, expected }) => {
  const tester = only ? test.only : skip ? test.skip : test;

  tester(name, () => {
    const output = calculateWithdrawal({
      poolUnits: new Fraction(input.poolUnits),
      wBasisPoints: new Fraction(input.wBasisPoints),
      asymmetry: new Fraction(input.asymmetry),
      lpUnits: new Fraction(input.lpUnits),
      externalAssetBalance: new Fraction(input.externalAssetBalance),
      nativeAssetBalance: new Fraction(input.nativeAssetBalance),
    });
    expect(output.withdrawExternalAssetAmount.toFixed(0)).toEqual(
      expected.withdrawExternalAssetAmount
    );
    expect(output.withdrawNativeAssetAmount.toFixed(0)).toEqual(
      expected.withdrawNativeAssetAmount
    );
  });
});
