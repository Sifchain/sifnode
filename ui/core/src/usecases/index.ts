import { Services, WithService } from "../services";
import { Store, WithStore } from "../store";
import ethWalletActions from "./ethWallet";
import clpActions from "./clp";
import walletActions from "./wallet";
import pegActions from "./peg";

export type UsecaseContext<
  T extends keyof Services = keyof Services,
  S extends keyof Store = keyof Store
> = WithService<T> & WithStore<S>;

export function createUsecases(context: UsecaseContext) {
  return {
    ethWallet: ethWalletActions(context),
    clp: clpActions(context),
    wallet: walletActions(context),
    peg: pegActions(context),
  };
}

export type Actions = ReturnType<typeof createUsecases>;
