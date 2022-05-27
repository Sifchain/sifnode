/**
 * This script will manually add/remove a single address to the OFAC blocklist and is for manual testing 
 * of the blocklist features only. It should not be used to modify the blocklist in any other way as any
 * changes will be updated automatically by the daily update script.
 */

import { ethers } from "hardhat";
import { Blocklist__factory } from "../build"

const BLOCKLIST_ADDRESS: string = process.env.BLOCKLIST_ADDRESS || "0x9C8a2011cCb697D7EDe3c94f9FBa5686a04DeACB";
const BLOCKLIST_ADMIN_PRIVATE_KEY: string = process.env.BLOCKLIST_ADMIN_PRIVATE_KEY || "";
const BLOCK_ADDRESS: string = process.env.BLOCK_ADDRESS || "";

async function add_address(address: string) {
    const admin = new ethers.Wallet(BLOCKLIST_ADMIN_PRIVATE_KEY);
    const blocklistFactory = await ethers.getContractFactory("Blocklist") as Blocklist__factory;
    const blocklist = await blocklistFactory.attach(BLOCKLIST_ADDRESS);
    await blocklist.connect(admin).removeFromBlocklist(BLOCKLIST_ADDRESS);
}

add_address(BLOCK_ADDRESS)
    .then(() => {console.log("Delete Address Operation Completed")})
    .catch((err) => console.error("Error occurred: ", err))