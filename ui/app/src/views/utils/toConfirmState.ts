import { TransactionStatus } from "ui-core";
import { ConfirmState } from "../../types";

// So this acts as an adapter and maps between states feel free to change the ConfirmState to match this
// for now we can convert between them like this
// this is a stopgap but better done in vue because the ConfirmState
// is not something core knows or cares about and it doesn't quite
// describe the general nature of a broadcast transaction state
// TODO: align these states based on TransactionStatus
export function toConfirmState(tx: TransactionStatus["state"]): ConfirmState {
  return {
    requested: "signing" as const,
    accepted: "confirmed" as const,
    rejected: "rejected" as const,
    failed: "failed" as const,
    out_of_gas: "out_of_gas" as const,
  }[tx];
}
