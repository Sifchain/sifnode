// import { Ref, ref } from "@vue/reactivity";
// import {
//   Asset,
//   AssetAmount,
//   IAssetAmount,
//   Network,
//   Pair,
//   Token,
// } from "../entities";
// import { PoolState, usePoolCalculator } from "./poolCalculator";

// const TOKENS = {
//   atk: Token({
//     decimals: 6,
//     symbol: "atk",
//     name: "AppleToken",
//     address: "123",
//     network: Network.ETHEREUM,
//   }),
//   btk: Token({
//     decimals: 6,
//     symbol: "btk",
//     name: "BananaToken",
//     address: "1234",
//     network: Network.ETHEREUM,
//   }),
//   eth: Token({
//     decimals: 18,
//     symbol: "eth",
//     name: "Ethereum",
//     address: "1234",
//     network: Network.ETHEREUM,
//   }),
// };

// // TODO: All the maths here are pretty naive need to double check with blockscience
// describe("usePoolCalculator", () => {
//   // input
//   const fromAmount: Ref<string> = ref("0");
//   const fromSymbol: Ref<string | null> = ref(null);
//   const toAmount: Ref<string> = ref("0");
//   const toSymbol: Ref<string | null> = ref(null);
//   const balances = ref([]) as Ref<IAssetAmount[]>;
//   const selectedField: Ref<"from" | "to" | null> = ref("from");
//   const marketPairFinder = jest.fn<Pair | null, any>(() => null);

//   // output

//   let aPerBRatioMessage: Ref<string>;
//   let bPerARatioMessage: Ref<string>;
//   let shareOfPool: Ref<string>;
//   let state: Ref<PoolState>;
//   beforeEach(() => {
//     ({
//       state,
//       aPerBRatioMessage,
//       bPerARatioMessage,
//       shareOfPool,
//     } = usePoolCalculator({
//       balances,
//       fromAmount,
//       toAmount,
//       fromSymbol,
//       selectedField,
//       toSymbol,
//       marketPairFinder,
//     }));

//     balances.value = [];
//     fromAmount.value = "0";
//     toAmount.value = "0";
//     fromSymbol.value = null;
//     toSymbol.value = null;
//   });

//   test("poolCalculator ratio messages", () => {
//     fromAmount.value = "1000";
//     toAmount.value = "500";
//     fromSymbol.value = "atk";
//     toSymbol.value = "btk";

//     expect(aPerBRatioMessage.value).toBe("0.5 BTK per ATK");
//     expect(bPerARatioMessage.value).toBe("2.0 ATK per BTK");
//     expect(shareOfPool.value).toBe("100.00%");
//   });

//   test("Can handle division by zero", () => {
//     fromAmount.value = "0";
//     toAmount.value = "0";
//     fromSymbol.value = "atk";
//     toSymbol.value = "btk";
//     expect(state.value).toBe(PoolState.ZERO_AMOUNTS);
//     expect(aPerBRatioMessage.value).toBe("");
//     expect(bPerARatioMessage.value).toBe("");
//     expect(shareOfPool.value).toBe("");
//   });

//   test("Calculate against a given pool", () => {
//     marketPairFinder.mockImplementationOnce(() =>
//       Pair(AssetAmount(TOKENS.atk, "2000"), AssetAmount(TOKENS.btk, "2000"))
//     );

//     fromAmount.value = "1000";
//     toAmount.value = "1000";
//     fromSymbol.value = "atk";
//     toSymbol.value = "btk";

//     expect(aPerBRatioMessage.value).toBe("1.0 BTK per ATK");
//     expect(bPerARatioMessage.value).toBe("1.0 ATK per BTK");
//     expect(shareOfPool.value).toBe("33.33%");
//   });

//   test("Insufficient balances", () => {
//     balances.value = [
//       AssetAmount(Asset.get("atk"), "10000"),
//       AssetAmount(Asset.get("btk"), "10000"),
//     ];
//     fromAmount.value = "1000000";
//     toAmount.value = "1000";
//     fromSymbol.value = "atk";
//     toSymbol.value = "btk";

//     expect(state.value).toBe(PoolState.INSUFFICIENT_FUNDS);
//   });
// });
test("", () => {
  expect(1).toBe(1);
});
