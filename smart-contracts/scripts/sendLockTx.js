module.exports = async (cb) => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");
  const BigNumber = require("bignumber.js");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const BridgeBank = artifacts.require("BridgeBank");

  const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";

  console.log("Expected usage: \nBRIDGEBANK_ADDRESS='0x9e8bd20374898f61b4e5bd32b880b7fe198e44a1' truffle exec scripts/sendLockTx.js --network ropsten sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace eth 100\n");
  /*******************************************
   *** Constants
   ******************************************/
  // Lock transaction default params
  const DEFAULT_COSMOS_RECIPIENT = Web3.utils.utf8ToHex(
    "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
  );
  const DEFAULT_ETH_DENOM = "eth";
  const DEFAULT_AMOUNT = 10;

  // Config values
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";
  const NETWORK_MAINNET =
    process.argv[4] === "--network" && process.argv[5] === "mainnet";

  const DEFAULT_PARAMS =
    process.argv[4] === "--default" ||
    (NETWORK_ROPSTEN && process.argv[6] === "--default");
  const NUM_ARGS = process.argv.length - 4;

  /*******************************************
   *** Command line argument error checking
   ***
   *** truffle exec lacks support for dynamic command line arguments:
   *** https://github.com/trufflesuite/truffle/issues/889#issuecomment-522581580
   ******************************************/
  if ((NETWORK_MAINNET || NETWORK_ROPSTEN) && DEFAULT_PARAMS) {
    if (NUM_ARGS !== 3) {
      return console.error(
        "Error: custom parameters are invalid on --default."
      );
    }
  } else if (NETWORK_ROPSTEN || NETWORK_MAINNET) {
    if (NUM_ARGS !== 2 && NUM_ARGS !== 5) {
      return console.error(
        "Error: invalid number of parameters, please try again."
      );
    }
  } else if (DEFAULT_PARAMS) {
    if (NUM_ARGS !== 1) {
      return console.error(
        "Error: custom parameters are invalid on --default."
      );
    }
  } else {
    if (NUM_ARGS !== 3) {
      return console.error(
        "Error: must specify recipient address, token address, and amount."
      );
    }
  }

  /*******************************************
   *** Lock transaction parameters
   ******************************************/
  let cosmosRecipient = DEFAULT_COSMOS_RECIPIENT;
  let coinDenom = DEFAULT_ETH_DENOM;
  let amount = DEFAULT_AMOUNT;

  if (!DEFAULT_PARAMS) {
    if (NETWORK_ROPSTEN || NETWORK_MAINNET) {
      cosmosRecipient = Web3.utils.utf8ToHex(process.argv[6]);
      coinDenom = process.argv[7];
      amount = new BigNumber(process.argv[8]);
    } else {
      cosmosRecipient = Web3.utils.utf8ToHex(process.argv[4]);
      coinDenom = process.argv[5];
      amount = new BigNumber(process.argv[6]);
    }
  }

  // Convert default 'eth' coin denom into null address
  if (coinDenom == "eth") {
    coinDenom = NULL_ADDRESS;
  }

  try {
    /*******************************************
     *** Contract interaction
     ******************************************/
    // Get current accounts
    const accounts = await web3.eth.getAccounts();
    const bank = await BridgeBank.at(process.env.BRIDGEBANK_ADDRESS);

    // Send lock transaction
    console.log("Connected to contract, sending lock...");
    let str = (await web3.eth.getTransactionCount(accounts[0])).toString()
    let nonceVal = Number(str);
    console.log("starting nonce: ", nonceVal)
    let numIterations = Number(process.env.COUNT)
    for (let x = 0; x < 1; x++) {
      const promises = [];
      for (let i = 0; i < numIterations; i++) {
        txResultPromise = bank.lock(cosmosRecipient, coinDenom, amount, {
          from: accounts[0],
          value: coinDenom === NULL_ADDRESS ? amount : 0,
          gas: 200000, // 300,000 Gwei,
          nonce: nonceVal,
          gasPrice: 2110000000
        });
        promises.push(txResultPromise);
        nonceVal++;
        console.log(`Sent lock... ${i}`);
      }
      allPromise = Promise.all(promises);
      doneAllPromise = await allPromise;
      console.log("Done all:");
      console.log(doneAllPromise.map(tx => ({ lockBurnNonce: tx.logs[0].args._nonce.toNumber(), tx_id: tx.tx, status: tx.receipt.status,  })));
    }
  } catch (error) {
    console.error({ error });
  }
  return cb();
};
