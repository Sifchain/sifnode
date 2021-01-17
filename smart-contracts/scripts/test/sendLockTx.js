module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const blockNumber = await web3.eth.getBlockNumber();
    logging.debug(`blockNumber: ${blockNumber}`);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.transactionYargOptions
    });

    logging.info(`starting sendLockTx, arguments are ${JSON.stringify(argv)}`);

    const bridgeBankContract = contractUtilites.buildContract(this, argv, "BridgeBank", argv.bridgebank_address);

    const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";

    /*******************************************
     *** Lock transaction parameters
     ******************************************/
    let cosmosRecipient = Web3.utils.utf8ToHex(argv.sifchain_address);
    let coinDenom = argv.symbol;
    let amount = argv.amount;

    // Convert default 'eth' coin denom into null address
    if (coinDenom === "eth") {
        coinDenom = NULL_ADDRESS;
    }

    try {
        const {logs} = await bridgeBankContract.then(function (instance) {
            let request = {
                from: argv.ethereum_address,
                value: coinDenom === NULL_ADDRESS ? amount : 0,
                gas: argv.gas
            };
            logging.debug(`Connected to contract, sending lock, request is ${JSON.stringify(request)}`)
            return instance.lock(cosmosRecipient, coinDenom, amount, request);
        });

        logging.debug(`logs are ${JSON.stringify(logs)}`);

        // Get event logs
        const event = logs.find(e => e.event === "LogLock");

        // Parse event fields
        const lockEvent = {
            to: event.args._to,
            from: event.args._from,
            symbol: event.args._symbol,
            token: event.args._token,
            value: Number(event.args._value),
            nonce: Number(event.args._nonce),
            argv: argv,
        };

        logging.debug(`lockEvent is ${JSON.stringify(lockEvent, undefined, 2)}`);
        console.log(JSON.stringify(lockEvent, undefined, 0))
    } catch (error) {
        logging.error(error.message);
        // stall so logger has time to write out errors
        await new Promise(resolve => setTimeout(resolve, 5000));
    }
    return cb();
};
