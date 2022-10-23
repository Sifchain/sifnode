import { HardhatRuntimeEnvironment } from "hardhat/types";
// import { Blocklist } from "../../build";
import { Wallet } from "ethers";
import { FetchWallet } from "../../scripts/helpers/KeyHandler";
import { Contract } from "ethers";
import {print} from "../../scripts/helpers/utils";

/**
 * Get a blocklist contract connected to the address provided
 * @param hre The hardhat Runtime Environment
 * @param contractAddress The address of the blocklist to connect to
 * @param wallet The ethers wallet to interact with the contract with
 * @returns a promise of the Blocklist object
 */
async function GetBlocklist(hre: HardhatRuntimeEnvironment, contractAddress: string, wallet: Wallet): Promise<Contract> {
    print("yellow", `Connecting to Blocklist at address ${contractAddress}`);
    const blocklistFactory = await hre.ethers.getContractFactory("Blocklist", wallet);
    const blocklist = await blocklistFactory.attach(contractAddress);
    print("success", `Connected to blocklist`);
    return blocklist;
}

/**
 * Add a user to the blocklist at provided contract address
 * @param hre The Hardhat Runtime Environment
 * @param contractAddress The address of the Blocklist Contract
 * @param user The address of the user to block
 * @returns True on success, False on failure
 */
export async function AddUser(hre: HardhatRuntimeEnvironment, contractAddress: string, user: string, walletName: string, walletPassword: string): Promise<boolean> {
    print("white", `Attempting to add user ${user} to blocklist at address ${contractAddress}`);
    try {
        const wallet = await FetchWallet(hre, walletName, walletPassword);
        if (wallet === false) {
            print("error", `Unable to open wallet named ${walletName}`);
            return false;
        }
        const blocklist = await GetBlocklist(hre, contractAddress, wallet);
        if (await blocklist.isBlocklisted(user)) {
            print("warn",`User (${user}) was already in Blocklist (${contractAddress}), doing nothing`);
            // If user is already blocked we can return success
            return true;
        }
        await blocklist.addToBlocklist(user);
        print("bigSuccess", "User added to the blocklist successfully");
        return true;
    } catch (error) {
        print("error", `Unable to add user to blocklist, error: ${error}`);
        return false;
    }
}

/**
 * Remove a user from the blocklist at provided contract address 
 * @param hre The Hardhat Runtime Environment
 * @param contractAddress The address of the Blocklist Contract
 * @param user The address of the user to remove from the blocklist
 * @returns True on success, False on failure
 */
export async function RemoveUser(hre: HardhatRuntimeEnvironment, contractAddress: string, user: string, walletName: string, walletPassword: string): Promise<boolean> {
    print("white", `Attempting to remove user ${user} from blocklist at address ${contractAddress}`);
    try {
        const wallet = await FetchWallet(hre, walletName, walletPassword);
        if (wallet === false) {
            print("error", `Unable to open wallet named ${walletName}`);
            return false;
        }
        const blocklist = await GetBlocklist(hre, contractAddress, wallet);
        if (!await blocklist.isBlocklisted(user)) {
            print("warn",`User (${user}) is not currently in Blocklist (${contractAddress}), doing nothing`);
            // If user is already not blocklisted we can return success
            return true;
        }
        await blocklist.removeFromBlocklist(user);
        print("bigSuccess", "User removed from the blocklist successfully");
        return true;
    } catch (error) {
        print("error", `Unable to remove user from blocklist, error: ${error}`);
        return false;
    }
}