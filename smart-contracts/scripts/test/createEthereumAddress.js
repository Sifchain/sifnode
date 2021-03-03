const BN = require('bn.js');

module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
    });

    const web3 = contractUtilites.buildWeb3(this, argv, logging);

    const newEtherumAccount = web3.eth.accounts.create();

    console.log(JSON.stringify(newEtherumAccount));

    return cb();
};
