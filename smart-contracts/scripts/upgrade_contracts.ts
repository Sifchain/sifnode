import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {DeployedBridgeBank, DeployedCosmosBridge, requiredEnvVar} from "../src/contractSupport";
import {DeploymentName, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {setupRopstenDeployment, setupSifchainMainnetDeployment} from "../src/hardhatFunctions";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import * as dotenv from "dotenv";

// Usage
//
// npx hardhat run scripts/upgrade_contracts.ts
//
// Uses these environment variables:
// (using the dotenv plugin, so you can have them in a .env file also)
//
// MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/...
// ROPSTEN_URL=https://eth-ropsten.alchemyapi.io/v2/...
// ROPSTEN_PROXY_ADMIN_PRIVATE_KEY=aaaa...
// DEPLOYMENT_NAME="sifchain-testnet-042-ibc"

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

    await hardhat.upgrades.upgradeProxy(existingBridgeBank, bridgeBankFactory.connect(proxyAdmin), {unsafeAllowCustomTypes: true})
    await hardhat.upgrades.upgradeProxy(existingCosmosBridge, cosmosBridgeFactory.connect(proxyAdmin), {unsafeAllowCustomTypes: true})
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
