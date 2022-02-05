import { boolean } from "fp-ts"
import { Exception } from "handlebars"
import { type } from "os"
import { SifEvent } from "../../src/watcher/watcher"
import { buildFailure, State, TransactionStep, VerbosityLevel } from "./context"

// Currently is a type alias, we're doing this because we might have complex verification step
type VerificationState = TransactionStep

export class StateMachineVerifierBuilder {
  private expectedStatesQueue: VerificationState[]
  private verboseLevel: VerbosityLevel
  constructor() {
    this.expectedStatesQueue = new Array()
    this.verboseLevel = "none"
  }

  // TODO: To be implemented
  setVerboseLevel(level: VerbosityLevel): void {
    this.verboseLevel = level
  }

  // Do we need this sugar?
  initial(step: TransactionStep): StateMachineVerifierBuilder {
    this.expectedStatesQueue.push(step)
    return this
  }

  // This can be overloaded to take simple step or complex step
  then(step: TransactionStep): StateMachineVerifierBuilder {
    this.expectedStatesQueue.push(step)
    return this
  }
  // ORRRRR... ps. merge this with .then by making assertionFn optional
  thenAndAssert(step: TransactionStep, assertionFn: (event: SifEvent) => boolean) {
    this.expectedStatesQueue.push(step)
    return this
  }

  stateEventAssertion(step: TransactionStep, assertionFn: (event: SifEvent) => boolean) {}

  // Do we need this sugar?
  finally(step: TransactionStep): StateMachineVerifierBuilder {
    return this
  }

  // TODO: Incomplete
  build(): StateMachineVerifier {
    return new StateMachineVerifier(this.expectedStatesQueue, this.verboseLevel)
  }
}

// TODO: Can we avoid maintaining this entirely with better typing?
// TODO: Where should we put this? Inside the class?
const ebrelayerstateToTxStep: Map<string, TransactionStep> = new Map()
ebrelayerstateToTxStep.set("ReceiveCosmosBurnMessage", TransactionStep.ReceiveCosmosBurnMessage)
ebrelayerstateToTxStep.set("WitnessSignProphecy", TransactionStep.WitnessSignProphecy)
ebrelayerstateToTxStep.set("ProphecyClaimSubmitted", TransactionStep.ProphecyClaimSubmitted)

const peggyEventToTxStep: Map<string, TransactionStep> = new Map()
peggyEventToTxStep.set("Burn", TransactionStep.Burn)
peggyEventToTxStep.set("GetCrossChainFeeConfig", TransactionStep.GetCrossChainFeeConfig)
peggyEventToTxStep.set("SendCoinsFromAccountToModule", TransactionStep.SendCoinsFromAccountToModule)
peggyEventToTxStep.set("BurnCoins", TransactionStep.BurnCoins)
peggyEventToTxStep.set("SetProphecy", TransactionStep.SetProphecy)
peggyEventToTxStep.set("PublishCosmosBurnMessage", TransactionStep.PublishCosmosBurnMessage)
peggyEventToTxStep.set("SetWitnessLockBurnNonce", TransactionStep.SetWitnessLockBurnNonce)
// TODO: ProphecyStatus is TERRIBLE NAMING!
peggyEventToTxStep.set("ProphecyStatus", TransactionStep.ProphecyStatus)

export class StateMachineVerifier {
  private verboseLevel: VerbosityLevel

  // Need function to convert
  private currentState: State = {
    value: { kind: "initialState" },
    createdAt: 0,
    currentHeartbeat: 0,
    transactionStep: TransactionStep.Initial,
    uniqueId: "eth to ceth",
  } as State

  // private terminalState: TransactionStep

  private expectedStateQueue: VerificationState[]

  // Option: Store expected state as map currentState -> expectedNextState
  // transactionStep -> assertionFunctions
  // private assertionFunctions: Map<TransactionStep, (event: SifEvent) => boolean>

  constructor(expectedStateQueue: VerificationState[], verboseLevel: VerbosityLevel) {
    this.expectedStateQueue = expectedStateQueue
    this.verboseLevel = verboseLevel
  }

  /**
   * Invoked in rx loop,
   *  takes in an event,
   *    verify it is right step,
   *    verify it has right values set
   *  outputs a context.State
   * @param event
   * @returns
   */
  verify(event: SifEvent): State {
    const currentTransactionStep = this.sifEventToTransactionStep(event)

    if (currentTransactionStep != this.expectedStateQueue[0]) {
      // Return failure state
      return buildFailure(
        this.currentState,
        event,
        `Bad Transition: Expected ${this.expectedStateQueue[0]}, Got ${currentTransactionStep}`
        // TODO: Support printing out the entire tree by storing successful steps
      )
    }
    // TODO: Lookup transaction step for extra verification

    // All verifications for this txStep was completed, moving for next iteration
    this.expectedStateQueue.shift()

    if (this.expectedStateQueue.length == 0) {
      // We have completed verification of all steps. We'll emit a successful state
      return {
        ...this.currentState,
        transactionStep: currentTransactionStep,
        value: {
          kind: "success",
        },
      }
    }
    return {
      ...this.currentState,
      transactionStep: currentTransactionStep,
    }
  }

  // This function would encapsulate all the case statement
  private sifEventToTransactionStep(event: SifEvent): TransactionStep {
    if (event.kind == "EbRelayerEvmStateTransition") {
      /**
       * TODO: We use "any" here because the it is object type, and we
       * need "kind" field
       */
      const eventData: any = event.data
      return ebrelayerstateToTxStep.get(eventData.kind)!
    }

    if (event.kind == "SifnodedPeggyEvent") {
      const eventData: any = event.data
      return peggyEventToTxStep.get(eventData.kind)!
    }

    if (event.kind == "EthereumMainnetLogUnlock") {
      return TransactionStep.EthereumMainnetLogUnlock
    }
    console.log("We SHOULD NOT HAVE REACHED HERE. Encountered unmapped sifEvent")
    throw Exception
    // return undefined
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

/**
 * What is a VerificationState?
 * Simple verification step:
 *  Dimension 1. At most once, at least once, exactly once
 *  Dimension 2. TrnasactionStep
 *  Dimension 3. Additional verification on the body itself
 *
 * Complex verification step:
 *  Composition of AND/OR of multiple simple verification steps
 *
 * StateMachineBuilder .then and .finally takes a VerificationState
 *
 */
