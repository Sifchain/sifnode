import * as fs from 'fs';
import {BaseContract, ContractFactory} from "ethers";
import {HardhatRuntimeEnvironment} from "hardhat/types";
import {inject, injectable, singleton} from "tsyringe";
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

@singleton()
export class DeployedBridgeBank extends DeployableContract<BridgeBank> {
    contractName() {
        return "BridgeBank"
    }
}

// Note that this class isn't injectable, since we don't have the right
// json artifacts for BridgeToken
export class DeployedBridgeToken extends DeployableContract<BridgeToken> {
    contractName() {
        return "BridgeToken"
    }
}

@singleton()
export class DeployedBridgeRegistry extends DeployableContract<BridgeRegistry> {
    contractName() {
        return "BridgeRegistry"
    }
}

@singleton()
export class DeployedCosmosBridge extends DeployableContract<CosmosBridge> {
    contractName() {
        return "CosmosBridge"
    }
}

/**
 * Throws an exception if name is not present in process.env
 *
 * @param name
 */
export function requiredEnvVar(name: string): string {
    const result = process.env[name]
    if (typeof result === 'string') {
        return result
    } else {
        throw `No setting for ${name} in environment`
    }
}
