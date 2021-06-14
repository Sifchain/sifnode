
module.exports = async (cb) => {
    const expectedUsage = () => {console.log("Expected usage:\nBRIDGEBANK_ADDRESS='insert bridgebank address' COSMOS_BRIDGE_ADDRESS='insert cosmosbridge address' truffle exec scripts/setBridgeBank.js --network mainnet\n")}
    try {
      /*******************************************
       *** Set up
      ******************************************/
      const ethers = require("ethers");
  
      const CosmosBridgeContract = artifacts.require("CosmosBridge");
  
      const bridgeBankContractAddress = process.env.BRIDGEBANK_ADDRESS;
  
      if (!bridgeBankContractAddress || bridgeBankContractAddress.length !== 42) {
        throw new Error("error, no bridgebank address")
      }
  
      if (!process.env.COSMOS_BRIDGE_ADDRESS || process.env.COSMOS_BRIDGE_ADDRESS.length !== 42) {
        throw new Error("error, no cosmos bridge address")
      }
      /*******************************************
       *** Constants
      ******************************************/
      // Config values
      const NETWORK_ROPSTEN =
        process.argv[4] === "--network" && process.argv[5] === "ropsten";
      const NETWORK_MAINNET =
        process.argv[4] === "--network" && process.argv[5] === "mainnet";
  
      /*******************************************
       *** Ether provider
      *** Set contract provider based on --network flag
      ******************************************/
      let provider;
      if (NETWORK_ROPSTEN) {
        provider = new ethers.providers.InfuraProvider("ropsten", process.env.INFURA_PROJECT_ID);
          
      } else if (NETWORK_MAINNET) {
        provider = new ethers.providers.InfuraProvider(null, process.env.INFURA_PROJECT_ID);
          
      } else {
        provider = new ethers.providers.JsonRpcProvider(process.env.LOCAL_PROVIDER);
      }
        
      try {
        /*******************************************
         *** Contract interaction
        ******************************************/
        // Get current accounts
        const accounts = await provider.listAccounts();
        let cosmosBridgeContract = await CosmosBridgeContract.at(process.env.COSMOS_BRIDGE_ADDRESS)
        // Set BridgeBank
        console.log("Loaded accounts and contract, setting bridgebank...");
  
        await cosmosBridgeContract.setBridgeBank(bridgeBankContractAddress, {
          from: accounts[0],
          value: 0,
          gas: 300000 // 300,000 Gwei
        });
  
        console.log("CosmosBridge's BridgeBank address set");
  
        cb();
      } catch (error) {
        expectedUsage();
        console.error({error})
        cb();
      }
    } catch (error) {
      expectedUsage();
      console.error({ error })
      return cb()
    }
  };
  
