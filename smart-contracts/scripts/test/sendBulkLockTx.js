const BN = require('bn.js');

module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.transactionYargOptions,
        'transactions': {
            type: "string",
            demandOption: true,
            description: 'json containing all the transactions to send.  Specify [{amount:, symbol:, sifchain_address:}].  The entries in --amount, --symbol, --sifchain_address are only used to estimate gas'
        },
        'lock_or_burn': {
            type: "string",
            default: "lock",
            description: 'set to either lock or burn'
        },
    });

    const bridgeBankContract = await contractUtilites.buildContract(this, argv, logging, "BridgeBank", argv.bridgebank_address);

    let cosmosRecipient = Web3.utils.utf8ToHex(argv.sifchain_address);
    let coinDenom = argv.symbol;
    let amount = argv.amount;

    let request = {
        from: argv.ethereum_address,
        value: coinDenom === sifchainUtilities.NULL_ADDRESS ? amount : 0,
        gas: argv.gas,
    };

    if (request.gas === 'estimate') {
        const estimate = await bridgeBankContract.lock.estimateGas(cosmosRecipient, coinDenom, amount, {
            ...request,
            gas: 6000000,
        });
        // increase by 10%
        request.gas = new BN(estimate, 10).mul(new BN(11)).div(new BN(10));
    }

    let transactions = JSON.parse(argv.transactions);
    const actions = [];
    for (const t of transactions) {
        logging.info(`calling bridgeBankContract.lock for ${t.sifchain_address}, amount is |${t.amount}|`);
        let lockResult;
        const amount = new BN(t.amount);
        try {
            if (argv.lock_or_burn === "lock")
                lockResult = bridgeBankContract.lock(Web3.utils.utf8ToHex(t.sifchain_address), t.symbol, amount, request);
            else
                lockResult = bridgeBankContract.burn(Web3.utils.utf8ToHex(t.sifchain_address), t.symbol, amount, request);
        } catch (error){
            logging.info(`goterror: ${error}`);
        }
        console.debug(`lockResult is ${lockResult}`);
        actions.push({lockResult, t});
    }
    const results = [];
    const blockCounts = {};
    for (const a of actions) {
        const result = await a["lockResult"];
        logging.info(`bridgeBankContract.lock result for ${a.t.sifchain_address}: ${JSON.stringify(result)}`);
        results.push(result);
        const blockNumber = result["receipt"]["blockNumber"];
        const existingBlockCount = blockCounts[blockNumber] || 0;
        blockCounts[blockNumber] = existingBlockCount + 1;
    }

    logging.info("all locks submitted");

    const web3 = contractUtilites.buildWeb3(this, argv, logging);

    const blockNumber = await web3.eth.getBlockNumber();

    console.log(JSON.stringify({
        blockNumber: blockNumber,
        blockCounts
        // results: results,
    }, undefined, 0))

    return cb();
};
