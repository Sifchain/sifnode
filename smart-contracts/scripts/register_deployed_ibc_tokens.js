require("dotenv").config();

const fs = require("fs");
const web3 = require("web3");
const { ethers } = require("hardhat");

const support = require("./helpers/forkingSupport");
const { print, estimateGasPrice } = require("./helpers/utils");

const USE_FORKING = !!process.env.USE_FORKING;

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.FORKING_CHAIN_ID || 1;

// Will get the token data from this file
// should be something like "./data/deployed_ibc_tokens_08_Oct_2021.json";
const TOKEN_DATA_SOURCE_FILENAME = process.env.REGISTER_TOKENS_SOURCE_FILENAME;

// Safeguard: if there is an error estimating gasPrice, it cannot go below this value
const MINIMUM_GAS_PRICE_IN_GWEI = 100;

let log = {};

async function main() {
  print("highlight", "~~~ REGISTER IBC TOKENS ~~~");

  // get tokens from file:
  print("yellow", `📢 Getting token data from ${TOKEN_DATA_SOURCE_FILENAME}`);
  const data = fs.readFileSync(TOKEN_DATA_SOURCE_FILENAME, "utf8");
  log = JSON.parse(data);

  print("yellow", `📢 Will register ${log.tokens.length} tokens`);

  // Create an instance of BridgeBank from the deployed code
  const { instance: bridgeBank } = await support.getDeployedContract(
    DEPLOYMENT_NAME,
    "BridgeBank",
    CHAIN_ID
  );

  // Get the current account
  const accounts = await ethers.getSigners();
  let activeAccount = accounts[0];

  // If we're forking, we want to impersonate the owner account
  if (USE_FORKING) {
    const signerOwner = await setupForking(bridgeBank);
    activeAccount = signerOwner;
  } else {
    print("cyan", `🤵 Active account is ${activeAccount.address}`);
  }

  // Get the current balance to be able to calculate how much it all cost
  const startingBalance = await ethers.provider.getBalance(
    activeAccount.address
  );

  // Get the ideal gasPrice
  const gasPrice = await estimateGasPrice(MINIMUM_GAS_PRICE_IN_GWEI);

  // Send transactions
  for (let i = 0; i < log.tokens.length; i++) {
    const tokenData = log.tokens[i];
    const tx = await bridgeBank
      .connect(activeAccount)
      .addExistingBridgeToken(tokenData.address, {
        gasPrice,
      });

    // Add this tx to the file
    tokenData.registerTx = tx.hash;

    print("green", `✅ ${tokenData.name} TX: ${tx.hash}`);
  }

  // Calculate and log the total cost
  await logTotalCost(activeAccount, startingBalance);

  // Save results to the souce file
  saveToFile();

  print("highlight", "~~~ DONE! ~~~");
}

// Will add the total cost of this script with the total cost from the deployment script
async function logTotalCost(account, startingBalance) {
  const endingBalance = await ethers.provider.getBalance(account.address);

  const cost = ethers.BigNumber.from(startingBalance).sub(
    ethers.BigNumber.from(endingBalance)
  );
  const costInEth = web3.utils.fromWei(cost.toString());

  const previousCost = web3.utils.toWei(log.totalCost);
  const totalCost = ethers.BigNumber.from(cost)
    .add(ethers.BigNumber.from(previousCost))
    .toString();

  const totalCostInEth = web3.utils.fromWei(totalCost);

  // Add to the previous cost to save to file
  log.totalCost = totalCostInEth;

  print("cyan", "----");
  print("cyan", `💵 Total ETH spent: ${costInEth}`);
  print("cyan", "----");
}

// Impersonate the owner account
async function setupForking(contractInstance) {
  const owner = await contractInstance.owner();

  const signerOwner = await support.impersonateAccount(
    owner,
    "10000000000000000000",
    "BridgeBank Owner"
  );

  return signerOwner;
}

function saveToFile() {
  fs.writeFileSync(TOKEN_DATA_SOURCE_FILENAME, JSON.stringify(log, null, 1));
  print(
    "yellow",
    `🧾 Results have been written to ${TOKEN_DATA_SOURCE_FILENAME}`
  );
}

function treatCommonErrors(e) {
  if (e.message.indexOf("Unsupported method") !== -1) {
    print(
      "h_red",
      "Error: if you are NOT trying to test this with a mainnet fork, please remove the variable USE_FORKING from your .env"
    );
  } else if (e.message.indexOf("insufficient funds") !== -1) {
    print(
      "h_red",
      "Error: insufficient funds. If you are using the correct private key, please refill your account with EVM native coins."
    );
  } else if (e.message.indexOf("caller is not the owner") !== -1) {
    print(
      "h_red",
      "Error: caller is not the owner. Either you have the wrong private key set in your .env, or you should add USE_FORKING=1 to your .env if you want to test the script."
    );
  } else {
    print("h_red", e.stack);
  }
}

main()
  .catch((error) => {
    treatCommonErrors(error);
    saveToFile();
    process.exit(0);
  })
  .finally(() => process.exit(0));
