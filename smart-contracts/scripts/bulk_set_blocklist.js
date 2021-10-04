require("dotenv").config();
const fs = require("fs");
const axios = require("axios");
const { ethers } = require("hardhat");

const { print, generateTodayFilename } = require("./helpers/utils");
const parser = require("./helpers/ofacParser");

// Defaults to the Ropsten address
const BLOCKLIST_ADDRESS =
  process.env.BLOCKLIST_ADDRESS || "0xbB4fa6cC28f18Ae005998a5336dbA3bC49e3dc57";

const state = {
  ofac: [],
  evm: [],
  toAdd: [],
  toRemove: [],
  blocklistInstance: null,
  owner: null,
};

async function main() {
  state.ofac = await parser.getList();
  state.evm = await fetchEvmBlocklist();
  calculateDiff();
}

async function fetchEvmBlocklist() {
  const blocklistFactory = await ethers.getContractFactory("Blocklist");
  state.blocklistInstance = await blocklistFactory.attach(BLOCKLIST_ADDRESS);
  const fullList = await state.blocklistInstance.getFullList();

  return fullList;
}

function calculateDiff() {
  // addresses that must be added don't exist on evm, but exist on ofac
  state.toAdd = state.ofac.filter((address) => !state.evm.includes(address));

  // addresses that must be removed exist on evm, but don't exist on ofac
  state.toRemove = state.evm.filter((address) => !state.ofac.includes(address));
}

// TODO: get the owner account
async function addToBlocklist() {
  if (state.toAdd.length === 0) {
    print("green", "The are no new addresses to add to the blocklist.");
    return;
  }

  // TODO: use Hardhat or Truffle!
  // 1) Would be nice to be able to impersonate an account here to test with a fork
  const accounts = await web3.eth.getAccounts();

  if (state.toAdd.length === 1) {
    await state.blocklistInstance.addToBlocklist(state.toAdd[0], {
      from: TODO,
    });
  } else {
    // there are many addresses to add
    await state.blocklistInstance.batchAddToBlocklist(state.toAdd, {
      from: TODO,
    });
  }
}

async function removeFromBlocklist() {}

main()
  .catch((error) => {
    console.error({ error });
  })
  .finally(() => process.exit(0));
