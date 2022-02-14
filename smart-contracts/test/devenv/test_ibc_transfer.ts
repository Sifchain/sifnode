import * as chai from "chai"
import {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe"
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import {BigNumber} from "ethers"
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities"
import {SifchainContractFactories, MINTER_ROLE} from "../../src/tsyringe/contracts"
import {buildDevEnvContracts, DevEnvContracts} from "../../src/contractSupport"
import web3 from "web3"
import {EbRelayerAccount} from "../../src/devenv/sifnoded"
import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import {ethers} from "hardhat"
import {SifnodedAdapter} from "./sifnodedAdapter"
import {SifchainAccountsPromise} from "../../src/tsyringe/sifchainAccounts"
import {executeLock, checkEvmLockState} from "./evm_lock_burn"
import {getDenomHash, ethDenomHash} from "./context"
import {checkSifnodeIbcTransferState} from "./sifnode_ibc_transfer"

chai.use(solidity)

const ibcDenom = "ibc-denom"
const srcPort = "port-0"
const srcChannel = "channel-0"

describe("ibc transfer tests", () => {
  dotenv.config()
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
  })

  it.only("should allow ibc transfer token to sifchain", async () => {    

    let testSifAccount: EbRelayerAccount = sifnodedAdapter.createTestSifAccount()
    const sendAmount = BigNumber.from(1234)

    // record the init balance before ibc transfer
    const initialReceiverBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      ibcDenom
    )

    let originalVerboseLevel: string | undefined = process.env["VERBOSE"]
    process.env["VERBOSE"] = "summary"

    await checkSifnodeIbcTransferState(sifnodedAdapter, srcPort, srcChannel, testSifAccount.account, sendAmount, ethDenomHash)

    // These are temporarily added to make the logging lvl lower
    process.env["VERBOSE"] = originalVerboseLevel

    console.log("ibc transfer complete")

    // get the balance after ibc transfer
    const finalReceiverBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      ibcDenom
    )

    console.log("Before lock the sender's balance is ", initialReceiverBalance)
    console.log("After lock the sender's balance is ", finalReceiverBalance)

    expect(initialReceiverBalance.add(sendAmount), "should be equal ").eq(
      finalReceiverBalance
    )
  })

})
