const BN = require('bn.js');

module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.transactionYargOptions
    });

    logging.debug(`sendLockTx arguments: ${JSON.stringify(argv, undefined, 2)}`);

    const bridgeBankContract = await contractUtilites.buildContract(this, argv, logging, "BridgeBank", argv.bridgebank_address);

    let cosmosRecipient = Web3.utils.utf8ToHex(argv.sifchain_address);
    let coinDenom = argv.symbol;
    let amount = new BN(argv.amount);

    let request = {
        from: argv.ethereum_address,
        value: coinDenom === sifchainUtilities.NULL_ADDRESS ? amount : 0,
    };

    if (argv.gas === 'estimate') {
        const estimate = await bridgeBankContract.lock.estimateGas(
            cosmosRecipient,
            coinDenom,
            amount,
            request,
        );
        // increase by 10%
        request.gas = new BN(estimate).mul(new BN(11)).div(new BN(10));
    } else {
        request.gas = argv.gas;
    }

    const lockResult = await bridgeBankContract.lock(cosmosRecipient, coinDenom, amount, request);

    console.log(JSON.stringify(lockResult, undefined, 0))

    return cb();
};
