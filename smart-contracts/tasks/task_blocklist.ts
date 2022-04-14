require("dotenv").config();

import { task, types } from "hardhat/config";
import { AddUser, RemoveUser } from "./blocklist/blocklist_operations";
import {SetupOFACBlocklist} from "./blocklist/deploy_ofac_contract";
import { SyncOfacBlocklist } from "./blocklist/sync_ofac_blocklist";

interface BlocklistGeneralArgs {
    wallet: string;
    password: string;
}

interface BlocklistDeployArgs extends BlocklistGeneralArgs {
    minimumBalance: string;
    ofacWebsite: string;
}

interface BlocklistSyncArgs extends BlocklistGeneralArgs {
    ofacWebsite: string;
    blocklistAddress: string;
}

interface BlocklistOperationArgs extends BlocklistGeneralArgs {
    blocklistAddress: string;
    userAddress: string;
}

task("blocklist:deploy", "Deploy and sync a new blocklist smart contract")
    .addOptionalParam("wallet", "The name of the wallet to generate/fetch. Defaults to 'blocklist'", "blocklist", types.string)
    .addOptionalParam("password", "The password to encrypt/decrypt the wallet with. Can also be set with WALLET_PASSWORD env variable.", "", types.string)
    .addOptionalParam("minimumBalance", "The minimum balance the scripts should look for before performing operations. Defaults to 2 ETH.", 2.0, types.float)
    .addOptionalParam("ofacWebsite", "The website to sync sanctioned addresses from. Defaults to `https://www.treasury.gov/ofac/downloads/sdnlist.txt`", "https://www.treasury.gov/ofac/downloads/sdnlist.txt", types.string)
    .setAction(async (args: BlocklistDeployArgs, hre) => {
        const password = process.env.WALLET_PASSWORD || args.password;
        await SetupOFACBlocklist(hre, args.wallet, password, hre.ethers.utils.parseEther(String(args.minimumBalance)), args.ofacWebsite);
    });

task("blocklist:sync", "Tools related to the deploying, management and syncing of OFAC blocklists")
    .addPositionalParam("blocklistAddress", "The contract address of the Blocklist")
    .addOptionalParam("wallet", "The name of the wallet to generate/fetch. Defaults to 'blocklist'", "blocklist", types.string)
    .addOptionalParam("password", "The password to encrypt/decrypt the wallet with. Can also be set with WALLET_PASSWORD env variable.", "", types.string)
    .addOptionalParam("ofacWebsite", "The website to sync sanctioned addresses from. Defaults to `https://www.treasury.gov/ofac/downloads/sdnlist.txt`", "https://www.treasury.gov/ofac/downloads/sdnlist.txt", types.string)
    .setAction(async (args: BlocklistSyncArgs, hre) => {
        const password = process.env.WALLET_PASSWORD || args.password;
        await SyncOfacBlocklist(hre, args.blocklistAddress, args.wallet, password, args.ofacWebsite);
    });

task("blocklist:add", "Tools related to the deploying, management and syncing of OFAC blocklists")
    .addPositionalParam("blocklistAddress", "The contract address of the Blocklist")
    .addPositionalParam("userAddress", "The address of the user to add to the blocklist")
    .addOptionalParam("wallet", "The name of the wallet to generate/fetch. Defaults to 'blocklist'", "blocklist", types.string)
    .addOptionalParam("password", "The password to encrypt/decrypt the wallet with. Can also be set with WALLET_PASSWORD env variable.", "", types.string)
    .setAction(async (args: BlocklistOperationArgs, hre) => {
        const password = process.env.WALLET_PASSWORD || args.password;
        await AddUser(hre, args.blocklistAddress, args.userAddress, args.wallet, password);

    });

task("blocklist:remove", "Tools related to the deploying, management and syncing of OFAC blocklists")
    .addPositionalParam("blocklistAddress", "The contract address of the Blocklist")
    .addPositionalParam("userAddress", "The address of the user to remove from the blocklist")
    .addOptionalParam("wallet", "The name of the wallet to generate/fetch. Defaults to 'blocklist'", "blocklist", types.string)
    .addOptionalParam("password", "The password to encrypt/decrypt the wallet with. Can also be set with WALLET_PASSWORD env variable.", "", types.string)
    .setAction(async (args: BlocklistOperationArgs, hre) => {
        const password = process.env.WALLET_PASSWORD || args.password;
        await RemoveUser(hre, args.blocklistAddress, args.userAddress, args.wallet, password);
    });