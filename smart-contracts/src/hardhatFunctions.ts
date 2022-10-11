import { BigNumber, BigNumberish } from "ethers"
import { HardhatRuntimeEnvironment } from "hardhat/types"
import { DependencyContainer } from "tsyringe"
import {
  BridgeBankMainnetUpgradeAdmin,
  CosmosBridgeMainnetUpgradeAdmin,
  DeploymentChainId,
  DeploymentDirectory,
  DeploymentName,
  HardhatRuntimeEnvironmentToken,
} from "./tsyringe/injectionTokens"
import { BridgeToken, BridgeToken__factory } from "../build"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { NotNativeCurrencyAddress } from "./ethereumAddress"
import { DeployedBridgeBank, DeployedBridgeToken } from "./contractSupport"
import { SifchainAccounts, SifchainAccountsPromise } from "./tsyringe/sifchainAccounts"

export const eRowanMainnetAddress = "0x07bac35846e5ed502aa91adf6a9e7aa210f2dcbe"

export async function impersonateAccount<T>(
  hre: HardhatRuntimeEnvironment,
  address: string,
  newBalance: BigNumberish | undefined,
  fn: (s: SignerWithAddress) => Promise<T>
) {
  await hre.network.provider.request({
    method: "hardhat_impersonateAccount",
    params: [address],
  })
  if (newBalance) {
    await setNewEthBalance(hre, address, newBalance)
  }
  const signer = await hre.ethers.getSigner(address)
  const result = await fn(signer)
  await hre.network.provider.request({
    method: "hardhat_stopImpersonatingAccount",
    params: [address],
  })
  return result
}

export async function startImpersonateAccount<T>(
  hre: HardhatRuntimeEnvironment,
  address: string,
  newBalance?: BigNumberish
): Promise<SignerWithAddress> {
  await hre.network.provider.request({
    method: "hardhat_impersonateAccount",
    params: [address],
  })
  if (newBalance) {
    await setNewEthBalance(hre, address, newBalance)
  }
  return await hre.ethers.getSigner(address)
}

export async function stopImpersonateAccount<T>(hre: HardhatRuntimeEnvironment, address: string) {
  await hre.network.provider.request({
    method: "hardhat_stopImpersonatingAccount",
    params: [address],
  })
}

export async function setNewEthBalance(
  hre: HardhatRuntimeEnvironment,
  address: string,
  newBalance: BigNumberish | undefined
) {
  const newValue = BigNumber.from(newBalance)
  await hre.network.provider.send("hardhat_setBalance", [
    address,
    newValue.toHexString().replace(/0x0+/, "0x"),
  ])
}

export async function setupDeployment(c: DependencyContainer) {
  const hre = c.resolve(HardhatRuntimeEnvironmentToken) as HardhatRuntimeEnvironment
  let deploymentName = process.env["DEPLOYMENT_NAME"]
  console.log("Attempted resolving deployment context. Deployment name:", deploymentName)
  switch (deploymentName) {
    case "sifchain":
    case "sifchain-1":
      setupSifchainMainnetDeployment(c, hre, deploymentName)
      break
    case "lance-deployment":
      setupLanceDeployment(c, hre, deploymentName)
      break
    case undefined:
      break
    default:
      setupRopstenDeployment(c, hre, deploymentName)
      break
  }
  console.log("Deployment setup complete")
}

export async function setupSifchainMainnetDeployment(
  c: DependencyContainer,
  hre: HardhatRuntimeEnvironment,
  deploymentName: "sifchain" | "sifchain-1"
) {
  c.register(DeploymentDirectory, { useValue: "./deployments" })
  c.register(DeploymentName, { useValue: deploymentName })
  // We'd like to be able to use chainId from the provider,
  // but it doesn't actually work.  It returns 1 even when
  // you're looking at forked ropsten.
  // const chainId = (await hre.ethers.provider.getNetwork()).chainId
  c.register(DeploymentChainId, { useValue: 1 })
  const bridgeTokenFactory = (await hre.ethers.getContractFactory(
    "BridgeToken"
  )) as BridgeToken__factory

  // BridgeToken for Rowan doesn't have a json file in deployments, so we need to build DeployedBridgeToken by hand
  // instead of using
  const existingRowanToken: BridgeToken = await bridgeTokenFactory.attach(eRowanMainnetAddress)
  const syntheticBridgeToken = {
    contract: Promise.resolve(existingRowanToken),
    contractName: () => "BridgeToken",
  }
  c.register(BridgeBankMainnetUpgradeAdmin, {
    useValue: "0x7c6c6ea036e56efad829af5070c8fb59dc163d88",
  })
  c.register(CosmosBridgeMainnetUpgradeAdmin, {
    useValue: "0x7c6c6ea036e56efad829af5070c8fb59dc163d88",
  })
  c.register(DeployedBridgeToken, { useValue: syntheticBridgeToken as DeployedBridgeToken })
}

async function setupLanceDeployment(
  c: DependencyContainer,
  hre: HardhatRuntimeEnvironment,
  deploymentName: "lance-deployment"
) {
  console.log("Resolving lance deployment setup")
  c.register(SifchainAccountsPromise, {useValue: new SifchainAccountsPromise(getSifchainAccounts(hre)) })
}

async function getSifchainAccounts(hardhat: HardhatRuntimeEnvironment): Promise<SifchainAccounts>{
  const operatorAccount = await hardhat.ethers.getSigner("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
  const ownerAccount = await hardhat.ethers.getSigner("0x70997970c51812dc3a010c7d01b50e0d17dc79c8")
  const pauserAccount = await hardhat.ethers.getSigner("0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc")
  const validatatorAccounts = await hardhat.ethers.getSigner("0x90f79bf6eb2c4f870365e785982e1f101e93b906")
  const extraAccounts = await hardhat.ethers.getSigner("0x15d34aaf54267db7d7c367839aaf71a00a2c6a65")

  console.log("Using hardcoded default evm addresses. FOR DEV ONLY")

  return new SifchainAccounts(operatorAccount, ownerAccount, pauserAccount, [validatatorAccounts], [extraAccounts])
}

export async function impersonateBridgeBankAccounts(
  c: DependencyContainer,
  hre: HardhatRuntimeEnvironment
) {
  const bridgeBank = await c.resolve(DeployedBridgeBank).contract
  const operator = await bridgeBank.operator()
  const owner = await bridgeBank.owner()
  const pauser = owner // TODO you can't look up pausers, so probably need to add a pauser
  startImpersonateAccount(hre, operator)
  startImpersonateAccount(hre, owner)
  startImpersonateAccount(hre, pauser)

  await setNewEthBalance(hre, owner, BigNumber.from("10000000000000000000000"))
  await setNewEthBalance(hre, operator, BigNumber.from("10000000000000000000000"))
  await setNewEthBalance(hre, pauser, BigNumber.from("10000000000000000000000"))
}

export async function setupRopstenDeployment(
  c: DependencyContainer,
  hre: HardhatRuntimeEnvironment,
  deploymentName: string
) {
  c.register(DeploymentDirectory, { useValue: "./deployments" })
  c.register(DeploymentName, { useValue: deploymentName })
  c.register(DeploymentChainId, { useValue: 3 })
}
