import * as hardhat from "hardhat"
import { container } from "tsyringe"
import { DeployedBridgeBank, requiredEnvVar } from "../src/contractSupport"
import { DeploymentName, HardhatRuntimeEnvironmentToken } from "../src/tsyringe/injectionTokens"
import {
  impersonateAccount,
  impersonateBridgeBankAccounts,
  setNewEthBalance,
  setupRopstenDeployment,
  setupSifchainMainnetDeployment,
} from "../src/hardhatFunctions"
import { SifchainContractFactories } from "../src/tsyringe/contracts"
import { buildIbcTokens, readTokenData } from "../src/ibcMatchingTokens"
import { IbcToken } from "../build"
import web3 from "web3"

const MINTER_ROLE: string = web3.utils.soliditySha3("MINTER_ROLE") ?? "0xBADBAD" // this should never fail
if (MINTER_ROLE == "0xBADBAD") throw Error("failed to get MINTER_ROLE")
const DEFAULT_ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000" // to bridgebank

async function main() {
  const [atomOwner] = await hardhat.ethers.getSigners()

  container.register(HardhatRuntimeEnvironmentToken, { useValue: hardhat })

  const deploymentName = requiredEnvVar("DEPLOYMENT_NAME")

  container.register(DeploymentName, { useValue: deploymentName })

  switch (hardhat.network.name) {
    case "ropsten":
      await setupRopstenDeployment(container, hardhat, deploymentName)
      break
    case "mainnet":
    case "hardhat":
    case "localhost":
      await setupSifchainMainnetDeployment(container, hardhat, deploymentName)
      break
  }

  const newToken = (await hardhat.ethers.getContractAt(
    "IbcToken",
    "0xAFd70A528cd5C172de51993C0C4734b205e40062"
  )) as IbcToken
  const bridgeBank = await container.resolve(DeployedBridgeBank).contract

  await newToken.grantRole(DEFAULT_ADMIN_ROLE, bridgeBank.address)
  console.log(
    JSON.stringify({ roleGrantedToBridgeBank: DEFAULT_ADMIN_ROLE, bridgeBank: bridgeBank.address })
  )
  await newToken.grantRole(MINTER_ROLE, bridgeBank.address)
  console.log(JSON.stringify({ roleGrantedToBridgeBank: MINTER_ROLE }))
  await newToken.renounceRole(MINTER_ROLE, await atomOwner.getAddress())
  console.log(JSON.stringify({ roleRenouncedByDeployer: MINTER_ROLE }))
  await newToken.renounceRole(DEFAULT_ADMIN_ROLE, await atomOwner.getAddress())
  console.log(JSON.stringify({ roleRenouncedByDeployer: DEFAULT_ADMIN_ROLE }))
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error)
    process.exit(1)
  })
