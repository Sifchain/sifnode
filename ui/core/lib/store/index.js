var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
import JSBI from "jsbi";
import { action, computed, observable } from "mobx";
import { MARKETCAP_TOKEN_ORDER } from "../constants";
import * as TOKENS from "../constants/tokens";
import { AssetAmount } from "../entities";
function getTokenBySymbol(symbol) {
    const tokenStore = TOKENS;
    return tokenStore[symbol];
}
// This is the reactive store that is shared with our frontend
// Trying to keep this flat
export class State {
    constructor(o) {
        // constants
        this.marketcapTokenOrder = MARKETCAP_TOKEN_ORDER;
        // reactive props
        this.userBalances = new Map();
        Object.assign(this, o);
    }
    get availableAssetAccounts() {
        const ordered = [];
        this.userBalances.forEach((balance) => {
            ordered.push(balance);
        });
        for (let i = 0; i < Math.min(this.marketcapTokenOrder.length, 20); i++) {
            const symbol = this.marketcapTokenOrder[i];
            const token = this.userBalances.get(symbol);
            if (!token)
                ordered.push(AssetAmount.create(getTokenBySymbol(symbol), JSBI.BigInt(0)));
        }
        return ordered;
    }
}
__decorate([
    observable
], State.prototype, "userBalances", void 0);
__decorate([
    computed
], State.prototype, "availableAssetAccounts", null);
// This is a bag of functions to mutate our state
// Lets break this up when it gets too big
export class StoreActions {
    constructor(state) {
        this.state = state;
    }
    setUserBalances(balances) {
        this.state.userBalances = balances.reduce((map, balance) => {
            map.set(balance.asset.symbol, balance);
            return map;
        }, new Map());
    }
}
__decorate([
    action.bound
], StoreActions.prototype, "setUserBalances", null);
export const store = createStore();
// Covenience function for creating the global store
export function createStore(state) {
    return new StoreActions(new State(state));
}
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
//# sourceMappingURL=index.js.map