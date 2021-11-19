import * as hardhat from "hardhat"
import { container } from "tsyringe"
import { DeployedBridgeBank, requiredEnvVar } from "../src/contractSupport"
import { DeploymentName, HardhatRuntimeEnvironmentToken } from "../src/tsyringe/injectionTokens"
import {
  impersonateBridgeBankAccounts,
  setupRopstenDeployment,
  setupSifchainMainnetDeployment,
} from "../src/hardhatFunctions"
import * as dotenv from "dotenv"
import { processTokenData } from "../src/ibcMatchingTokens"

const envconfig = dotenv.config()

async function main() {
  const [bridgeBankOwner] = await hardhat.ethers.getSigners()

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

  const useForking = process.env["USE_FORKING"]
  if (useForking) await impersonateBridgeBankAccounts(container, hardhat, deploymentName)

  const bridgeBank = await container.resolve(DeployedBridgeBank).contract
  const owner = await bridgeBank.owner()
  const ownerAsSigner = await hardhat.ethers.getSigner(owner)

  await processTokenData(
    bridgeBank.connect(ownerAsSigner),
    requiredEnvVar("TOKEN_ADDRESS_FILE"),
    container
  )
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error)
    process.exit(1)
  })
