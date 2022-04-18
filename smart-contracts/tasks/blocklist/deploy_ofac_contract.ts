import {FetchWallet, GenerateWallet} from "../../scripts/helpers/KeyHandler";
import {hasSameElementsAndLength, print, sleep} from "../../scripts/helpers/utils";
import {BigNumberish, Wallet} from "ethers";
import {getList} from "./ofacParser";
import { HardhatRuntimeEnvironment } from "hardhat/types";
/**
 * Function waits for a wallet to be filled with the minimum balance amount, returns the balance that is
 * above or equal to minimum balance.
 * @param wallet The wallet to query balance on
 * @param minimumBalance The balance to resolve on
 * @returns The balance in the wallet after funding
 */
async function WaitForFunding(wallet: Wallet, minimumBalance: BigNumberish): Promise<BigNumberish> {
    for(;;) {
        const balance = await wallet.getBalance();
        if (balance >= minimumBalance) {
            return balance;
        }
        // Wait 15 seconds between balance checks to not get rate limited
        await sleep(15_000);
    }
}

/**
 * Async function that will take a wallet name, a password, and what balance it should trigger on. Once run it will fetch/create a wallet 
 * of the given name and password, and then will monitor the balance of the account for the minimum balance. Once it detects the minimum 
 * balance it will deploy the OFAC blocklist contracts and then sync the list with the current sanctioned addresses from the OFAC website.
 * @param hre The hardhat runtime environment
 * @param walletName The name of the wallet to fetch/generate
 * @param password The password to encrypt/decrypt the wallet with
 * @param minimumBalance The balance at which to monitor for
 * @returns Promise<void>
 */
export async function SetupOFACBlocklist(hre: HardhatRuntimeEnvironment, walletName: string, password: string, minimumBalance: BigNumberish, ofacURL: string) {   
    const ethers = hre.ethers;
    print("white", "Attempting to generate/fetch Private Key");
    const address = await GenerateWallet(hre, walletName, password);
    const wallet = await FetchWallet(hre, walletName, password);
    if (address === false || wallet === false) {
        print("error", "ERROR: Could not generate/read wallet file");
        return;
    }
    print("success", `Generated/Fetched OFAC wallet named: ${walletName}, public address of wallet is: ${address}`);
    let prettyBalance = ethers.utils.formatEther(minimumBalance);
    print("h_yellow", `Waiting for account to be funded with ${prettyBalance} ETH before continuing`);
    const balance = await WaitForFunding(wallet, minimumBalance);
    prettyBalance = ethers.utils.formatEther(balance);
    print("success", `Account has been funded with ${prettyBalance} ETH, starting the deploy process`);
    const blocklistFactory = await ethers.getContractFactory("Blocklist", wallet);
    try {
        const blocklist = await blocklistFactory.deploy();
        print("success", `Blocklist was deployed to EVM chain, Contract Address is ${blocklist.address}`);
        const ofacAddresses = await getList(ofacURL);
        await blocklist.batchAddToBlocklist(ofacAddresses);
        const addedAddresses = await blocklist.getFullList();
        if (!hasSameElementsAndLength(ofacAddresses, addedAddresses)) {
            print("error", "Not all OFAC addresses where added to the blocklist, manually investigate what has gone wrong.");
            return;
        }
        print("success", `All ${addedAddresses.length} OFAC sanctioned addresses have been added to the blocklist contract`);
        print("yellow", `Sanctioned Addresses: [${await blocklist.getFullList()}]`);
        prettyBalance = ethers.utils.formatEther(await wallet.getBalance());
        print("cyan", `Remaining ethereum after deploy and list sync: ${prettyBalance} ETH`);
        print("bigSuccess", `Blocklist deploy and setup was successful, contract address: ${blocklist.address}`);
    } catch (error) {
        print("error", "Error while attempting to deploy blocklist contract");
        console.error(error);
    }
}