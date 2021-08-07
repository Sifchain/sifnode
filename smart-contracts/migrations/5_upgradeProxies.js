/**
 * Upgrades BridgeBank and CosmosBridge
 * 
 * Expected usage (option 1): yarn peggy:upgradeProxies 'sifchain-testnet' --network ropsten --skip-dry-run
 * Expected usage (option 2): truffle migrate -f 5 --to 5 sifchain-testnet --network ropsten --skip-dry-run
 * Change 'sifchain-testnet' to the desired deployment
 */

const { upgradeProxy, prepareUpgrade } = require('@openzeppelin/truffle-upgrades');
const truffleContract = require("truffle-contract");
const fs = require('fs');

const c = require('../scripts/constants').get();

const deploymentsPath = './deployments';

const state = {
  bridgeBank: {
    current: {
      contract: null,
      parsedArtifacts: null,
      address: null,
    },
    new: null,
  },
  cosmosBridge: {
    current: {
      contract: null,
      parsedArtifacts: null,
      address: null,
    },
    new: null,
  }
};

function setup() {
  const folder = c.argv[7]; // something like 'sifchain-testnet'
  const basePathForDeployedContracts = `${deploymentsPath}/${folder}/`;
  const basePathForTruffleArtifacts = `../${deploymentsPath}/${folder}/`;

  console.log(`UPGRADING :: deployment: ${folder} | chainId: ${c.chainId} ...`);

  state.bridgeBank.current.contract = truffleContract(require(`${basePathForTruffleArtifacts}BridgeBank`));
  const bbConfig = fs.readFileSync(`${basePathForDeployedContracts}BridgeBank.json`, 'utf8');
  state.bridgeBank.current.parsedArtifacts = JSON.parse(bbConfig);
  state.bridgeBank.current.address = state.bridgeBank.current.parsedArtifacts.networks[c.chainId].address;

  state.cosmosBridge.current.contract = truffleContract(require(`${basePathForTruffleArtifacts}CosmosBridge`));
  const cbConfig = fs.readFileSync(`${basePathForDeployedContracts}CosmosBridge.json`, 'utf8');
  state.cosmosBridge.current.parsedArtifacts = JSON.parse(cbConfig);
  state.cosmosBridge.current.address = state.cosmosBridge.current.parsedArtifacts.networks[c.chainId].address;

  console.log(`-> Current BridgeBank at: ${state.bridgeBank.current.address}`);
  console.log(`-> Current CosmosBridge at: ${state.cosmosBridge.current.address}`);
}

module.exports = async function (deployer, network, accounts) {
  deployer.then(async () => {
    function setTxSpecifications(gasAmount, from, deployObject) {
      const txObj = {
        gas: gasAmount,
        from: from,
        unsafeAllowCustomTypes: true,
        deployer: deployObject,
      }

      if (c.env.mainnetGasPrice) {
        txObj.gasPrice = c.env.mainnetGasPrice;
      }

      return txObj;
    }

    try {
      setup();

      state.bridgeBank.current.contract.setProvider(c.web3.currentProvider);
      state.cosmosBridge.current.contract.setProvider(c.web3.currentProvider);

      const currentBbInstance = await state.bridgeBank.current.contract.at(state.bridgeBank.current.address);
      const currentCbInstance = await state.cosmosBridge.current.contract.at(state.cosmosBridge.current.address);

      state.bridgeBank.new = await prepareUpgrade(currentBbInstance.address, c.BridgeBankContract, setTxSpecifications(3000000, accounts[0], deployer));
      state.cosmosBridge.new = await prepareUpgrade(currentCbInstance.address, c.CosmosBridgeContract, setTxSpecifications(3000000, accounts[0], deployer));

      console.log("--> Prepared BridgeBank:", state.bridgeBank.new.address);
      console.log("--> Prepared CosmosBridge:", state.cosmosBridge.new.address);
    } catch (e) {
      console.log({ e });
    }
  });
}