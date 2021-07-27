import * as fs from 'fs';
import {BaseContract, BigNumber, BigNumberish, ContractFactory} from "ethers";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DependencyContainer, inject, injectable} from "tsyringe";
import {
    BridgeBankMainnetUpgradeAdmin,
    DeploymentChainId,
    DeploymentDirectory,
    DeploymentName,
    HardhatRuntimeEnvironmentToken
} from "./tsyringe/injectionTokens";
import {
    BridgeBank,
    BridgeRegistry,
    BridgeToken,
    BridgeToken__factory,
    CosmosBridge,
    ERC20
} from "../build";
import * as path from "path"
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {EthereumAddress, NotNativeCurrencyAddress} from "./ethereumAddress";
import * as hardhat from "hardhat";
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
    const btf = await hre.ethers.getContractFactory("BridgeToken") as BridgeToken__factory
    const existingRowanToken: BridgeToken = await btf.attach(eRowanMainnetAddress)
    const syntheticBridgeToken = {
        contract: Promise.resolve(existingRowanToken),
        contractName: () => "BridgeToken"
    }
    c.register(BridgeBankMainnetUpgradeAdmin, {useValue: "0x7c6c6ea036e56efad829af5070c8fb59dc163d88"})
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
