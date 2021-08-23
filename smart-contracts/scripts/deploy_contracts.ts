import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {DeployedBridgeBank, DeployedCosmosBridge, requiredEnvVar} from "../src/contractSupport";
import {DeploymentName, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {setupRopstenDeployment, setupSifchainMainnetDeployment} from "../src/hardhatFunctions";
import {BridgeTokenSetup, SifchainContractFactories} from "../src/tsyringe/contracts";
import * as dotenv from "dotenv";

// Usage
//
// npx hardhat run scripts/deploy_contracts.ts

async function main() {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    await container.resolve(BridgeTokenSetup).complete
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
