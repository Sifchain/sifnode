import {BigNumber, BigNumberish} from "ethers";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DependencyContainer} from "tsyringe";
import {
    BridgeBankMainnetUpgradeAdmin,
    CosmosBridgeMainnetUpgradeAdmin,
    DeploymentChainId,
    DeploymentDirectory,
    DeploymentName
} from "./tsyringe/injectionTokens";
import {BridgeToken, BridgeToken__factory, ERC20} from "../build";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {NotNativeCurrencyAddress} from "./ethereumAddress";
import {DeployedBridgeToken} from "./contractSupport";

export const eRowanMainnetAddress = "0x07bac35846e5ed502aa91adf6a9e7aa210f2dcbe"

export async function impersonateAccount<T>(
    hre: HardhatRuntimeEnvironment,
    address: string,
    newBalance: BigNumberish | undefined,
    fn: (s: SignerWithAddress) => Promise<T>
) {
    await hre.network.provider.request({
        method: "hardhat_impersonateAccount",
        params: [address]
    });
    if (newBalance) {
        await setNewEthBalance(hre, address, newBalance)
    }
    const signer = await hre.ethers.getSigner(address)
    const result = await fn(signer)
    await hre.network.provider.request({
        method: "hardhat_stopImpersonatingAccount",
        params: [address],
    });
    return result
}

export async function startImpersonateAccount<T>(
    hre: HardhatRuntimeEnvironment,
    address: string,
    newBalance?: BigNumberish
): Promise<SignerWithAddress> {
    await hre.network.provider.request({
        method: "hardhat_impersonateAccount",
        params: [address]
    });
    if (newBalance) {
        await setNewEthBalance(hre, address, newBalance)
    }
    return await hre.ethers.getSigner(address)
}

export async function stopImpersonateAccount<T>(
    hre: HardhatRuntimeEnvironment,
    address: string,
) {
    await hre.network.provider.request({
        method: "hardhat_stopImpersonatingAccount",
        params: [address],
    });
}

export async function setNewEthBalance(
    hre: HardhatRuntimeEnvironment,
    address: string,
    newBalance: BigNumberish | undefined,
) {
    const newValue = BigNumber.from(newBalance)
    await hre.network.provider.send("hardhat_setBalance", [
        address,
        newValue.toHexString().replace(/0x0+/, "0x")
    ]);
}

export async function setupSifchainMainnetDeployment(c: DependencyContainer, hre: HardhatRuntimeEnvironment) {
    c.register(DeploymentDirectory, {useValue: "./deployments"})
    c.register(DeploymentName, {useValue: "sifchain"})
    c.register(DeploymentChainId, {useValue: 1})
    const bridgeTokenFactory = await hre.ethers.getContractFactory("BridgeToken") as BridgeToken__factory

    // BrideToken for Rowan doesn't have a json file in deployments, so we need to build DeployedBridgeToken by hand
    // instead of using
    const existingRowanToken: BridgeToken = await bridgeTokenFactory.attach(eRowanMainnetAddress)
    const syntheticBridgeToken = {
        contract: Promise.resolve(existingRowanToken),
        contractName: () => "BridgeToken"
    }
    c.register(BridgeBankMainnetUpgradeAdmin, {useValue: "0x7c6c6ea036e56efad829af5070c8fb59dc163d88"})
    //TODO this is probably the wrong address?
    c.register(CosmosBridgeMainnetUpgradeAdmin, {useValue: "0x7c6c6ea036e56efad829af5070c8fb59dc163d88"})
    c.register(DeployedBridgeToken, {useValue: syntheticBridgeToken as DeployedBridgeToken})
}

export async function approveThenDo<T>(
    token: ERC20,
    agent: NotNativeCurrencyAddress,
    amount: BigNumberish,
    account: SignerWithAddress,
    fn: () => Promise<T>
) {
    await token.connect(account).approve(agent.address, amount)
    return fn()
}
