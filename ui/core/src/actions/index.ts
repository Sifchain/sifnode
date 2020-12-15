import { Api, WithApi } from "../api";
import { Store, WithStore } from "../store";
import ethWalletActions from "./ethWallet";
import clpActions from "./clp";
import walletActions from "./wallet";

export type ActionContext<
  T extends keyof Api = keyof Api,
  S extends keyof Store = keyof Store
> = WithApi<T> & WithStore<S>;

export function createActions(context: ActionContext) {
  return {
    ethWallet: ethWalletActions(context),
    clp: clpActions(context),
    wallet: walletActions(context),
  };
}

export type Actions = ReturnType<typeof createActions>;
