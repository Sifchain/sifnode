import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {
    DeployedBridgeBank,
    DeployedBridgeRegistry,
    DeployedCosmosBridge,
    requiredEnvVar
} from "../src/contractSupport";
import {DeploymentName, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {setupRopstenDeployment, setupSifchainMainnetDeployment} from "../src/hardhatFunctions";
import {
    BridgeBankProxy, BridgeRegistryProxy,
    BridgeTokenSetup,
    RowanContract,
    SifchainContractFactories
} from "../src/tsyringe/contracts";
import * as dotenv from "dotenv";


export type DeployedContractAddresses = {
    bridgeBank: string,
    bridgeRegistry: string,
    rowanContract: string,
}
// Usage
//
// npx hardhat run scripts/deploy_contracts.ts

async function main() {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    await container.resolve(BridgeTokenSetup).complete
    const bridgeBank = await container.resolve(BridgeBankProxy).contract
    const bridgeRegistry = await container.resolve(BridgeRegistryProxy).contract
    const rowanContract = await container.resolve(RowanContract).contract
    const result: DeployedContractAddresses = {
        bridgeBank: bridgeBank.address,
        bridgeRegistry: bridgeRegistry.address,
        rowanContract: rowanContract.address
    }
    console.log(JSON.stringify(result))
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
