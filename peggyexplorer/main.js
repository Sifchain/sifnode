const { _ } = require('lodash');
const { LcdClient } = require('@cosmjs/launchpad');
const cosmosClient = require('./cosmosClient');

const BRIDGE_REGISTRY_CONTRACT_ABI = require('../smart-contracts/build/contracts/BridgeRegistry').abi;
const BANK_STORAGE_CONTRACT_ABI = require('../smart-contracts/build/contracts/BankStorage').abi;
const BRIDGE_BANK_CONTRACT_ABI = require('../smart-contracts/build/contracts/BridgeBank').abi;
const Web3 = require('web3');
const Contract = require('web3-eth-contract');

const ROPSTEN_ETHEREUM_PROVIDER_URL = 'wss://ropsten.infura.io/ws/v3/f1f4e06cebc8462b846a67328cb67e90';
const LOCAL_ETHEREUM_PROVIDER_URL = 'ws://0.0.0.0:7545';
const ETHEREUM_PROVIDER_URL = LOCAL_ETHEREUM_PROVIDER_URL;

const BRIDGE_BANK_CONTRACT_ADDRESS_ROPSTEN = '0x979F0880de42A7aE510829f13E66307aBb957f13';
const BRIDGE_BANK_CONTRACT_ADDRESS_LOCAL = '0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4';
const BRIDGE_BANK_CONTRACT_ADDRESS = BRIDGE_BANK_CONTRACT_ADDRESS_LOCAL;
const STARTING_ETHEREUM_BLOCK = 0;

const web3 = new Web3(Web3.givenProvider || ETHEREUM_PROVIDER_URL);
Contract.setProvider(ETHEREUM_PROVIDER_URL);

const bridgeBankContract = new web3.eth.Contract(BRIDGE_BANK_CONTRACT_ABI, BRIDGE_BANK_CONTRACT_ADDRESS);

const allEthereumEvents = [];
const allEthereumBlocksChecked = [];
const allCosmosTxs = cosmosClient.allCosmosTxs;

const express = require('express')
var cors = require('cors')

const app = express()
app.use(cors())

const port = 5000

app.get('/dump', (req, res) => {
  res.json({
    allEthereumEvents, allEthereumBlocksChecked, allCosmosTxs
  })
})

app.listen(port, () => {
  console.log(`Listening at http://localhost:${port}`)
})

start();

async function start() {
  console.log(`Loaded bridge bank contract at ${bridgeBankContract._address}`);
  startWatchingBlocks();
  startWatchingCosmos();
}

async function startWatchingBlocks(blockNumber) {
  web3.eth.subscribe('newBlockHeaders', (error, result) => {
    if (error) {
      console.log(error);
    } else {
      console.log(`Discovered new ethereum block: ${result.number}`);
      if (allEthereumBlocksChecked.length === 0) {
        console.log("First block found, populating from beginning...")
        populatePastEventsFrom(STARTING_ETHEREUM_BLOCK, result.number);
      } else {
        populatePastEventsFrom(result.number, result.number);
      }
    }
  })
}

async function populatePastEventsFrom(startingBlock, endingBlock) {
  allEvents = await bridgeBankContract.getPastEvents("allEvents", {
    fromBlock: startingBlock,
    toBlock: endingBlock
  });

  pushEthereumEvents(allEvents);
  for (let i = startingBlock; i <= endingBlock; i++) {
    pushEthereumBlock(i);
  }

  const pastLockEvents = _.filter(allEvents, { 'event': 'LogLock' });
  const pastUnlockEvents = _.filter(allEvents, { 'event': 'LogUnlock' });
  const pastBurnEvents = _.filter(allEvents, { 'event': 'LogBurn' });
  const pastMintEvents = _.filter(allEvents, { 'event': 'LogBridgeTokenMint' });

  const blocksPopulated = startingBlock === endingBlock ? startingBlock : `${startingBlock} to ${endingBlock}`;
  console.log(`Populated ethereum events from block ${blocksPopulated}:
    locks: ${pastLockEvents.length},
    unlocks: ${pastUnlockEvents.length},
    burns: ${pastBurnEvents.length},
    mints: ${pastMintEvents.length}`
  );
}

pushEthereumEvents = events => {
  allEthereumEvents.push(...allEvents);
}

pushEthereumBlock = block => {
  allEthereumBlocksChecked.push(block);
}

async function startWatchingCosmos() {
  console.log("Starting to watch cosmos:");
  const info = await cosmosClient.kickoff();
  console.log({ info });
}
