/**
 * This script will deploy N new BridgeTokens to an EVM network.
 *
 * Before executing this script, add the following variables to your .env,
 * changing the values to your actual mainnet Alchemy URL and Private Key:
 * MAINNET_URL=https://eth-mainnet.alchemyapi.io/v2/XXXXXXXXXXXXXXXXXXXXXXXX
 * MAINNET_PRIVATE_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
 *
 * Then, create or edit the file data/ibc_tokens_to_deploy.json so that it
 * has only the IbcTokens that you want to deploy.
 * Example:
 * [
 *   {
 *     "name": "Alice Token",
 *     "symbol": "ALI",
 *     "decimals": 10,
 *     "denom": ""
 *   },
 *   {
 *     "name": "Bob Token",
 *     "symbol": "BOB",
 *     "decimals": 18,
 *     "denom": "Bob denom"
 *   }
 * ]
 *
 * Note that the `denom` field is optional. If you don't have that information,
 * you may leave it as an empty string.
 *
 * Finally, run the command `yarn deployIbcTokens:run`.
 *
 * A new file will be created with the results. It's name will be something like
 * data/deployed_ibc_tokens_07_Oct_2021.json, but with today's date.
 *
 * @dev If you want to TEST this script, run the command `yarn deployIbcTokens:test`
 *
 * That's it.
 */

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
const TOKEN_DATA_SOURCE_FILENAME = "./data/ibc_tokens_to_deploy.json";

// Safeguard: if there is an error estimating gasPrice, it cannot go below this value
const MINIMUM_GAS_PRICE_IN_GWEI = 100;

const log = {
  tokens: [],
};

async function main() {
  print("highlight", "~~~ DEPLOY IBC TOKENS ~~~");

  // get tokens from file:
  print("yellow", `📢 Getting token data from ${TOKEN_DATA_SOURCE_FILENAME}`);
  const data = fs.readFileSync(TOKEN_DATA_SOURCE_FILENAME, "utf8");
  const tokensToDeploy = JSON.parse(data);

  print("yellow", `📢 Will deploy ${tokensToDeploy.length} tokens`);

  // Get the current account
  const accounts = await ethers.getSigners();
  const activeAccount = accounts[0];

  // If we're forking, we want to impersonate the owner account
  if (USE_FORKING) {
    await setupForking();
  } else {
    print("cyan", `🤵 Active account is ${activeAccount.address}`);
  }

  // Create an instance of BridgeBank from the deployed code, to have access to its address
  const { instance: bridgeBank } = await support.getDeployedContract(
    DEPLOYMENT_NAME,
    "BridgeBank",
    CHAIN_ID
  );

  const startingBalance = await ethers.provider.getBalance(
    activeAccount.address
  );

  const gasPrice = await estimateGasPrice(MINIMUM_GAS_PRICE_IN_GWEI);

  // Deploy each token in the list
  for (let i = 0; i < tokensToDeploy.length; i++) {
    const tokenData = tokensToDeploy[i];
    await deployToken({
      name: tokenData.name,
      symbol: tokenData.symbol,
      decimals: tokenData.decimals,
      denom: tokenData.denom,
      deployer: activeAccount,
      bridgeBank,
      gasPrice,
    });
  }

  // Calculate and log the total cost
  await logTotalCost(activeAccount, startingBalance);

  // Write logs to the data folder
  logFilename = saveToFile();

  print("yellow", `🧾 Results have been written to ${logFilename}`);
  print("highlight", "~~~ DONE! ~~~");
}

async function deployToken({
  name,
  symbol,
  decimals,
  denom,
  deployer,
  bridgeBank,
  gasPrice,
}) {
  sanityCheck({ name, symbol, decimals, deployer });

  print("yellow", `🔑 Deploying ${name}. Please wait...`);

  // Deploy the token
  const bridgeTokenFactory = await ethers.getContractFactory("BridgeToken");
  const bridgeToken = await bridgeTokenFactory.deploy(
    name,
    symbol,
    decimals,
    denom,
    { gasPrice }
  );
  await bridgeToken.deployed();
  print(
    "green",
    `✅ 1/5: The token ${name} is deployed at ${bridgeToken.address}`
  );

  // Grant the minter role to BridgeBank
  const grantMinterTx = await bridgeToken.grantRole(
    MINTER_ROLE,
    bridgeBank.address,
    { gasPrice }
  );
  print("green", `✅ 2/5: Grant Minter TX: ${grantMinterTx.hash}`);

  // Grant the admin role to BridgeBank
  const grantAdminTx = await bridgeToken.grantRole(
    ADMIN_ROLE,
    bridgeBank.address,
    { gasPrice }
  );
  print("green", `✅ 3/5: Grant Admin TX: ${grantAdminTx.hash}`);

  // Renounce the minter role
  const renounceMinterTx = await bridgeToken.renounceRole(
    MINTER_ROLE,
    deployer.address,
    { gasPrice }
  );
  print("green", `✅ 4/5: Renounce Minter TX: ${renounceMinterTx.hash}`);

  // Renounce the admin role
  const renounceAdminTx = await bridgeToken.renounceRole(
    ADMIN_ROLE,
    deployer.address,
    { gasPrice }
  );
  print("green", `✅ 5/5: Renounce Admin TX: ${renounceAdminTx.hash}`);

  const receipt = {
    name: name,
    symbol: symbol,
    decimals: decimals,
    denom: denom,
    address: web3.utils.toChecksumAddress(bridgeToken.address),
    grantMinterRoleTx: grantMinterTx.hash,
    grantAdminRoleTx: grantAdminTx.hash,
    renounceMinterRoleTx: renounceMinterTx.hash,
    renounceAdminRoleTx: renounceAdminTx.hash,
  };
  log.tokens.push(receipt);
}

async function setupForking() {
  // Impersonate the admin account
  await support.impersonateAccount(
    support.PROXY_ADMIN_ADDRESS,
    "10000000000000000000",
    "Proxy Admin"
  );
}

function sanityCheck({ name, symbol, decimals, deployer }) {
  if (!deployer) {
    throw new Error("💥 CRITICAL: MISSING DEPLOYER ACCOUNT");
  }
  if (!name || !symbol || !decimals) {
    print(
      "h_red",
      `Missing token data! Name: ${name} | symbol: ${symbol} | decimals: ${decimals}`
    );
  }
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

function saveToFile() {
  // write the log to a file named after today
  const logFileName = generateTodayFilename({
    prefix: "deployed_ibc_tokens",
    extension: "json",
    directory: "data",
  });
  fs.writeFileSync(logFileName, JSON.stringify(log, null, 1));

  return logFileName;
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
    logFilename = saveToFile();
    print(
      "yellow",
      `🧾 There was an error, but results have been written to ${logFilename}`
    );
    process.exit(0);
  })
  .finally(() => process.exit(0));
