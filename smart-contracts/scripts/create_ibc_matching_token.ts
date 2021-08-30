import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {DeployedBridgeBank, requiredEnvVar} from "../src/contractSupport";
import {DeploymentName, HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {
    impersonateAccount, setNewEthBalance,
    setupRopstenDeployment,
    setupSifchainMainnetDeployment
} from "../src/hardhatFunctions";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import {buildIbcTokens, readTokenData} from "../src/ibcMatchingTokens";

async function main() {
    const [bridgeBankOwner] = await hardhat.ethers.getSigners();

    hardhat.network.config
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

    const factories = await container.resolve(SifchainContractFactories) as SifchainContractFactories

    const ibcTokenFactory = (await factories.ibcToken).connect(bridgeBankOwner)

    const bridgeBank = await container.resolve(DeployedBridgeBank).contract
    const bridgeBankOperater = await bridgeBank.operator()

    if (hardhat.network.name == "localhost") {
        await impersonateAccount(hardhat, await bridgeBank.operator(), undefined, async impersonatedAccount => {
            await buildIbcTokens(ibcTokenFactory, await readTokenData(requiredEnvVar("TOKEN_FILE")), bridgeBank.connect(impersonatedAccount))
        })
    } else {
        await buildIbcTokens(ibcTokenFactory, await readTokenData(requiredEnvVar("TOKEN_FILE")), bridgeBank)
    }
    console.log("created ibc tokens on this network: ", hardhat.network.name)
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
