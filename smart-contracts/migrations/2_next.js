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
  let owner = process.env.OWNER;
  let pauser = process.env.PAUSER;
  let initialValidators = process.env.INITIAL_VALIDATOR_ADDRESSES.split(",");
  let initialPowers = process.env.INITIAL_VALIDATOR_POWERS.split(",");
  const tokenAmount = web3.utils.toWei("120000000");

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
      setTxSpecifications(6721975, operator, deployer)
    );

    console.log("cosmosBridge address: ", cosmosBridge.address)

    // 2. Deploy BridgeBank contract:
    //    Gas used:        4,823,348 Gwei
    //    Total cost:    0.09646696 Ether
    const bridgeBank = await deployProxy(
      BridgeBank,
      [
        operator,
        CosmosBridge.address,
        owner,
        pauser
      ],
      setTxSpecifications(6721975, operator, deployer)
    );
    console.log("bridgeBank address: ", bridgeBank.address)

    // 3. Deploy BridgeRegistry contract:
    //    Gas used:          363,370 Gwei
    //    Total cost:     0.0072674 Ether
    await deployProxy(
      BridgeRegistry,
      [
        CosmosBridge.address,
        BridgeBank.address
      ],
      setTxSpecifications(6721975, operator, deployer)
    );

    await cosmosBridge.setBridgeBank(bridgeBank.address, 
      setTxSpecifications(600000, operator)
    );

    if (network === 'mainnet') {
      return console.log("Network is mainnet, not going to deploy token");
    }

    const erowan = await deployer.deploy(eRowan, "erowan", setTxSpecifications(4612388, operator));

    await erowan.addMinter(BridgeBank.address, setTxSpecifications(4612388, operator));

    await bridgeBank.addExistingBridgeToken(erowan.address, setTxSpecifications(4612388, operator));

    const tokenAddress = "0x0000000000000000000000000000000000000000";

    // allow 10 eth to be sent at once
    await bridgeBank.updateTokenLockBurnLimit(tokenAddress, '10000000000000000000', setTxSpecifications(4612388, operator));
    await bridgeBank.updateTokenLockBurnLimit(erowan.address, '10000000000000000000', setTxSpecifications(4612388, operator));
    await erowan.approve(bridgeBank.address, '10000000000000000000', setTxSpecifications(4612388, operator));

    console.log("erowan token address: ", erowan.address);

    const bnAmount = web3.utils.toWei("100", "ether");

    await erowan.mint(operator, bnAmount, setTxSpecifications(4612388, operator));

    if (network === "develop") {
      await erowan.mint(accounts[1], bnAmount, setTxSpecifications(4612388, operator));
    }

    return;
  });
};