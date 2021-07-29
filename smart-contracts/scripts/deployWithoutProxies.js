const c = require('./constants').get();
const state = {
    deployedCount: 0,
    bridgeBank: null,
    cosmosBridge: null,
};
const EXPECTED_DEPLOYMENTS = 2;

module.exports = async () => {
    try {
        // Deploy BridgeBank
        const bbHash = await deploy(
            c.BridgeBankContract.abi,
            c.BridgeBankContract.bytecode,
            [],
            onBridgeBankDeployed
        );
        console.log(`BridgeBank deploy TxHash: ${bbHash}`);

        // Deploy CosmosBridge
        const cbHash = await deploy(
            c.CosmosBridgeContract.abi,
            c.CosmosBridgeContract.bytecode,
            [],
            onCosmosBridgeDeployed
        );
        console.log(`CosmosBridge deploy TxHash: ${cbHash}`);
    } catch (error) {
        console.error({error})
    }
}

async function deploy(abi, bytecode, arguments, callback) {
    const contract = new c.web3.eth.Contract(abi);
    const hash = await signAndSendTransaction(
        contract.deploy({ data: bytecode, arguments }),
        callback
    );

    return hash;
}

async function signAndSendTransaction(transaction, callback) {
    const gas = await transaction.estimateGas({ from: c.env.owner });
    const gasPrice = c.mainnetGasPrice;
    const nonce = await c.web3.eth.getTransactionCount(
      c.env.owner,
      "pending"
    );

    const options = {
      data: transaction.encodeABI(),
      gas,
      gasPrice,
      nonce,
    };

    const signedTransaction = await c.web3.eth.accounts.signTransaction(
      options,
      c.env.ethereumPrivateKey
    );

    return new Promise((resolve, reject) => {
        c.web3.eth.sendSignedTransaction(
          signedTransaction.rawTransaction,
          (err, data) => {
            if (err) {
              return reject(err);
            }
          }
        )
        .off("sending", () => {})
        .off("sent", () => {})
        .off("confirmation", () => {})
        .once("transactionHash", (hash) => {
          resolve(hash);
        })
        .once("receipt", (receipt) => {
          if (callback) {
            callback(receipt);
          }
        })
        .catch((e) => {
          console.log(`FAILED: ${e.message}`);
        });
    });
}

async function onBridgeBankDeployed(receipt) {
    console.log(`BridgeBank deployed at ${receipt.contractAddress}`);

    state.bridgeBank = await c.BridgeBankContract.at(receipt.contractAddress);

    state.deployedCount++;
    if(state.deployedCount === EXPECTED_DEPLOYMENTS) {
        initializeContracts();
    }
}

async function onCosmosBridgeDeployed(receipt) {
    console.log(`CosmosBridge deployed at ${receipt.contractAddress}`);

    state.cosmosBridge = await c.CosmosBridgeContract.at(receipt.contractAddress);

    state.deployedCount++;
    if(state.deployedCount === EXPECTED_DEPLOYMENTS) {
        initializeContracts();
    }
}

async function initializeContracts() {
    try {
        const bbInitReceipt = await state.bridgeBank.methods['initialize(address,address,address,address)']
        (
            c.env.operator,
            state.cosmosBridge.address,
            c.env.owner,
            c.env.pauser,
            { from: c.env.operator }
        );

        console.log(`\nBridgeBank initialization receipt: ${JSON.stringify(bbInitReceipt, null, 1)}`);
            
        const cbInitReceipt = await state.cosmosBridge.initialize(
            c.env.operator,
            c.env.consensusThreshold,
            c.env.initialValidatorAddresses,
            c.env.initialValidatorPowers,
            { from: c.env.operator }
        );

        console.log(`\nCosmosBridge initialization receipt: ${JSON.stringify(cbInitReceipt, null, 1)}`);

        console.log('\n\nSUCCESS: Contracts initialized!');

        process.exit(0);
    } catch (error) {
        console.error({error})
    }
}