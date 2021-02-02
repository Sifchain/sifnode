const { _ } = require('lodash');
const BRIDGE_REGISTRY_CONTRACT_ABI = require('../smart-contracts/build/contracts/BridgeRegistry').abi;
const BANK_STORAGE_CONTRACT_ABI = require('../smart-contracts/build/contracts/BankStorage').abi;
const BRIDGE_BANK_CONTRACT_ABI = require('../smart-contracts/build/contracts/BridgeBank').abi;
const Web3 = require('web3');
const Contract = require('web3-eth-contract');

const ETHEREUM_PROVIDER_URL = 'wss://ropsten.infura.io/ws/v3/f1f4e06cebc8462b846a67328cb67e90';
const SIFNODE_RPC_URL = 'http://rpc-sandpit.sifchain.finance:26657';
const BRIDGE_BANK_CONTRACT_ADDRESS = '0x979F0880de42A7aE510829f13E66307aBb957f13';
const STARTING_ETHEREUM_BLOCK = 0;

const web3 = new Web3(Web3.givenProvider || ETHEREUM_PROVIDER_URL);
Contract.setProvider(ETHEREUM_PROVIDER_URL);

const bridgeBankContract = new web3.eth.Contract(BRIDGE_BANK_CONTRACT_ABI, BRIDGE_BANK_CONTRACT_ADDRESS);

const allEthereumEvents = [];
const allEthereumBlocksChecked = [];

start();

async function start() {
  console.log(`Loaded bridge bank contract at ${bridgeBankContract._address}`);
  startWatchingBlocks();
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
  allEthereumEvents.push(allEvents);
}

pushEthereumBlock = block => {
  allEthereumBlocksChecked.push(block);
}
