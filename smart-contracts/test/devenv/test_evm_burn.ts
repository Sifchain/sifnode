import * as chai from "chai"
import {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe"
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import {BigNumber} from "ethers"
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities"
import {SifchainContractFactories, MINTER_ROLE, ADMIN_ROLE} from "../../src/tsyringe/contracts"
import {buildDevEnvContracts, DevEnvContracts} from "../../src/contractSupport"
import web3 from "web3"
import {EbRelayerAccount} from "../../src/devenv/sifnoded"
import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import {ethers} from "hardhat"
import {SifnodedAdapter, DEFAULT_PREPAY_AMOUNT} from "./sifnodedAdapter"
import {SifchainAccountsPromise} from "../../src/tsyringe/sifchainAccounts"
import {executeBurn, checkEvmBurnState} from "./evm_burn"
import {getDenomHash, ethDenomHash, Direction} from "./context"

chai.use(solidity)

describe("burn rowan tests", () => {
  dotenv.config()
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")
  const networkDescriptor = devEnvObject?.ethResults?.chainId ?? 31337

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
  })

  it.only("should allow rowan to sifchain", async () => {
    // TODO: Could these be moved out of the test fx? and instantiated via beforeEach?
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )

    const sendAmount = BigNumber.from("50") // 3500 gwei

    // burn from evm don't need prepaid rowan in sifnode side
    let testSifAccount: EbRelayerAccount = await sifnodedAdapter.createTestSifAccount(false)

    // grant the miner
    const sifchainAccountsPromise = container.resolve(SifchainAccountsPromise)
    const ownerAccount = (await sifchainAccountsPromise.accounts).ownerAccount
    await contracts.rowanContract.grantRole(String(MINTER_ROLE), ownerAccount.address)
    await contracts.rowanContract.grantRole(String(ADMIN_ROLE), ownerAccount.address)


    const senderEthereumAccount = ethereumAccounts.availableAccounts[0]
    // mint token to sender
    await contracts.rowanContract.connect(ownerAccount).mint(senderEthereumAccount.address, sendAmount)
    await contracts.rowanContract.connect(ownerAccount).setDenom("rowan")

    // record the init balance before lock
    const initialRowanSenderBalance = await contracts.rowanContract.balanceOf(senderEthereumAccount.address)
    const initialContractBalance = await contracts.rowanContract.balanceOf(contracts.bridgeBank.address)
    const initialRowanReceiverBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      "rowan"
    )

    let originalVerboseLevel: string | undefined = process.env["VERBOSE"]
    process.env["VERBOSE"] = "summary"

    // burn rowan
    const tx = await executeBurn(
      contracts,
      sendAmount,
      senderEthereumAccount,
      web3.utils.utf8ToHex(testSifAccount.account),
      contracts.rowanContract,
    )

    await checkEvmBurnState(contracts, tx, sendAmount, "rowan")

    // These are temporarily added to make the logging lvl lower
    process.env["VERBOSE"] = originalVerboseLevel

    console.log("Lock complete")

    // get the balance after lock
    const finalRowanSenderBalance = await contracts.rowanContract.balanceOf(senderEthereumAccount.address)
    const finalContractBalance = await contracts.rowanContract.balanceOf(contracts.bridgeBank.address)
    const finalRowanReceiverBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      "rowan"
    )

    console.log("Before lock the sender's balance is ", initialRowanSenderBalance)
    console.log("Before lock the contract's balance is ", initialContractBalance)
    console.log("Before lock the receiver's balance is ", initialRowanReceiverBalance)

    console.log("After lock the sender's balance is ", finalRowanSenderBalance)
    console.log("After lock the contract's balance is ", finalContractBalance)
    console.log("After lock the receiver's balance is ", finalRowanReceiverBalance)

    expect(initialRowanSenderBalance.sub(sendAmount), "should be equal ").eq(
      finalRowanSenderBalance
    )
    expect(initialRowanReceiverBalance.add(sendAmount), "should be equal ").eq(
      finalRowanReceiverBalance,
    )
  })
})
