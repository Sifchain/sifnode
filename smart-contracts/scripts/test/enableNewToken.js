const BN = require('bn.js');

module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.bridgeBankAddressYargOptions,
        ...sifchainUtilities.symbolYargOption,
        ...sifchainUtilities.amountYargOption,
        'limit_amount': {
            describe: 'an amount',
            demandOption: true
        },
    });

    const ethMultiplier = (new BN("10", 10)).pow(new BN(18));

    const amount = new BN(argv.amount, 10);
    const limitAmount = new BN(argv.limit_amount, 10);

    const BridgeToken = artifacts.require("BridgeToken");
    const newToken = await BridgeToken.new(argv.symbol);

    const bridgeBankContract = await contractUtilites.buildContract(this, argv, logging, "BridgeBank", argv.bridgebank_address);

    const accounts = await web3.eth.getAccounts();
    const operator_address = accounts[0]; // operator

    logging.info(`bridgeBankContract.updateEthWhiteList ${newToken.address}, true, from ${operator_address}`);
    const updateWhiteListResult = await bridgeBankContract.updateEthWhiteList(newToken.address, true, {
        from: operator_address
    });
    logging.info(`bridgeBankContract.updateEthWhiteList result ${JSON.stringify(updateWhiteListResult)}`);

    logging.info(`bridgeBankContract.updateTokenLockBurnLimit address ${newToken.address}, limitAmount ${limitAmount}, from ${operator_address}`);
    await bridgeBankContract.updateTokenLockBurnLimit(newToken.address, limitAmount, {
        from: operator_address
    });

    const token_destination = accounts[0];

    logging.info(`newToken.mint to destination ${token_destination}, amount ${amount}, from ${operator_address}`);
    await newToken.mint(token_destination, amount, {
        from: operator_address
    });

    logging.info(`newToken.approve address ${bridgeBankContract.address}, from ${token_destination}`);
    await newToken.approve(bridgeBankContract.address, amount.toString(), {
        from: token_destination
    });

    const result = {
        destination: token_destination,
        "amount": amount.toString(),
        "newtoken_address": newToken.address,
        "newtoken_symbol": argv.symbol,
    }
    console.log(JSON.stringify(result));

    return cb();
};
