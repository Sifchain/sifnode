module.exports = async (cb) => {
    try {


    const CosmosBridge = artifacts.require("CosmosBridge");

    /*******************************************
     *** Set up
     ******************************************/
    require("dotenv").config();
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");

    // Contract abstraction

    /*******************************************
     *** Constants
     ******************************************/
    const NETWORK_ROPSTEN =
      process.argv[4] === "--network" && process.argv[5] === "ropsten";

    /*******************************************
     *** Web3 provider
     *** Set contract provider based on --network flag
     ******************************************/
    let provider;
    if (NETWORK_ROPSTEN) {
      provider = new HDWalletProvider(
        process.env.ETHEREUM_PRIVATE_KEY,
        process.env['WEB3_PROVIDER']
      );
    } else {
      provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
    }
    const cosmosBridge = await CosmosBridge.at("0x6A85ABc7a7400520e7E454ef4ABbB8AEE8b156bE");

    const web3 = new Web3(provider);

    const CLAIM_TYPE_BURN = 1;
    const symbol = "ETH";
    const cosmosSender = "0x736966316e78363530733871397732386632673374397a74787967343875676c64707475777a70616365";
    const cosmosSenderSequence = 1;
    const amount = 0;
    const ethereumReceiver = "0xf17f52151EbEF6C7334FAD080c5704D77216b732";

    let estimatedGas = await cosmosBridge.newProphecyClaim.estimateGas(
        CLAIM_TYPE_BURN,
        cosmosSender,
        cosmosSenderSequence,
        ethereumReceiver,
        symbol,
        amount,
        {
            from: "0x1Aa97F2463A78364F6D3Da90EEb99F8CDb9392f4"
        }
    );
    console.log("Estimated gas cost: ", estimatedGas);

    estimatedGas = await cosmosBridge.newProphecyClaim.estimateGas(
        CLAIM_TYPE_BURN,
        cosmosSender,
        cosmosSenderSequence,
        ethereumReceiver,
        "eth",
        amount,
        {
            from: "0x1Aa97F2463A78364F6D3Da90EEb99F8CDb9392f4"
        }
    );
    console.log("Estimated gas cost: ", estimatedGas);
    cb();
    } catch (error) {
        console.error("Error: ", error)
        cb();
    }
  };
