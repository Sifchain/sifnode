import { boolean } from "fp-ts"
import { SifEvent } from "../../src/watcher/watcher"
import { State, TransactionStep } from "./test_lockburn"

export class StateMachineVerifierBuilder {
  constructor() {}

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

  // Option: Store expected state as map currentState -> expectedNextState

  // transactionStep -> assertionFunctions

  verify(event: SifEvent): State {
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
 *                    .finally(TransactionStep.ProphecyClaimSubmitted)
 */
