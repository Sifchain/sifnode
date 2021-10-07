/**
 * This script will deploy a new BridgeToken to an EVM network.
 // TODO runbook
 // TODO get list from file
 // TODO write the attach script
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

const MINIMUM_GAS_PRICE_IN_GWEI = 100;

const log = {
  tokens: [],
};

const tokensToDeploy = [
  {
    name: "Dan Token",
    symbol: "DAN",
    decimals: 6,
    denom: "",
  },
  {
    name: "Alice Token",
    symbol: "ALI",
    decimals: 10,
    denom: "Alice denom",
  },
  {
    name: "Bruce Token",
    symbol: "BRU",
    decimals: 18,
    denom: "Bruce denom",
  },
];

async function main() {
  print("highlight", "~~~ DEPLOY IBC TOKENS ~~~");

  print("yellow", `📢 Will deploy ${tokensToDeploy.length} tokens.`);

  // If we're forking, we want to impersonate the owner account
  if (USE_FORKING) await setupForking();

  // Create an instance of BridgeBank from the deployed code, to have access to its address
  const { instance: bridgeBank } = await support.getDeployedContract(
    DEPLOYMENT_NAME,
    "BridgeBank",
    CHAIN_ID
  );

  // Get the current account
  const accounts = await ethers.getSigners();
  const activeAccount = accounts[0];
  print("cyan", `🤵 Active account is ${activeAccount.address}`);

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

  print("yellow", `🔑 Deploying ${name} token. Please wait...`);

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
    address: bridgeToken.address,
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
  print("cyan", `Total ETH spent: ${log.totalCost}`);
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
  })
  .finally(() => process.exit(0));
