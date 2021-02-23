module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const bridgeContract = truffleContract(
    require("../build/contracts/BridgeBank.json")
  );
  const tokenContract = truffleContract(
    require("../build/contracts/BridgeToken.json")
  );

  console.log("Expected usage: \n truffle exec scripts/sendApproveTx.js --network ropsten 10 0xdDA6327139485221633A1FcD65f4aC932E60A2e1");

  /*******************************************
   *** Constants
   ******************************************/
  // Config values
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";
  const DEFAULT_PARAMS =
    process.argv[4] === "--default" ||
    (NETWORK_ROPSTEN && process.argv[6] === "--default");
  const NUM_ARGS = process.argv.length - 4;

  // Default transaction parameters
  const DEFAULT_TOKEN_AMOUNT = 100;

  /*******************************************
   *** Command line argument error checking
   ***
   *** truffle exec lacks support for dynamic command line arguments:
   *** https://github.com/trufflesuite/truffle/issues/889#issuecomment-522581580
   ******************************************/
  if (NETWORK_ROPSTEN) {
    if (NUM_ARGS !== 3 && NUM_ARGS !== 4) {
      return console.error(
        "Error: Must specify token amount if using the Ropsten network."
      );
    }
  } else {
    if (NUM_ARGS !== 1 && NUM_ARGS !== 2) {
      return console.error("Error: Must specify token amount or --default.");
    }
  }

  /*******************************************
   *** Approve transaction parameters
   ******************************************/
  let tokenAmount;

  if (NETWORK_ROPSTEN) {
    tokenAmount = process.argv[6];
  } else {
    if (!DEFAULT_PARAMS) {
      tokenAmount = process.argv[4];
    } else {
      tokenAmount = DEFAULT_TOKEN_AMOUNT;
    }
  }


  /*******************************************
   *** Approve transaction parameters
   ******************************************/
  let tokenAddress;

  if (NETWORK_ROPSTEN) {
    tokenAddress = process.argv[7];
  } else {
    if (!DEFAULT_PARAMS) {
      tokenAddress = process.argv[5];
    } else {
      tokenAddress = false;
    }
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

  bridgeContract.setProvider(web3.currentProvider);
  tokenContract.setProvider(web3.currentProvider);
  try {
    /*******************************************
     *** Contract interaction
    ******************************************/
    // Get current accounts
    const accounts = await web3.eth.getAccounts();

    const bridgeContractAddress = await bridgeContract
      .deployed()
      .then(function(instance) {
        return instance.address;
      });
    
    let instance
    if (tokenAddress) {
      instance = await tokenContract.at(tokenAddress)
    } else {
      instance = await tokenContract.deployed()
    }

    // Send lock transaction
    const { logs } = await instance.approve(bridgeContractAddress, tokenAmount, {
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei
    });

    // Get event logs
    const event = logs.find(e => e.event === "Approval");

    // Parse event fields
    const approvalEvent = {
      owner: event.args.owner,
      spender: event.args.spender,
      value: Number(event.args.value)
    };

    console.log(approvalEvent);
  } catch (error) {
    console.error({error})
  }
  return;
};
