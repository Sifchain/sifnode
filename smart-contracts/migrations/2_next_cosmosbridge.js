require("dotenv").config();

const { deployProxy } = require('@openzeppelin/truffle-upgrades');

const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");
const BridgeRegistry = artifacts.require("BridgeRegistry");
const eRowan = artifacts.require("BridgeToken");

module.exports = function(deployer, network, accounts) {
  /*******************************************
   *** Input validation of contract params
   ******************************************/
  if (process.env.CONSENSUS_THRESHOLD.length === 0) {
    return console.error(
      "Must provide consensus threshold as environment variable."
    );
  }
  if (process.env.OPERATOR.length === 0) {
    return console.error(
      "Must provide operator address as environment variable."
    );
  }
  if (process.env.OWNER.length === 0) {
    return console.error(
      "Must provide owner address as environment variable."
    );
  }

  let consensusThreshold = process.env.CONSENSUS_THRESHOLD;
  let operator = process.env.OPERATOR;
  let initialValidators = process.env.INITIAL_VALIDATOR_ADDRESSES.split(",");
  let initialPowers = process.env.INITIAL_VALIDATOR_POWERS.split(",");

  if (!initialPowers.length || !initialValidators.length) {
    return console.error(
      "Must provide validator and power."
    );
  }
  if (initialPowers.length !== initialValidators.length) {
    return console.error(
      "Each initial validator must have a corresponding power specified."
    );
  }

  /*******************************************************
   *** Contract deployment summary
   ***
   *** Total deployments:       7 (includes Migrations.sol)
   *** Gas price (default):                       20.0 Gwei
   *** Final cost:                         0.25369878 Ether
   *******************************************************/
  deployer.then(async () => {

    function setTxSpecifications(gasAmount, from, deployObject) {
      const txObj = {
        gas: gasAmount,
        from: from,
        unsafeAllowCustomTypes: true,
        deployer: deployObject
      }

      if (process.env.MAINNET_GAS_PRICE) {
        txObj.gasPrice = process.env.MAINNET_GAS_PRICE
      }

      return txObj;
    }

    // 1. Deploy CosmosBridge contract:
    //    Gas used:       2,649,300 Gwei
    //    Total cost:     0.052986 Ether
    const cosmosBridge = await deployProxy(
      CosmosBridge,
      [
        operator,
        consensusThreshold,
        initialValidators,
        initialPowers
      ],
      setTxSpecifications(6721975, accounts[0], deployer)
    );

    console.log("cosmosBridge address: ", cosmosBridge.address)

    return;
  });
};