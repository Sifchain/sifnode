import { boolean } from "fp-ts"
import { SifEvent } from "../../src/watcher/watcher"
import { State, TransactionStep } from "./test_lockburn"

export class StateMachineVerifierBuilder {
  constructor() {}

  setVerboseLevel(level: string): void {}

  initial(step: TransactionStep): StateMachineVerifierBuilder {
    return this
  }

  then(step: TransactionStep): StateMachineVerifierBuilder {
    return this
  }

  stateEventAssertion(step: TransactionStep, assertionFn: (event: SifEvent) => boolean) {}

  finally(step: TransactionStep): StateMachineVerifierBuilder {
    return this
  }

  // TODO: Incomplete
  build(): StateMachineVerifier {
    return null
  }
}

export class StateMachineVerifier {
  // Need function to convert
  private currentState: State
  private terminalState: TransactionStep

  // Option: Store expected state as map currentState -> expectedNextState

  // transactionStep -> assertionFunctions
  private assertionFunctions: Map<TransactionStep, (event: SifEvent) => boolean>

  verify(event: SifEvent): State {
    const currentTransactionStep = this.sifEventToTransactionStep(event)

    // TODO: Lookup transaction step for extra verification

    if (currentTransactionStep == this.terminalState) {
      return {
        ...this.currentState,
        value: {
          kind: "terminate",
        },
      }
    }
    return null
  }

  // This function would encapsulate all the case statement
  private sifEventToTransactionStep(event: SifEvent): TransactionStep {
    return null
  }
}

// class State {}
/**
 * Demo
 *
 * const stateMachienBuilder = new StateMachineVerifierBuilder()
 * stateMachineBuilder.then(TransactionStep.Prophecy)
 *                    .then(TransactionStep.EthereumMainnetLogUnlock)
 *                    .assert((sifevent) -> {})
 *                    .assert()
 *                   .then(TransactionStep.fdfs)
 *                      .assert(//fdfs is good)
 *                    .finally(TransactionStep.ProphecyClaimSubmitted)
 *
 * let verifier: StateMachineVerifier = stateMachienBuilder.build()
 */
