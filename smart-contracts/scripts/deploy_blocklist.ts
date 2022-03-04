/**
 * This script will deploy a new Blocklist contract to the requested EVM chain
 **/

import { ethers } from "hardhat";
import { Blocklist__factory, Blocklist } from "../build"

const BLOCKLIST_ADMIN_PRIVATE_KEY: string = process.env.BLOCKLIST_ADMIN_PRIVATE_KEY || "";

export async function deploy_blocklist() : Promise<Blocklist> {
    const admin = new ethers.Wallet(BLOCKLIST_ADMIN_PRIVATE_KEY);
    const blocklistFactory = await ethers.getContractFactory("Blocklist") as Blocklist__factory;
    const blocklist = await blocklistFactory.deploy()
    console.log("Blocklist deployed to: ", blocklist.address);
    return blocklist
}

deploy_blocklist()
    .then(() => {console.log("Deploy blocklist successfully")})
    .catch((err) => console.error("Error occurred: ", err))