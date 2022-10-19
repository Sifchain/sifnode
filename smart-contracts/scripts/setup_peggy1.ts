import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { ethers, upgrades } from "hardhat";
import { BridgeBank, CosmosBridge } from "../build/"

const ETH = (value: string) => ethers.utils.parseEther(value);

async function fundAccounts(admin: SignerWithAddress, ...others: SignerWithAddress[]) {
    const adminBalance = await admin.getBalance();
    if (adminBalance.lt(ETH("3"))) {
        throw Error(`Insufficient Balance to deploy contracts! Current Balance: ${ethers.utils.formatEther(adminBalance)} ETH`);
    }
    console.log(`Admin account hash ETH balance of ${ethers.utils.formatEther(adminBalance)} ETH which is sufficient to fund all accounts and continue.`);
    for (const account of others) {
        const balance = await account.getBalance();
        if (balance.lt(ETH("0.25"))) {
            // If account has less then 0.25 ETH send the difference to bring them to 0.25 from the admin account
            const diff = ETH("0.25").sub(balance);
            await admin.sendTransaction({from: admin.address, to: account.address, value: diff});
        }
    }
}

async function setup() {
    const [admin, operator, owner, pauser, blocklistAdmin] = await ethers.getSigners();
    // Make sure accounts are funded, if not transfer funds from the admin
    console.log("Funding Accounts");
    await fundAccounts(admin, operator, owner, pauser, blocklistAdmin);
    console.log("All accounts are funded");

    console.log("Deploying Cosmos Bridge");
    const consensusThreshold = 70 // Percentage of relayers to confirm a transaction
    const cosmosBridgeFactory = await ethers.getContractFactory("CosmosBridge", admin);
    const cosmosBridge = await upgrades.deployProxy(cosmosBridgeFactory, [operator.address, consensusThreshold, [], []]) as CosmosBridge;
    console.log("CosmosBridge deployed at: ", cosmosBridge.address);

    console.log("Deploying BridgeBank");
    const bridgeBankFactory = await ethers.getContractFactory("BridgeBank", admin);
    const bridgeBank = await upgrades.deployProxy(
        bridgeBankFactory, 
        [
            operator.address, 
            cosmosBridge.address, 
            owner.address, 
            pauser.address
        ],
        {
            initializer: "initialize(address,address,address,address)"
        }
    ) as BridgeBank;
    console.log("BridgeBank deployed at: ", bridgeBank.address);

    console.log("Setting Bridgebank in cosmosbridge");
    const setBridgeBank = await cosmosBridge.connect(operator).setBridgeBank(bridgeBank.address);
    console.log("Bridgebank successfully set in cosmosBridge, TX: ", setBridgeBank.hash);

    
    console.log("Deploying eRowan ERC20 Token");
    const eRowanFactory = await ethers.getContractFactory("BridgeToken", admin);
    const eRowan = await eRowanFactory.deploy("erowan");
    console.log("eRowan token deployed at address: ", eRowan.address);
    
    console.log("Minting an initial supply of eRowan to the admin account");
    const mintTx = await eRowan.mint(admin.address, ETH("100000000"));
    console.log("Initial minting of eRowan successful, Tx: ", mintTx.hash);

    console.log("Adding the bridgeBank as a minter of eRowan");
    const addBridgeBankTx = await eRowan.addMinter(bridgeBank.address);
    console.log("eRowan add new minter successful, Tx: ", addBridgeBankTx.hash);

    console.log("Renouncing the minter role of the admin address so only bridgebank can mint eRowan");
    const renounceTx = await eRowan.renounceMinter();
    console.log("Renouncement successful, Tx: ", renounceTx.hash);

    console.log("Adding eRowan as an existing Bridge Token of BridgeBank");
    const addERowanTx = await bridgeBank.connect(owner).addExistingBridgeToken(eRowan.address);
    console.log("Adding successfull, Tx: ", addERowanTx.hash);

    console.log("Deploying and initializing BridgeRegistry Contract");
    const bridgeRegistryFactory = await ethers.getContractFactory("BridgeRegistry");
    const bridgeRegistry = await bridgeRegistryFactory.deploy()
    const initBRTX = await bridgeRegistry.initialize(cosmosBridge.address, bridgeBank.address);
    console.log("BridgeRegistry deployed to address: ", bridgeRegistry.address);
    console.log("BridgeRegistry Init transaction submitted successfully, Tx: ", initBRTX.hash);

    console.log("Deploying Blocklist Contract");
    const blocklistFactory = await ethers.getContractFactory("Blocklist", blocklistAdmin);
    const blocklist = await blocklistFactory.deploy()
    console.log("Blocklist address set as: ", blocklist.address);
   
    console.log("Setting Blocklist in BridgeBank");
    const setBlocklistTx = await bridgeBank.connect(operator).setBlocklist(blocklist.address);
    console.log("BridgeBank set blocklist sucessfully submitted, TX: ", setBlocklistTx.hash);
}

setup().then(() => console.log('Setup script has completed')).catch((error) => console.error("Error encountered: ", error));