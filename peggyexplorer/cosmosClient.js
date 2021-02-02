const fetch = require('node-fetch');
const WebSocket = require('ws');

const SIFNODE_RPC_URL = 'http://0.0.0.0:26657/';

const ws = new WebSocket(SIFNODE_RPC_URL + 'websocket');

const allCosmosTxs = [];
exports.allCosmosTxs = allCosmosTxs;

let syncHeight = 0;
let syncPageNumber = 1;

exports.kickoff = async function () {

  ws.on('open', function open() {
    console.log("Socket open")
    ws.send(JSON.stringify({ "jsonrpc": "2.0", "method": "status", "params": [], "id": 'sync' }));
  });

  ws.on('message', function incoming(message) {
    const messageJSON = JSON.parse(message);
    const id = messageJSON && messageJSON.id;
    switch (id) {
      case 'sync':
        const sync = messageJSON.result.sync_info;
        processSync(sync);
        break;
      case 'tx_sync':
        if (messageJSON.error && messageJSON.error.data.includes('page should be')) {
          console.log("Sync complete");
          startSubscribe();
        } else {
          const txSync = messageJSON.result;
          processTxSync(txSync);
        }
        break;
      case 'tx_subscription':
        const txSub = messageJSON.result;
        if (messageJSON.result && messageJSON.result.data && messageJSON.result.data.type === 'tendermint/event/Tx') {
          console.log(messageJSON.result);
          const hash = messageJSON.result.events['tx.hash'][0];
          processTxSub(messageJSON.result.data.value, hash);
        }
        break;
      default:
        console.log("Unknown message received: ");
        console.log(messageJSON);
        break;
    }
  });
}

const startSubscribe = _ => {
  console.log(`Subscribing to txs after block ${syncHeight}`);
  ws.send(JSON.stringify({
    "jsonrpc": "2.0", "method": "subscribe",
    "params": [
      `tm.event = 'Tx' AND tx.height > ${syncHeight}`
    ],
    "id": 'tx_subscription'
  }));
}

const processTxSync = txSync => {
  console.log("New txs received during sync");
  const total = txSync.total_count;
  console.log(`Total txs to sync: ${total}`);
  const totalPages = Math.ceil(total / 100);
  console.log(`Page ${syncPageNumber} of ${totalPages}`)
  const txs = txSync.txs;
  processNewTxs(txs);
  syncPageNumber++;
  console.log(`Processed, querying page ${syncPageNumber}`);
  ws.send(JSON.stringify({
    "jsonrpc": "2.0",
    "method": "tx_search",
    "params": [
      `tx.height<=${syncHeight}`,
      false,
      '' + syncPageNumber,
      '100',
      "asc"
    ]
    , "id": `tx_sync`
  }));
}

const processTxSub = (txSub, hash) => {
  console.log("New tx received from subscription");
  const tx = Object.assign({}, txSub.TxResult, { tx_result: txSub.TxResult.result, result: undefined, hash });
  allCosmosTxs.push(tx);
}

const processSync = sync => {
  syncHeight = sync.latest_block_height;
  console.log(`New sync, processing with height: ${syncHeight}`);
  console.log('Querying all txs below sync point.');
  console.log(`Querying page ${syncPageNumber}`);
  ws.send(JSON.stringify({
    "jsonrpc": "2.0",
    "method": "tx_search",
    "params": [
      `tx.height<=${syncHeight}`,
      false,
      '' + syncPageNumber,
      '100',
      "asc"
    ]
    , "id": `tx_sync`
  }));
}

const processNewTxs = txs => {
  console.log(`Processing ${txs.length} new txs`);
  allCosmosTxs.push(...txs);
}

const processNewBlock = block => {
  console.log(`New cosmos block received: Block ${block.header.height}`);
  console.log(`Pulling all transactions before ${block.header.height}`);
}
