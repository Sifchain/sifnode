module.exports = async (cb) => {
    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const waitForBlocksArgs = {
        'block_number': {
            type: "number",
            demandOption: true
        },
    };
    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...waitForBlocksArgs
    });

    const web3 = contractUtilites.buildWeb3(this, argv);

    const logging = sifchainUtilities.configureLogging(this);

    let waitTime = 2000;
    switch(argv.ethereum_network) {
        case "ropsten":
        case "mainnet":
            waitTime = 60000;
            break;
    }
    try {
        let blockNumber = await web3.eth.getBlockNumber();
        do {
            blockNumber = await web3.eth.getBlockNumber();
            logging.debug(`waiting for block ${argv.block_number}, current block is ${blockNumber}`)
            await new Promise(resolve => setTimeout(resolve, 2000));
        } while (blockNumber < argv.block_number);
    } catch (error) {
        console.log(error);
        // stall so logger has time to write out errors
        await new Promise(resolve => setTimeout(resolve, 2000));
    }

    return cb();
};
