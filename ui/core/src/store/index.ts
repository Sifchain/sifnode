import JSBI from "jsbi";
import * as TOKENS from "../constants/tokens";
import { AssetAmount, Asset } from "../entities";
import { reactive, computed } from "@vue/reactivity";

function getTokenBySymbol(symbol: string) {
  const tokenStore = TOKENS as { [symbol: string]: Asset };
  return tokenStore[symbol];
}

type AssetAmountMap = Map<string, AssetAmount>;

export type State = {
  userBalances: AssetAmountMap;
  marketcapTokenOrder: string[];
  availableAssetAccounts: readonly AssetAmount[];
};

export type Actions = {
  setUserBalances: (balances: AssetAmount[]) => void;
};

export function createStore(initialState?: Partial<State>) {
  const availableAssetAccounts = (computed<AssetAmount[]>(() => {
    const ordered: AssetAmount[] = [];
    for (const balance of state.userBalances) {
      ordered.push(balance[1]);
    }

    for (let i = 0; i < Math.min(state.marketcapTokenOrder.length, 20); i++) {
      const symbol = state.marketcapTokenOrder[i];
      const token = state.userBalances.get(symbol);
      const tokenSymbol = getTokenBySymbol(symbol);
      if (!token && tokenSymbol) {
        ordered.push(
          AssetAmount.create(getTokenBySymbol(symbol), JSBI.BigInt(0))
        );
      }
    }

    return ordered;
  }) as unknown) as AssetAmount[];

  const state = reactive<State>({
    marketcapTokenOrder: [],
    userBalances: new Map(),
    availableAssetAccounts,
    ...initialState,
  }) as State;

  const setUserBalances = (balances: AssetAmount[]) => {
    state.userBalances = balances.reduce((map, balance) => {
      map.set(balance.asset.symbol, balance);
      return map;
    }, new Map<string, AssetAmount>());
  };

  return {
    state,
    setUserBalances,
  };
}
export type Store = ReturnType<typeof createStore>;

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
