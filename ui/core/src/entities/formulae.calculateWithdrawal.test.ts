import { Amount } from "./Amount";
import { calculateWithdrawal } from "./formulae";

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
      withdrawExternalAssetAmount: "182794",
      withdrawNativeAssetAmount: "0",
    },
  },
  {
    name: "Calculates asymmetry -10000",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "-10000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "500000",
    },
    expected: {
      withdrawExternalAssetAmount: "0",
      withdrawNativeAssetAmount: "90500",
    },
  },
  {
    name: "Calculates asymmetry -5000",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "-5000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "500000",
    },
    expected: {
      withdrawExternalAssetAmount: "50000",
      withdrawNativeAssetAmount: "72438",
    },
  },
  {
    name: "Calculates asymmetry 0",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "0",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "500000",
    },
    expected: {
      withdrawExternalAssetAmount: "100000",
      withdrawNativeAssetAmount: "50000",
    },
  },
  {
    name: "Calculates asymmetry 5000",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "5000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "500000",
    },
    expected: {
      withdrawExternalAssetAmount: "144875",
      withdrawNativeAssetAmount: "25000",
    },
  },
  {
    name: "Calculates asymmetry 10000",
    input: {
      poolUnits: "1000000",
      wBasisPoints: "10000",
      asymmetry: "10000",
      lpUnits: "100000",
      externalAssetBalance: "1000000",
      nativeAssetBalance: "500000",
    },
    expected: {
      withdrawExternalAssetAmount: "181000",
      withdrawNativeAssetAmount: "0",
    },
  },
  {
    name: "Even: Calculates asymmetry 5000",
    input: {
      asymmetry: "5000",
      externalAssetBalance: "1000000",
      lpUnits: "1000000",
      nativeAssetBalance: "1000000",
      poolUnits: "1000000",
      wBasisPoints: "10000",
    },
    expected: {
      withdrawExternalAssetAmount: "1000000",
      withdrawNativeAssetAmount: "500000",
    },
  },
  {
    name: "Even: Calculates asymmetry 10000",
    input: {
      asymmetry: "10000",
      externalAssetBalance: "1000000",
      lpUnits: "1000000",
      nativeAssetBalance: "1000000",
      poolUnits: "1000000",
      wBasisPoints: "10000",
    },
    expected: {
      withdrawExternalAssetAmount: "1000000",
      withdrawNativeAssetAmount: "0",
    },
  },
];

tests.forEach(({ name, only, skip, input, expected }) => {
  const tester = only ? test.only : skip ? test.skip : test;

  tester(name, () => {
    const output = calculateWithdrawal({
      poolUnits: Amount(input.poolUnits),
      wBasisPoints: Amount(input.wBasisPoints),
      asymmetry: Amount(input.asymmetry),
      lpUnits: Amount(input.lpUnits),
      externalAssetBalance: Amount(input.externalAssetBalance),
      nativeAssetBalance: Amount(input.nativeAssetBalance),
    });
    expect(output.withdrawExternalAssetAmount.toBigInt().toString()).toEqual(
      expected.withdrawExternalAssetAmount,
    );
    expect(output.withdrawNativeAssetAmount.toBigInt().toString()).toEqual(
      expected.withdrawNativeAssetAmount,
    );
  });
});
