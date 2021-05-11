import { Services, WithService } from "../services";
import { Store, WithStore } from "../store";
import ethWalletActions from "./wallet/eth";
import clpActions from "./clp";
import walletActions from "./wallet/sif";
import pegActions from "./peg";

export type UsecaseContext<
  T extends keyof Services = keyof Services,
  S extends keyof Store = keyof Store
> = WithService<T> & WithStore<S>;

export function createUsecases(context: UsecaseContext) {
  return {
    clp: clpActions(context),
    wallet: {
      sif: walletActions(context),
      eth: ethWalletActions(context),
    },
    peg: pegActions(context),
  };
}

export type Actions = ReturnType<typeof createUsecases>;
