import * as fs from 'fs';
import {BaseContract, ContractFactory} from "ethers";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {inject, injectable} from "tsyringe";
import {
    DeploymentChainId,
    DeploymentDirectory,
    DeploymentName,
    HardhatRuntimeEnvironmentToken
} from "./tsyringe/injectionTokens";
import {BridgeBank, BridgeRegistry, BridgeToken, CosmosBridge} from "../build";
import * as path from "path"

export async function getContractFromTruffleArtifact<T extends BaseContract>(
    hre: HardhatRuntimeEnvironment,
    filename: string,
    chainId: number
): Promise<T> {
    const artifactContents = fs.readFileSync(filename, {encoding: "utf-8"})
    const parsedArtifactContents = JSON.parse(artifactContents)
    const truffle = require("@truffle/contract")
    const truffleContract = (truffle as any)(parsedArtifactContents)
    const contractData = truffleContract.networks[chainId]
    const ethersContract = await hre.ethers.getContractAt(truffleContract.abi, contractData.address)
    return ethersContract as T
}

@injectable()
export class DeployableContract<T extends BaseContract> {
    readonly contract: Promise<T>
    contractName(): string {
        return "must override this"
    }

    constructor(
        @inject(HardhatRuntimeEnvironmentToken) hre: HardhatRuntimeEnvironment,
        @inject(DeploymentDirectory) deploymentDirectory: string,
        @inject(DeploymentName) deploymentName: string,
        @inject(DeploymentChainId) deploymentChainId: number,
    ) {
        const t = this
        const r = typeof this
        const n = this.contractName()
        this.contract = getContractFromTruffleArtifact(
            hre,
            path.join(deploymentDirectory, deploymentName, `${n}.json`),
            deploymentChainId
        )
    }
}

@injectable()
export class DeployedBridgeBank extends DeployableContract<BridgeBank> {
    contractName() {
        return "BridgeBank"
    }
}

@injectable()
export class DeployedBridgeToken extends DeployableContract<BridgeToken> {
    contractName() {
        return "BridgeToken"
    }
}

@injectable()
export class DeployedBridgeRegistry extends DeployableContract<BridgeRegistry> {
    contractName() {
        return "BridgeRegistry"
    }
}

@injectable()
export class DeployedCosmosBridge extends DeployableContract<CosmosBridge> {
    contractName() {
        return "CosmosBridge"
    }
}
