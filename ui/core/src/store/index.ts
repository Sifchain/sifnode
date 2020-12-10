import { reactive } from "@vue/reactivity";
import { wallet, WalletStore } from "./wallet";
import { asset, AssetStore } from "./asset";
import { pools, PoolStore } from "./pools";

export type Store = {
  wallet: WalletStore;
  asset: AssetStore;
  pools: PoolStore;
};

export function createStore() {
  const state = reactive<Store>({
    wallet,
    asset,
    pools,
  }) as Store;

  return state;
}

export type WithStore<T extends keyof Store = keyof Store> = {
  store: Pick<Store, T>;
};

//
// ==============================================================================
//
// For reference here is Uniswaps redux store shape:
//
// NOTE: Uniswap attempt to reuse their redux state for both pool and swap
//       not sure why exactly yet
//
// const state = {
//   application: {
//     blockNumber: {
//       "1": 11000440,
//     },
//     popupList: [],
//     openModal: null,
//   },
//   user: {
//     userDarkMode: null,
//     matchesDarkMode: false,
//     userExpertMode: false,
//     userSlippageTolerance: 50,
//     userDeadline: 1200,
//     tokens: {},
//     pairs: {},
//     timestamp: 1601963554327,
//     URLWarningVisible: true,
//     lastUpdateVersionTimestamp: 1601963554074,
//   },
//   transactions: {},
//   swap: {
//     INPUT: {
//       currencyId: "ETH",
//     },
//     OUTPUT: {
//       currencyId: "0xfC1E690f61EFd961294b3e1Ce3313fBD8aa4f85d",
//     },
//     independentField: "INPUT",
//     typedValue: "0.003",
//     recipient: null,
//   },
//   mint: {
//     independentField: "CURRENCY_A",
//     typedValue: "",
//     otherTypedValue: "",
//   },
//   burn: {
//     independentField: "LIQUIDITY_PERCENT",
//     typedValue: "0",
//   },
//   multicall: {
//     // serialized data from multicall
//     callResults: {
//       "1": {
//         "0x6C3e4cb2E96B01F4b866965A91ed4437839A121a-0x18160ddd": {
//           data:
//             "0x0000000000000000000000000000000000000000000000008505bfc91777ee70",
//           blockNumber: 11000528,
//         },
//       },
//       // ... lots of this stuff dont under stand the multicall stuff but possibly
//       // looks like this could be local blockchain data?
//     },
//   },
//   lists: {
//     selectedListUrl: "tokens.uniswap.eth",
//   },
// };
