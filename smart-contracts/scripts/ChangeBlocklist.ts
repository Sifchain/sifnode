import { ethers } from "hardhat";
import {BridgeBank__factory} from "../build";

const BridgeBankAddress = "0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8";
const BlockListAddress = "0xa74d631Ac62F028b839a60251E9e3Cf905736826";

async function finishUpdate() {
    const [operator] = await ethers.getSigners();
    const factory = await ethers.getContractFactory("BridgeBank") as BridgeBank__factory;
    const bridgebank = await factory.attach(BridgeBankAddress);

    // Set the blocklist variable
    console.log("Setting the blocklist variable");
    await bridgebank.connect(operator).setBlocklist(BlockListAddress);
}

finishUpdate()
    .then(() => {console.log("Bridgebank Blocklist Updated")})
    .catch((err)=> console.error("Encountered an error:", err))