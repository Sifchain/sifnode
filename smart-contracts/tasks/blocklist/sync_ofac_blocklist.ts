require("dotenv").config();

import {print} from "../../scripts/helpers/utils";
import {getList} from "./ofacParser";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BigNumberish, Wallet, Contract } from "ethers";
// import { Blocklist } from "../../build";
import { HardhatRuntimeEnvironment } from "hardhat/types";
import {FetchWallet} from "../../scripts/helpers/KeyHandler";

// Defaults to the Ethereum Mainnet address
const BLOCKLIST_ADDRESS =
  process.env.BLOCKLIST_ADDRESS || "0x9C8a2011cCb697D7EDe3c94f9FBa5686a04DeACB";

const USE_FORKING = !!process.env.USE_FORKING;

// Will estimate gas and multiply the result by this value (wiggle room)
const GAS_PRICE_BUFFER = 1.2;

interface State {
  ofac: string[];
  evm: string[];
  toAdd: string[];
  toRemove: string[];
  blocklistInstance: Contract;
  activeAccount: Wallet;
  idealGasPrice: BigNumberish;
}

export async function SyncOfacBlocklist(hre: HardhatRuntimeEnvironment, blocklistAddress: string, walletName: string, walletPassword: string, ofacURL: string) {
  print("highlight", "~~~ SYNC OFAC BLOCKLIST ~~~");

  // Fetches lists, compares them and figures out what has to be added or removed
  const state = await setupState(hre, blocklistAddress, walletName, walletPassword, ofacURL);

    print("cyan", `🤵 Active account is ${state.activeAccount.address}`);

  // Add addresses to the blocklist
  await addToBlocklist(state);
  print("cyan", `----`);

  // Remove addresses from the blocklist
  await removeFromBlocklist(state);
  print("cyan", `----`);

  // Print success
  print("h_green", "Our EVM blocklist is synced with OFAC's blocklist");
  print("highlight", "~~~ DONE ~~~");
}

async function setupState(hre: HardhatRuntimeEnvironment, blocklistAddress: string, walletName: string, walletPassword: string, ofacURL: string) : Promise<State> {
  const ethers = hre.ethers;
  
  const wallet = await FetchWallet(hre, walletName, walletPassword)
  if (wallet === false) {
    print("error", "Could not fetch wallet, exiting");
    throw(`Could not fetch walletName: ${walletName}`);
    }
  const activeAccount = wallet
  // Set the deployed blocklist instance
  const blocklistFactory = await ethers.getContractFactory("Blocklist", wallet);
  const blocklistInstance = await blocklistFactory.attach(blocklistAddress);

   // Estimate gasPrice:
  const idealGasPrice = await estimateGasPrice(hre);

  // Set the OFAC list
  const ofac = await getList(ofacURL);
  print("cyan", `OFAC LIST: ${ofac}`);
  print("cyan", `----`);

  // Set the EVM list
  print("yellow", "Fetching EVM blocklist...");
  const evm: string[] = await blocklistInstance.getFullList();
  print("cyan", `EVM LIST : ${evm}`);
  print("cyan", `----`);

  // Find out what the diff between lists is
  print("yellow", "Calculating Diff...");

  // Addresses that must be added don't exist on evm, but exist on ofac
  const toAdd = ofac.filter((address) => !evm.includes(address));
  print("cyan", `Will add: ${toAdd}`);

  // Addresses that must be removed exist on evm, but don't exist on ofac
  const toRemove = evm.filter((address) => !ofac.includes(address));
  print("cyan", `Will remove: ${toRemove}`);
  print("cyan", "----");
  
  return {
    ofac, 
    toAdd, 
    toRemove, 
    idealGasPrice, 
    blocklistInstance, 
    evm,
    activeAccount,
  }
}

async function estimateGasPrice(hre: HardhatRuntimeEnvironment) {
  console.log("Estimating ideal Gas price, please wait...");

  const gasPrice = await hre.ethers.provider.getGasPrice();
  const finalGasPrice = Math.round(gasPrice.toNumber() * GAS_PRICE_BUFFER);

  console.log(`Using ideal Gas price: ${hre.ethers.utils.formatUnits(finalGasPrice, "gwei")} GWEI`);

  return finalGasPrice;
}

async function addToBlocklist(state: State) {
  if (state.toAdd.length === 0) {
    print("yellow", "The are no new addresses to add to the blocklist");
    return;
  }

  print("yellow", "Adding addresses to the blocklist. Please wait...");

  let tx;
  if (state.toAdd.length === 1) {
    tx = await state.blocklistInstance
      .connect(state.activeAccount)
      .addToBlocklist(state.toAdd[0], { gasPrice: state.idealGasPrice, gasLimit: 6000000 })
      .catch((e: Error) => {
        throw e;
      });
  } else {
    // there are many addresses to add
    tx = await state.blocklistInstance
      .connect(state.activeAccount)
      .batchAddToBlocklist(state.toAdd, { gasPrice: state.idealGasPrice, gasLimit: 6000000 })
      .catch((e: Error) => {
        throw e;
      });
  }

  print("cyan", `Added ${state.toAdd} to the blocklist.`);
  print("h_green", `TX Hash: ${tx.hash}`);
}

async function removeFromBlocklist(state: State) {
  if (state.toRemove.length === 0) {
    print("yellow", "The are no addresses to remove from the blocklist");
    return;
  }

  print("yellow", "Removing addresses from the blocklist. Please wait...");

  let tx;
  if (state.toRemove.length === 1) {
    tx = await state.blocklistInstance
      .connect(state.activeAccount)
      .removeFromBlocklist(state.toRemove[0], { gasPrice: state.idealGasPrice, gasLimit: 6000000 })
      .catch((e: Error) => {
        throw e;
      });
  } else {
    // there are many addresses to remove
    tx = await state.blocklistInstance
      .connect(state.activeAccount)
      .batchRemoveFromBlocklist(state.toRemove, { gasPrice: state.idealGasPrice, gasLimit: 6000000 })
      .catch((e: Error) => {
        throw e;
      });
  }

  print("cyan", `Removed ${state.toRemove} from the blocklist.`);
  print("h_green", `TX Hash: ${tx.hash}`);
}

function treatCommonErrors(e: Error) {
  if (e.message.indexOf("getFullList") !== -1) {
    print(
      "h_red",
      "Error: cannot execute functions on the blocklist contract. Are you sure you have the right address in your .env variables?"
    );
  } else if (e.message.indexOf("Unsupported method") !== -1) {
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
    console.error({ e });
  }
}