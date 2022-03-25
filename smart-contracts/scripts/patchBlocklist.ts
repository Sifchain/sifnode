import { ethers } from "hardhat";
import {BridgeBank__factory} from "../build";

const BridgeBankAddress = "0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8";
const BlockListAddress = "0x9C8a2011cCb697D7EDe3c94f9FBa5686a04DeACB";

async function finishUpdate() {
    const [admin, operator, pauser] = await ethers.getSigners();
    const factory = await ethers.getContractFactory("BridgeBank") as BridgeBank__factory;
    const bridgebank = await factory.attach(BridgeBankAddress);

    // Set the blocklist variable
    console.log("Setting the blocklist variable");
    await bridgebank.connect(operator).setBlocklist(BlockListAddress);

    // Unpause the bridgebank
    console.log("Unpausing the variable");
    await bridgebank.connect(pauser).unpause();
}

finishUpdate()
    .then(() => {console.log("Tests completed")})
    .catch((err)=> console.error("Encountered an error:", err))