import { promises } from "fs";
import path from "path";
import os from "os";
import { Wallet } from "ethers";
import { HardhatRuntimeEnvironment } from "hardhat/types";

const KeyDir = `${os.homedir()}/.sifnode/wallets/`;

function GetWalletDir(walletName: string): string {
    const walletPath = `${KeyDir}${walletName.toLowerCase()}`;
    return path.resolve(walletPath)
}

/**
 * Generates a new wallet and stores the private data under the users home directory 
 * @param hre The hardhat runtime environment so that this code can be run inside hardhat tasks
 * @param walletName Name of the wallet to create
 * @param password A password to encrypt the wallet with (required to read it back in)
 * @returns (string) Wallet address or False if the wallet could not be created
 */
export async function GenerateWallet(hre: HardhatRuntimeEnvironment, walletName: string, password: string = ""): Promise<string | false> {
    // Get the directory wallets are stored
    const walletDir = GetWalletDir(walletName);
    const walletFile = `${walletDir}/key`;
    // Check if key exists first before generating a wallet over it
    try {
        await promises.access(walletFile);
        console.log("Wallet already exists, not overriding", walletFile);
        // If we found an already existing wallet we will report that wallets address
        const wallet = await FetchWallet(hre, walletName, password)
        if (wallet === false) {
            // If the wallet does not parse we should act as though wallet generation has failed
            return false
        } else {
            return wallet.address;
        }
    } catch {
        // We could not find the wallet so its safe to generate one
    }
    // Generate a Sifnode Key Directory if it does not exist
    try {
        await promises.mkdir(walletDir, {mode: 0o700, recursive: true});
        // Generate the new wallet
        const wallet = await ethers.Wallet.createRandom()
        // Generate the private key either from 
        const jsonWallet = await wallet.encrypt(password);
        // Store the new wallet file with read only permissions
        await promises.writeFile(walletFile, jsonWallet, {mode: 0o400});
        // Disable writing permissions to the key directory
        await promises.chmod(walletDir, 0o500);
        // Return the wallets public address
        return wallet.address;
    } catch (error) {
            return false;
        }
}

/**
 * Reads a stored private key from storage and generates a wallet to use for signing transactions and
 * deploying code.
 * @param hre The hardhat runtime environment so that this code can be run inside hardhat tasks
 * @param walletName Name of the wallet to create
 * @param password A password to decrypt the wallet with (required if wallet was generated with a password)
 * @returns false if wallet could not be opened, otherwise returns a ethers wallet instance
 */
export async function FetchWallet(hre: HardhatRuntimeEnvironment, walletName: string, password: string = ""): Promise<Wallet | false> {
    // Get hardhat ethers instance
    const ethers = hre.ethers;
    // Lookup the Sifnode Key Directory
    const walletDir = GetWalletDir(walletName);
    // Read a private key file into a wallet
    try {
        const privateKey = String(await promises.readFile(`${walletDir}/key`));
        // We decrypt the wallet then use this wallet to create a ethers wallet file that has the hardhat provider
        const wallet = await ethers.Wallet.fromEncryptedJson(privateKey, password);
        // Regenerate the wallet with the hardhat provider so that we can send things over the network through hardhat
        return new ethers.Wallet(wallet.privateKey, ethers.provider)
    } catch {
        // Could not find a wallet for that name, return False
        return false;
    }
}