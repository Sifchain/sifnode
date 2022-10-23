import * as chai from "chai"
import {expect} from "chai"
import {solidity} from "ethereum-waffle"
import {container} from "tsyringe"
import {HardhatRuntimeEnvironmentToken} from "../../src/tsyringe/injectionTokens"
import * as hardhat from "hardhat"
import {BigNumber} from "ethers"
import {ethereumResultsToSifchainAccounts, readDevEnvObj} from "../../src/tsyringe/devenvUtilities"
import {SifchainContractFactories, MINTER_ROLE} from "../../src/tsyringe/contracts"
import {buildDevEnvContracts} from "../../src/contractSupport"
import web3 from "web3"
import {EbRelayerAccount, crossChainFeeBase, crossChainBurnFee} from "../../src/devenv/sifnoded"

import * as dotenv from "dotenv"
import "@nomiclabs/hardhat-ethers"
import {SifnodedAdapter} from "./sifnodedAdapter"
import {getDenomHash, ethDenomHash} from "./context"
import {checkSifnodeBurnState} from "./sifnode_lock_burn"
import {executeLock, checkEvmLockState} from "./evm_lock_burn"
import {SifchainAccountsPromise} from "../../src/tsyringe/sifchainAccounts"

chai.use(solidity)

describe("lock and burn tests", () => {
  dotenv.config()
  // This test only works when devenv is running, and that requires a connection to localhost
  expect(hardhat.network.name, "please use devenv").to.eq("localhost")

  const devEnvObject = readDevEnvObj("environment.json")
  const networkDescriptor = devEnvObject?.ethResults?.chainId ?? 9999

  const sifnodedAdapter: SifnodedAdapter = new SifnodedAdapter(
    devEnvObject!.sifResults!.adminAddress!.homeDir,
    devEnvObject!.sifResults!.adminAddress!.account,
    process.env["GOBIN"]
  )

  before("register HardhatRuntimeEnvironmentToken", async () => {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
  })

  it.only("should allow erc20 back to Ethereum", async () => {
    // TODO: Could these be moved out of the test fx? and instantiated via beforeEach?
    const factories = container.resolve(SifchainContractFactories)
    const contracts = await buildDevEnvContracts(devEnvObject, hardhat, factories)

    // deploy a new erc20 token
    const bridgeToken = await factories.bridgeToken
    const erc20 = await bridgeToken.deploy("erc20", "erc20", 18, "erc20denom")

    const ethereumAccounts = await ethereumResultsToSifchainAccounts(
      devEnvObject.ethResults!,
      hardhat.ethers.provider
    )
    const destinationEthereumAddress = ethereumAccounts.availableAccounts[0]

    // const sendAmount = BigNumber.from(5 * ETH) // 3500 gwei
    const sendAmount = BigNumber.from("5000000000000000000") // 3500 gwei

    let testSifAccount: EbRelayerAccount = sifnodedAdapter.createTestSifAccount()
    let originalVerboseLevel: string | undefined = process.env["VERBOSE"]
    process.env["VERBOSE"] = "summary"
    // Need to have a burn of eth happen at least once or there's no data about eth in the token metadata
    let tx = await executeLock(
      contracts,
      sendAmount,
      ethereumAccounts.availableAccounts[1],
      web3.utils.utf8ToHex(testSifAccount.account)
    )

    await checkEvmLockState(contracts, tx, sendAmount, ethDenomHash)

    console.log("lock eth done")

    // grant the miner
    const erc20Denom = getDenomHash(networkDescriptor, erc20.address.toString())

    const sifchainAccountsPromise = container.resolve(SifchainAccountsPromise)
    const ownerAccount = (await sifchainAccountsPromise.accounts).ownerAccount
    await erc20.grantRole(String(MINTER_ROLE), ownerAccount.address)

    // mint token to sender
    await erc20
      .connect(ownerAccount)
      .mint(ethereumAccounts.availableAccounts[1].address, sendAmount)

    // lock the erc20 token
    tx = await executeLock(
      contracts,
      sendAmount,
      ethereumAccounts.availableAccounts[1],
      web3.utils.utf8ToHex(testSifAccount.account),
      erc20,
    )

    await checkEvmLockState(contracts, tx, sendAmount, erc20Denom)
    console.log("lock erc20 done")

    // record the init balance before lock
    const initialErc20SenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      erc20Denom
    )
    const initialErc20ReceiverBalance = await erc20.balanceOf(destinationEthereumAddress.address)

    // These are temporarily added to make the logging lvl lower
    process.env["VERBOSE"] = originalVerboseLevel

    console.log("Lock complete")

    let crossChainCethFee = crossChainFeeBase * crossChainBurnFee
    let burnAmount = BigNumber.from("2300000000000000000") // 2300 gwei

    await checkSifnodeBurnState(
      sifnodedAdapter,
      contracts,
      testSifAccount,
      destinationEthereumAddress,
      burnAmount,
      erc20Denom,
      String(crossChainCethFee),
      networkDescriptor
    )

    // Here we verify the user balance is correct
    // get the balance after burn
    const finalErc20SenderBalance = await sifnodedAdapter.getBalance(
      testSifAccount.account,
      erc20Denom
    )
    const finalErc20ReceiverBalance = await erc20.balanceOf(destinationEthereumAddress.address)

    console.log("Before burn the sender's balance is ", initialErc20SenderBalance)
    console.log("Before burn the receiver's balance is ", initialErc20ReceiverBalance)

    console.log("After burn the sender's balance is ", finalErc20SenderBalance)
    console.log("After burn the receiver's balance is ", finalErc20ReceiverBalance)

    expect(initialErc20SenderBalance.sub(burnAmount), "should be equal ").eq(
      finalErc20SenderBalance
    )
    expect(initialErc20ReceiverBalance.add(burnAmount), "should be equal ").eq(
      finalErc20ReceiverBalance
    )
  })
})
