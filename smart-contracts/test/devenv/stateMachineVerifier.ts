import { SifEvent } from "../../src/watcher/watcher"
import { State, TransactionStep } from "./test_lockburn"

class StateMachineVerifierBuilder {
  constructor() {}

  initial(step: TransactionStep): StateMachineVerifierBuilder {
    return this
  }

  then(step: TransactionStep): StateMachineVerifierBuilder {
    return this
  }

  finally(step: TransactionStep): StateMachineVerifierBuilder {
    return this
  }

  // TODO: Incomplete
  build(): StateMachineVerifier {
    return null
  }
}

class StateMachineVerifier {
  // Need function to convert
  private currentState: State

  // Option: Store expected state as map currentState -> expectedNextState

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
