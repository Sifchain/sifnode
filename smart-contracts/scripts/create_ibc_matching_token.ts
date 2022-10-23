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

  const useForking = !!process.env["USE_FORKING"]
  const mintTokens = !!process.env["MINT_TOKENS"]
  if (useForking) await impersonateBridgeBankAccounts(container, hardhat, deploymentName)

  const factories = (await container.resolve(
    SifchainContractFactories
  )) as SifchainContractFactories

  const ibcTokenFactory = (await factories.ibcToken).connect(bridgeBankOwner)

  const bridgeBank = await container.resolve(DeployedBridgeBank).contract

  const startingBalance = await hardhat.ethers.provider.getBalance(bridgeBankOwner.address)
  await buildIbcTokens(
    ibcTokenFactory,
    await readTokenData(requiredEnvVar("TOKEN_FILE")),
    bridgeBank,
    mintTokens
  )
  const endingBalance = await hardhat.ethers.provider.getBalance(bridgeBankOwner.address)
  console.log(
    JSON.stringify({
      createdTokensOnNetwork: hardhat.network.name,
      cost: hardhat.ethers.utils.formatEther(startingBalance.sub(endingBalance)),
    })
  )
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error)
    process.exit(1)
  })
