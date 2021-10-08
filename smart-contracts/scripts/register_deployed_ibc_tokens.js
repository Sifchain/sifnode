require("dotenv").config();

const fs = require("fs");
const web3 = require("web3");
const { ethers } = require("hardhat");

const support = require("./helpers/forkingSupport");
const {
  print,
  generateTodayFilename,
  estimateGasPrice,
  MINTER_ROLE,
  ADMIN_ROLE,
} = require("./helpers/utils");

const USE_FORKING = !!process.env.USE_FORKING;

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.FORKING_CHAIN_ID || 1;

// Will get the token data from this file
// TODO get it from .env
const TOKEN_DATA_SOURCE_FILENAME =
  "./data/deployed_ibc_tokens_07_Oct_2021.json";

// Safeguard: if there is an error estimating gasPrice, it cannot go below this value
const MINIMUM_GAS_PRICE_IN_GWEI = 100;

const log = {};

async function main() {
  print("highlight", "~~~ DEPLOY IBC TOKENS ~~~");

  // get tokens from file:
  print("yellow", `📢 Getting token data from ${TOKEN_DATA_SOURCE_FILENAME}`);
  const data = fs.readFileSync(TOKEN_DATA_SOURCE_FILENAME, "utf8");
  const tokensToRegister = JSON.parse(data).tokens;

  print("yellow", `📢 Will register ${tokensToRegister.length} tokens`);

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
  let signerOwner;
  if (USE_FORKING) {
    signerOwner = await setupForking(bridgeBank);
    activeAccount = signerOwner;
  } else {
    signerOwner = ethers.getSigner(activeAccount.address);
    print("cyan", `🤵 Active account is ${activeAccount.address}`);
  }

  const startingBalance = await ethers.provider.getBalance(
    activeAccount.address
  );

  const gasPrice = await estimateGasPrice(MINIMUM_GAS_PRICE_IN_GWEI);

  const txs = [];
  for (let i = 0; i < tokensToRegister.length; i++) {
    const tokenData = tokensToRegister[i];
    const tx = await bridgeBank
      .connect(signerOwner)
      .addExistingBridgeToken(tokenData.address, {
        gasPrice,
      });

    txs.push({
      token: tokenData.name,
      registerTx: tx.hash,
    });

    print("green", `✅ ${tokenData.name} TX: ${tx.hash}`);
  }

  // Calculate and log the total cost
  await logTotalCost(activeAccount, startingBalance);

  print("highlight", "~~~ DONE! ~~~");
}

async function logTotalCost(account, startingBalance) {
  const endingBalance = await ethers.provider.getBalance(account.address);
  const totalCost = ethers.BigNumber.from(startingBalance).sub(
    ethers.BigNumber.from(endingBalance)
  );

  log.totalCost = web3.utils.fromWei(totalCost.toString());
  print("cyan", "----");
  print("cyan", `💵 Total ETH spent: ${log.totalCost}`);
  print("cyan", "----");
}

async function setupForking(contractInstance) {
  const owner = await contractInstance.owner();

  // Impersonate the admin account
  const signerOwner = await support.impersonateAccount(
    owner,
    "10000000000000000000",
    "BridgeBank Owner"
  );

  return signerOwner;
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
  })
  .finally(() => process.exit(0));
