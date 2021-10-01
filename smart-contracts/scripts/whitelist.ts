import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {DeployedBridgeBank, requiredEnvVar} from "../src/contractSupport";
import {DeploymentName, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {setupRopstenDeployment, setupSifchainMainnetDeployment} from "../src/hardhatFunctions";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import {getWhitelistItems} from "../src/whitelist";

async function main() {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    const deploymentName = requiredEnvVar("DEPLOYMENT_NAME")
    container.register(DeploymentName, {useValue: deploymentName})

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

    const bridgeBank = await container.resolve(DeployedBridgeBank).contract
    const btf = await container.resolve(SifchainContractFactories).bridgeToken
    const result = await getWhitelistItems(bridgeBank, btf)
    console.log(JSON.stringify(result))
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
