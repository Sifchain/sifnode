import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {DeployedBridgeBank, DeployedCosmosBridge, requiredEnvVar} from "../src/contractSupport";
import {DeploymentName, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {setupRopstenDeployment, setupSifchainMainnetDeployment} from "../src/hardhatFunctions";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import * as dotenv from "dotenv";

async function main() {
    const [proxyAdmin] = await hardhat.ethers.getSigners();

    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})

    container.register(DeploymentName, {useValue: requiredEnvVar("DEPLOYMENT_NAME")})

    switch (hardhat.network.name) {
        case "ropsten":
            await setupRopstenDeployment(container, hardhat, "sifchain-testnet-042-ibc")
            break
        case "mainnet":
            await setupSifchainMainnetDeployment(container, hardhat)
            break
    }

    // upgradeProxy wants two things: a ContractFactory to build the new logic contract,
    // and an existing contract that will be replaced.

    const factories = await container.resolve(SifchainContractFactories) as SifchainContractFactories
    const bridgeBankFactory = await factories.bridgeBank
    const cosmosBridgeFactory = await factories.cosmosBridge

    const existingBridgeBank = await container.resolve(DeployedBridgeBank).contract
    const existingCosmosBridge = await container.resolve(DeployedCosmosBridge).contract

    await hardhat.upgrades.upgradeProxy(existingBridgeBank, bridgeBankFactory.connect(proxyAdmin))
    await hardhat.upgrades.upgradeProxy(existingCosmosBridge, cosmosBridgeFactory.connect(proxyAdmin))
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
