require("dotenv").config();

const web3 = require("web3");
const { ethers } = require("hardhat");
const support = require("./helpers/forkingSupport");
const { print } = require("./helpers/utils");

const USE_FORKING = !!process.env.USE_FORKING;

// If there is no DEPLOYMENT_NAME env var, we'll use the mainnet deployment
const DEPLOYMENT_NAME = process.env.DEPLOYMENT_NAME || "sifchain-1";

// If there is no FORKING_CHAIN_ID env var, we'll use the mainnet id
const CHAIN_ID = process.env.FORKING_CHAIN_ID || 1;

// Roles to grant and renounce
const MINTER_ROLE = web3.utils.soliditySha3("MINTER_ROLE");
const ADMIN_ROLE =
  "0x0000000000000000000000000000000000000000000000000000000000000000";

const state = {
  token: {
    name: process.env.TOKEN_NAME,
    symbol: process.env.TOKEN_SYMBOL,
    decimals: process.env.TOKEN_DECIMALS,
    denom: process.env.TOKEN_DENOM || "",
  },
};

async function main() {
  print("highlight", "~~~ DEPLOY IBC TOKENS ~~~");

  // Guarantee we have all the needed variables
  sanityCheck();

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

  // Deploy the token
  const bridgeTokenFactory = await ethers.getContractFactory("BridgeToken");
  const bridgeToken = await bridgeTokenFactory.deploy(
    state.token.name,
    state.token.symbol,
    state.token.decimals,
    state.token.denom
  );
  await bridgeToken.deployed();
  print(
    "green",
    `✅ 1/5: The token ${state.token.name} is deployed at ${bridgeToken.address}`
  );

  // Grant the minter role to BridgeBank
  const grantMinterTx = await bridgeToken.grantRole(
    MINTER_ROLE,
    bridgeBank.address
  );
  print("green", `✅ 2/5: Grant Minter TX: ${grantMinterTx.hash}`);

  // Grant the admin role to BridgeBank
  const grantAdminTx = await bridgeToken.grantRole(
    ADMIN_ROLE,
    bridgeBank.address
  );
  print("green", `✅ 3/5: Grant Admin TX: ${grantAdminTx.hash}`);

  // Renounce the minter role
  const renounceMinterTx = await bridgeToken.renounceRole(
    MINTER_ROLE,
    activeAccount.address
  );
  print("green", `✅ 4/5: Renounce Minter TX: ${renounceMinterTx.hash}`);

  // Renounce the admin role
  const renounceAdminTx = await bridgeToken.renounceRole(
    ADMIN_ROLE,
    activeAccount.address
  );
  print("green", `✅ 5/5: Renounce Admin TX: ${renounceAdminTx.hash}`);

  // Should estimate gas?!

  print("highlight", "~~~ DONE! ~~~");
}

async function setupForking() {
  // Impersonate the admin account
  state.proxyAdmin = await support.impersonateAccount(
    support.PROXY_ADMIN_ADDRESS,
    "10000000000000000000",
    "Proxy Admin"
  );
}

function sanityCheck() {
  if (!state.token.name || !state.token.symbol || !state.token.decimals) {
    print(
      "h_red",
      `Token name: ${state.token.name} | Token symbol: ${state.token.symbol} | Token decimals: ${state.token.decimals}`
    );
    throw new Error("💥 CRITICAL: MISSING ENV VARIABLES");
  }
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
