module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const contract = truffleContract(
    require("../build/contracts/BridgeBank.json")
  );

  const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";

  // Config values
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";
  const NUM_ARGS = process.argv.length - 4;

  /*******************************************
   *** Command line argument error checking
   ***
   *** truffle exec lacks support for dynamic command line arguments:
   *** https://github.com/trufflesuite/truffle/issues/889#issuecomment-522581580
   ******************************************/
  if (NETWORK_ROPSTEN) {
    if (NUM_ARGS !== 2 && NUM_ARGS !== 5) {
      return console.error(
        "Error: invalid number of parameters, please try again."
      );
    }
  } else {
    if (NUM_ARGS !== 2) {
      return console.error(
        "Error: must specify token address, and new value."
      );
    }
  }

  console.log("Expected usage: \n truffle exec scripts/sendUpdateWhiteList.js --network ropsten 0xdDA6327139485221633A1FcD65f4aC932E60A2e1 true");

  /*******************************************
   *** Lock transaction parameters
   ******************************************/
  let coinDenom = NULL_ADDRESS;
  let inList = false;

  if (NETWORK_ROPSTEN) {
    coinDenom = process.argv[6];
    inList = process.argv[7];
  } else {
    coinDenom = process.argv[4];
    inList = process.argv[5];
  }

  // Convert default 'eth' coin denom into null address
  if (coinDenom == "eth") {
    coinDenom = NULL_ADDRESS;
  }

  /*******************************************
   *** Web3 provider
   *** Set contract provider based on --network flag
   ******************************************/
  let provider;
  if (NETWORK_ROPSTEN) {
    provider = new HDWalletProvider(
      process.env.ETHEREUM_PRIVATE_KEY,
      "https://ropsten.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
    );
  } else {
    provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
  }

  const web3 = new Web3(provider);
  contract.setProvider(web3.currentProvider);
  try {
    /*******************************************
     *** Contract interaction
     ******************************************/
    // Get current accounts
    const accounts = await web3.eth.getAccounts();

    // Send update whitelist transation
    console.log("Connecting to contract....");
    const {
      logs
    } = await contract.deployed().then(function (instance) {
      console.log("Connected to contract, sending lock...");
      return instance.updateEthWhiteList(coinDenom, inList, {
        from: accounts[0],
        gas: 300000 // 300,000 Gwei
      });
    });

    console.log("Sent update white list...");

    // Get event logs
    const event = logs.find(e => e.event === "LogWhiteListUpdate");

    // Parse event fields
    const whiteListUpdateEvent = {
      token: event.args._token,
      in_list: event.args._value,
    };

    console.log(whiteListUpdateEvent);
  } catch (error) {
    console.error({
      error
    });
  }
  return;
};
