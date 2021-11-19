import chai, { expect } from "chai";
import { solidity } from "ethereum-waffle";
import { container } from "tsyringe";
import { HardhatRuntimeEnvironmentToken } from "../../src/tsyringe/injectionTokens";
import * as hardhat from "hardhat";
import { BigNumber } from "ethers";
import {
  ethereumResultsToSifchainAccounts,
  readDevEnvObj,
} from "../../src/tsyringe/devenvUtilities";
import { SifchainContractFactories } from "../../src/tsyringe/contracts";
import { buildDevEnvContracts, DevEnvContracts } from "../../src/contractSupport";
import web3 from "web3";
import * as ethereumAddress from "../../src/ethereumAddress";
import { SifEvent, SifHeartbeat, sifwatch, sifwatchReplayable } from "../../src/watcher/watcher";
import * as rxjs from "rxjs";
import {
  defer,
  distinctUntilChanged,
  lastValueFrom,
  Observable,
  scan,
  Subscription,
  takeWhile,
} from "rxjs";
import { EbRelayerEvmEvent } from "../../src/watcher/ebrelayer";
import { EthereumMainnetEvent } from "../../src/watcher/ethereumMainnet";
import { filter } from "rxjs/operators";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import * as ChildProcess from "child_process";
import { crossChainBurnFee, crossChainFeeBase, EbRelayerAccount } from "../../src/devenv/sifnoded";
import { v4 as uuidv4 } from "uuid";
import * as dotenv from "dotenv";
import deepEqual = require("deep-equal");

// The hash value for ethereum on mainnet
const ethDenomHash = "sif5ebfaf95495ceb5a3efbd0b0c63150676ec71e023b1043c40bcaaf91c00e15b2";

chai.use(solidity);

interface Failure {
  kind: "failure";
  value: SifEvent | "timeout";
  message: string;
}

interface Success {
  kind: "success";
}

interface InitialState {
  kind: "initialState";
}

interface Terminate {
  kind: "terminate";
}

interface State {
  value: SifEvent | EthereumMainnetEvent | Success | Failure | InitialState | Terminate;
  createdAt: number;
  currentHeartbeat: number;
  fromEthereumAddress: string;
  ethereumNonce: BigNumber;
  denomHash: string;
  ethereumLockBurnSequence: BigNumber;
  transactionStep: TransactionStep;
  uniqueId: string;
}

enum TransactionStep {
  Initial = "Initial",
  SawLogLock = "SawLogLock",
  SawProphecyClaim = "SawProphecyClaim",
  SawEthbridgeClaimArray = "SawEthbridgeClaimArray",
  BroadcastTx = "BroadcastTx",
  EthBridgeClaimArray = "EthBridgeClaimArray",
  CreateEthBridgeClaim = "CreateEthBridgeClaim",
  AddTokenMetadata = "AddTokenMetadata",
  AppendValidatorToProphecy = "AppendValidatorToProphecy",
  ProcessSuccessfulClaim = "ProcessSuccessfulClaim",
  CoinsSent = "CoinsSent",

  Burn = "Burn",
  GetTokenMetadata = "GetTokenMetadata",
  CosmosEvent = "CosmosEvent",
  PublishedProphecy = "PublishedProphecy",
  LogBridgeTokenMint = "LogBridgeTokenMint",
}

function isTerminalState(s: State) {
  switch (s.value.kind) {
    case "success":
    case "failure":
      return true;
    default:
      return s.transactionStep === TransactionStep.CoinsSent;
  }
}

function isNotTerminalState(s: State) {
  return !isTerminalState(s);
}

type VerbosityLevel = "summary" | "full" | "none";

function verbosityLevel(): VerbosityLevel {
  switch (process.env["VERBOSE"]) {
    case undefined:
      return "none";
    case "summary":
      return "summary";
    default:
      return "full";
  }
}

function attachDebugPrintfs<T>(xs: Observable<T>, verbosity: VerbosityLevel): Subscription {
  return xs.subscribe({
    next: (x) => {
      switch (verbosity) {
        case "full":
          console.log(JSON.stringify(x));
          break;
        case "summary":
          const p = x as any;
          console.log(
            `${p.currentHeartbeat}\t${p.transactionStep}\t${p.value?.kind}\t${p.value?.data?.kind}`
          );
          break;
      }
    },
    error: (e) => console.log("goterror: ", e),
    complete: () => console.log("alldone"),
  });
}

function hasDuplicateNonce(a: EbRelayerEvmEvent, b: EbRelayerEvmEvent): boolean {
  return a.data.event.Nonce === b.data.event.Nonce;
}

const gobin = process.env["GOBIN"];

describe("lock and burn tests", () => {
  dotenv.config();

  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost");

  const devEnvObject = readDevEnvObj("environment.json");
  // a generic sif address, nothing special about it
  const recipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace");

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat });
  });

  function ensureCorrectTransition(
    acc: State,
    v: SifEvent,
    predecessor: TransactionStep | TransactionStep[],
    successor: TransactionStep
  ): State {
    var stepIsCorrect: boolean;
    if (Array.isArray(predecessor)) {
      stepIsCorrect = (predecessor as string[]).indexOf(acc.transactionStep) >= 0;
    } else {
      stepIsCorrect = predecessor === acc.transactionStep;
    }
    if (stepIsCorrect)
      return {
        ...acc,
        value: v,
        createdAt: acc.currentHeartbeat,
        transactionStep: successor,
      };
    else
      return buildFailure(
        acc,
        v,
        `bad transition: expected ${predecessor}, got ${acc.transactionStep} before transition to ${successor}`
      );
  }

  function buildFailure(acc: State, v: SifEvent, message: string): State {
    return {
      ...acc,
      value: {
        kind: "failure",
        value: v,
        message: message,
      },
    };
  }

  // TODO: This ISNT an ebrelayer Account. it is a SIFACCOUNT
  function createTestSifAccount(): EbRelayerAccount {
    // Generate uuid
    let testSifAccountName = uuidv4();
    let cmd: string = `${gobin}/sifnoded keys add ${testSifAccountName} --home ${
      devEnvObject!.sifResults!.adminAddress!.homeDir
    } --keyring-backend test --output json`;
    let responseString: string = ChildProcess.execSync(cmd, { encoding: "utf8" });
    let responseJson = JSON.parse(responseString);

    // console.log("CreateTestAccount Response: ", responseJson)
    return {
      name: responseJson.name,
      account: responseJson.address,
      homeDir: "",
    };
  }

  // TODO: Move all these sif TS SDK to it's own class. I think it should go to smart-contract/devenv
  // TODO: Rethink naming. SendToSif?
  function fundSifAccount(
    adminAccount: string,
    destination: string,
    amount: number,
    symbol: string,
    homeDir: string
  ): object {
    // sifnoded tx bank send adminAccount testAccountToBeFunded --keyring-backend test --chain-id localnet concat(amount,symbol) --gas-prices=0.5rowan --gas-adjustment=1.5 --home <homeDir> --gas auto -y
    let sifnodedCmd: string = `${gobin}/sifnoded tx bank send ${adminAccount} ${destination} --keyring-backend test --chain-id localnet ${amount}${symbol} --gas-prices=0.5rowan --gas-adjustment=1.5 --home ${homeDir} --gas auto -y`;
    let responseString: string = ChildProcess.execSync(sifnodedCmd, { encoding: "utf8" });
    return JSON.parse(responseString);
  }

  // TODO: This is placed here now because devObject is available in this scope
  async function sifTransfer(
    sender: string,
    destination: SignerWithAddress,
    amount: BigNumber,
    symbol: string,
    // TODO: What is correct value for corsschainfee?
    crossChainFee: string,
    netwrokDescriptor: number
  ) {}

  async function executeSifBurn(
    sender: EbRelayerAccount,
    destination: SignerWithAddress,
    amount: BigNumber,
    symbol: string,
    // TODO: What is correct value for corsschainfee?
    crossChainFee: string,
    netwrokDescriptor: number
  ): Promise<object> {
    let sifnodedCmd: string = `${gobin}/sifnoded tx ethbridge burn ${sender.account} ${
      destination.address
    } ${amount} ${symbol} ${crossChainFee} --network-descriptor ${netwrokDescriptor} --keyring-backend test --gas-prices=0.5rowan --gas-adjustment=1.5 --chain-id localnet --home ${
      devEnvObject!.sifResults!.adminAddress!.homeDir
    } --from ${sender.name} -y `;

    let responseString = ChildProcess.execSync(sifnodedCmd, { encoding: "utf8" });
    return JSON.parse(responseString);
  }

  // Wrap an async function into an Observable<T>
  function deferAsync<T>(fn: () => Promise<T>): Observable<T> {
    return defer(() => rxjs.from(fn()));
  }

  async function executeLock(
    contracts: DevEnvContracts,
    smallAmount: BigNumber,
    sender1: SignerWithAddress,
    sifchainRecipient: string,
    verbose: boolean,
    identifier: string
  ) {
    const [evmRelayerEvents, replayedEvents] = sifwatchReplayable(
      {
        evmrelayer: "/tmp/sifnode/evmrelayer.log",
        sifnoded: "/tmp/sifnode/sifnoded.log",
      },
      hardhat,
      contracts.bridgeBank
    );

    const tx = await contracts.bridgeBank
      .connect(sender1)
      .lock(sifchainRecipient, ethereumAddress.eth.address, smallAmount, {
        value: smallAmount,
      });

    const states: Observable<State> = evmRelayerEvents
      .pipe(filter((x) => x.kind !== "SifnodedInfoEvent"))
      .pipe(
        scan(
          (acc: State, v: SifEvent) => {
            if (isTerminalState(acc))
              // we've reached a decision
              return { ...acc, value: { kind: "terminate" } as Terminate };
            switch (v.kind) {
              case "EbRelayerError":
              case "SifnodedError":
                // if we get an actual error, that's always a failure
                return { ...acc, value: { kind: "failure", value: v, message: "simple error" } };
              case "SifHeartbeat":
                // we just store the heartbeat
                return { ...acc, currentHeartbeat: v.value } as State;
              case "EthereumMainnetLogLock":
                // we should see exactly one lock
                let ethBlock = v.data.block as any;
                if (ethBlock.transactionHash === tx.hash && v.data.value.eq(smallAmount)) {
                  const newAcc: State = {
                    ...acc,
                    fromEthereumAddress: v.data.from,
                    ethereumNonce: BigNumber.from(v.data.nonce),
                  };
                  return ensureCorrectTransition(
                    newAcc,
                    v,
                    TransactionStep.Initial,
                    TransactionStep.SawLogLock
                  );
                } else
                  return {
                    ...acc,
                    value: {
                      kind: "failure",
                      value: v,
                      message: "incorrect EthereumMainnetLogLock",
                    },
                  };
              case "EbRelayerEvmStateTransition":
                switch ((v.data as any).kind) {
                  case "EthereumProphecyClaim":
                    const d = v.data as any;
                    if (
                      d.prophecyClaim.ethereum_sender == acc.fromEthereumAddress &&
                      BigNumber.from(d.event.Nonce).eq(acc.ethereumNonce)
                    )
                      return ensureCorrectTransition(
                        {
                          ...acc,
                          denomHash: d.prophecyClaim.denom_hash,
                        },
                        v,
                        TransactionStep.SawLogLock,
                        TransactionStep.SawProphecyClaim
                      );
                    break;
                  case "EthBridgeClaimArray":
                    let claims = (v.data as any).claims as any[];
                    const matchingClaim = claims.find(
                      (claim) => claim.denom_hash === acc.denomHash
                    );
                    if (matchingClaim)
                      return ensureCorrectTransition(
                        acc,
                        v,
                        TransactionStep.SawProphecyClaim,
                        TransactionStep.EthBridgeClaimArray
                      );
                    break;
                  case "BroadcastTx":
                    const messages = (v.data as any).messages as any[];
                    const matchingMessage = messages.find(
                      (msg) => msg.eth_bridge_claim.denom_hash === acc.denomHash
                    );
                    if (matchingMessage)
                      return ensureCorrectTransition(
                        acc,
                        v,
                        TransactionStep.EthBridgeClaimArray,
                        TransactionStep.BroadcastTx
                      );
                }
              case "SifnodedPeggyEvent":
                switch ((v.data as any).kind) {
                  case "coinsSent":
                    const coins = ((v.data as any).coins as any)[0];
                    if (coins["denom"] === ethDenomHash && smallAmount.eq(coins["amount"]))
                      return ensureCorrectTransition(
                        acc,
                        v,
                        TransactionStep.ProcessSuccessfulClaim,
                        TransactionStep.CoinsSent
                      );
                    else return buildFailure(acc, v, "incorrect hash or amount");
                  // TODO these steps need validation to make sure they're happing in the right order with the right data
                  case "CreateEthBridgeClaim":
                    let newSequenceNumber = (v.data as any).msg.Interface.eth_bridge_claim
                      .ethereum_lock_burn_sequence;
                    if (acc.ethereumNonce?.eq(newSequenceNumber))
                      return ensureCorrectTransition(
                        acc,
                        v,
                        [TransactionStep.BroadcastTx, TransactionStep.AppendValidatorToProphecy],
                        TransactionStep.CreateEthBridgeClaim
                      );
                    break;
                  case "AppendValidatorToProphecy":
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.CreateEthBridgeClaim,
                      TransactionStep.AppendValidatorToProphecy
                    );
                  case "ProcessSuccessfulClaim":
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.AppendValidatorToProphecy,
                      TransactionStep.ProcessSuccessfulClaim
                    );
                  case "AddTokenMetadata":
                    return ensureCorrectTransition(
                      acc,
                      v,
                      TransactionStep.ProcessSuccessfulClaim,
                      TransactionStep.AddTokenMetadata
                    );
                }
                return { ...acc, value: v, createdAt: acc.currentHeartbeat };
              default:
                // we have a new value (of any kind) and it should use the current heartbeat as its creation time
                return { ...acc, value: v, createdAt: acc.currentHeartbeat };
            }
          },
          {
            value: { kind: "initialState" },
            createdAt: 0,
            currentHeartbeat: 0,
            transactionStep: TransactionStep.Initial,
            uniqueId: identifier,
          } as State
        )
      );

    // it's useful to skip debug prints of states where only the heartbeat changed
    const withoutHeartbeat = states.pipe(
      distinctUntilChanged<State>((a, b) => {
        return deepEqual({ ...a, currentHeartbeat: 0 }, { ...b, currentHeartbeat: 0 });
      })
    );

    const verboseSubscription = attachDebugPrintfs(withoutHeartbeat, verbosityLevel());

    const lv = await lastValueFrom(states.pipe(takeWhile((x) => x.value.kind !== "terminate")));

    expect(
      lv.transactionStep,
      `did not get CoinsSent, last step was ${JSON.stringify(lv, undefined, 2)}`
    ).to.eq(TransactionStep.CoinsSent);

    verboseSubscription.unsubscribe();
    replayedEvents.unsubscribe();
  }

  it("should allow ceth to eth tx", async () => {
    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    );
    const factories = container.resolve(SifchainContractFactories);
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories);
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0];
    const sendAmount = BigNumber.from(3500);
    const networkDescriptor = devEnvObject!.ethResults!.chainId;

    let testSifAccount: EbRelayerAccount = createTestSifAccount();
    fundSifAccount(
      devEnvObject!.sifResults!.adminAddress!.account,
      testSifAccount!.account,
      10000000000,
      "rowan",
      devEnvObject!.sifResults!.adminAddress!.homeDir
    );

    // Need to have a burn of eth happen at least once or there's no data about eth in the token metadata
    await executeLock(
      contracts,
      sendAmount,
      ethereumAccounts.availableAccounts[0],
      web3.utils.utf8ToHex(testSifAccount.account),
      false,
      "ceth to eth"
    );

    const evmRelayerEvents = sifwatch(
      {
        evmrelayer: "/tmp/sifnode/evmrelayer.log",
        sifnoded: "/tmp/sifnode/sifnoded.log",
      },
      hardhat,
      contracts.bridgeBank
    ).pipe(filter((x) => x.kind !== "SifnodedInfoEvent"));

    const states: Observable<State> = evmRelayerEvents.pipe(
      scan(
        (acc: State, v: SifEvent) => {
          if (isTerminalState(acc))
            // we've reached a decision
            return { ...acc, value: { kind: "terminate" } as Terminate };
          switch (v.kind) {
            case "EbRelayerError":
            case "SifnodedError":
              // if we get an actual error, that's always a failure
              return { ...acc, value: { kind: "failure", value: v, message: "simple error" } };
            case "SifHeartbeat":
              // we just store the heartbeat
              return { ...acc, currentHeartbeat: v.value } as State;

            // Ebrelayer side log assertions
            case "EbRelayerEvmStateTransition": {
              let ebrelayerEvent: any = v.data;
              switch (ebrelayerEvent.kind) {
                case "CosmosEvent": {
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.Burn,
                    TransactionStep.CosmosEvent
                  );
                }

                case "PublishedProphecy": {
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.CosmosEvent,
                    TransactionStep.PublishedProphecy
                  );
                }
              }
            }
            // Sifnoded side log assertions
            case "SifnodedPeggyEvent": {
              let sifnodedEvent: any = v.data;
              switch (sifnodedEvent.kind) {
                case "Burn":
                  let cosmos_sender = sifnodedEvent.msg.Interface.cosmos_sender;
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.Initial,
                    TransactionStep.Burn
                  ); // v.data

                case "GetTokenMetadata":
                  return ensureCorrectTransition(
                    acc,
                    v,
                    TransactionStep.Burn,
                    TransactionStep.GetTokenMetadata
                  );

                case "SignProphecy":
                  return {} as State;

                case "ProphecyStatus":
                  // Assert it is successful. But we dn need to coz thats the only end state
                  return {} as State;
              }
              return { ...acc, value: v, createdAt: acc.currentHeartbeat };
            }

            default:
              // we have a new value (of any kind) and it should use the current heartbeat as its creation time
              return { ...acc, value: v, createdAt: acc.currentHeartbeat };
          }
        },
        {
          value: { kind: "initialState" },
          createdAt: 0,
          currentHeartbeat: 0,
          transactionStep: TransactionStep.Initial,
        } as State
      )
    );

    // it's useful to skip debug prints of states where only the heartbeat changed
    const withoutHeartbeat = states.pipe(
      distinctUntilChanged<State>((a, b) => {
        return deepEqual({ ...a, currentHeartbeat: 0 }, { ...b, currentHeartbeat: 0 });
      })
    );

    const verboseSubscription = attachDebugPrintfs(withoutHeartbeat, verbosityLevel());

    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee;

    await executeSifBurn(
      testSifAccount,
      destinationEthereumAddress,
      sendAmount.sub(crossChainCethFee),
      ethDenomHash,
      String(crossChainCethFee),
      networkDescriptor
    );

    const lv = await lastValueFrom(states.pipe(takeWhile((x) => x.value.kind !== "terminate")));
    expect(
      lv.transactionStep,
      `did not get a LogBridgeTokenMint, last step was ${JSON.stringify(lv, undefined, 2)}`
    ).to.eq(TransactionStep.LogBridgeTokenMint);

    verboseSubscription.unsubscribe();
  });

  it.only("should send two locks of ethereum", async () => {
    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    );
    const factories = container.resolve(SifchainContractFactories);
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories);
    const sender1 = ethereumAccounts.availableAccounts[0];
    const smallAmount = BigNumber.from(1017);

    // Do two locks of ethereum
    await executeLock(contracts, smallAmount, sender1, recipient, true, "lock of eth");
    await executeLock(contracts, smallAmount, sender1, recipient, true, "second lock of eth");
  });
});
