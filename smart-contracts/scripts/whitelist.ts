import * as hardhat from "hardhat";
import {container} from "tsyringe";
import {DeployedBridgeBank} from "../src/contractSupport";
import {HardhatRuntimeEnvironmentToken} from "../src/tsyringe/injectionTokens";
import {setupDeployment} from "../src/hardhatFunctions";
import {SifchainContractFactories} from "../src/tsyringe/contracts";
import {getWhitelistItems} from "../src/whitelist";

async function main() {
    container.register(HardhatRuntimeEnvironmentToken, {useValue: hardhat})
    await setupDeployment(container)

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
