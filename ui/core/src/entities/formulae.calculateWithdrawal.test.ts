import { calculateWithdrawal } from "./formulae";
import { Fraction } from "./fraction/Fraction";

const tests = [
  {
    skip: false,
    only: false,
    name: "even",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "0",
      lpUnits: "10000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "100000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "10000000000",
      withdrawNativeAssetAmount: "10000000000",
    },
  },
  {
    name: "all external",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "10000",
      lpUnits: "10000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "100000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "18100000000",
      withdrawNativeAssetAmount: "0",
    },
  },
  {
    name: "all native",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "-10000",
      lpUnits: "1000000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "100000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "0",
      withdrawNativeAssetAmount: "18100000000",
    },
  },
  {
    name: "25% native",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "-5000",
      lpUnits: "10000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "100000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "5000000000",
      withdrawNativeAssetAmount: "14487500000",
    },
  },
  {
    name: "25% external",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "5000",
      lpUnits: "10000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "10000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "14487500000",
      withdrawNativeAssetAmount: "5000000000",
    },
  },
  {
    name: "external worth half",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "0",
      lpUnits: "10000000000",
      externalAssetBalance: "2000000000000",
      nativeAssetBalance: "100000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "20000000000",
      withdrawNativeAssetAmount: "10000000000",
    },
  },
  {
    name: "Calculates withdrawal correctly",
    input: {
      poolUnits: "100009900000",
      wBasisPoints: "10000",
      asymmetry: "10000",
      lpUnits: "10000000000",
      externalAssetBalance: "101000000000",
      nativeAssetBalance: "99019800000",
    },
    expected: {
      withdrawExternalAssetAmount: "18279400000",
      withdrawNativeAssetAmount: "0",
    },
  },
  {
    name: "Calculates asymmetry -10000",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "-10000",
      lpUnits: "10000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "50000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "0",
      withdrawNativeAssetAmount: "9050000000",
    },
  },
  {
    name: "Calculates asymmetry -5000",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "-5000",
      lpUnits: "10000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "50000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "5000000000",
      withdrawNativeAssetAmount: "7243800000",
    },
  },
  {
    name: "Calculates asymmetry 0",
    input: {
      poolUnits: "100000000000",
      wBasisPoints: "10000",
      asymmetry: "0",
      lpUnits: "10000000000",
      externalAssetBalance: "100000000000",
      nativeAssetBalance: "50000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "10000000000",
      withdrawNativeAssetAmount: "5000000000",
    },
  },
  {
    name: "Calculates asymmetry 5000",
    input: {
      poolUnits: "10000000000",
      wBasisPoints: "10000",
      asymmetry: "5000",
      lpUnits: "10000000000",
      externalAssetBalance: "10000000000",
      nativeAssetBalance: "5000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "1448750000",
      withdrawNativeAssetAmount: "250000000",
    },
  },
  {
    name: "Calculates asymmetry 10000",
    input: {
      poolUnits: "10000000000",
      wBasisPoints: "10000",
      asymmetry: "10000",
      lpUnits: "1000000000",
      externalAssetBalance: "10000000000",
      nativeAssetBalance: "5000000000",
    },
    expected: {
      withdrawExternalAssetAmount: "1810000000",
      withdrawNativeAssetAmount: "0",
    },
  },
  {
    name: "Even: Calculates asymmetry 5000",
    input: {
      asymmetry: "5000",
      externalAssetBalance: "1000000000",
      lpUnits: "1000000000",
      nativeAssetBalance: "1000000000",
      poolUnits: "1000000000",
      wBasisPoints: "10000",
    },
    expected: {
      withdrawExternalAssetAmount: "1000000000",
      withdrawNativeAssetAmount: "500000000",
    },
  },
  {
    name: "Even: Calculates asymmetry 10000",
    input: {
      asymmetry: "10000",
      externalAssetBalance: "1000000000",
      lpUnits: "1000000000",
      nativeAssetBalance: "1000000000",
      poolUnits: "1000000000",
      wBasisPoints: "10000",
    },
    expected: {
      withdrawExternalAssetAmount: "1000000000",
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
