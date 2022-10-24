// Prints all of the token whitelisting events like this:

// npx truffle exec scripts/test/whitelistedTokens.js --json_path /home/james/workspace/sifnode/smart-contracts/deployments/sandpit  --ethereum_network ropsten --bridgebank_address 0x979F0880de42A7aE510829f13E66307aBb957f13
//
// {"token":"0xA3D31ee81Ec2a898B4CF7A67a0851086e4Da7af3","value":true,"symbol":"erowan","name":"erowan"}
// {"token":"0xfA8fC9C22C33FE62BabD5D92DD38Aa27B730d562","value":true,"symbol":"dtoken","name":"dtoken"}

const BN = require('bn.js');

module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.bridgeBankAddressYargOptions,
    });

    const bridgeBankContract = await contractUtilites.buildContract(this, argv, logging, "BridgeBank", argv.bridgebank_address);

    const whitelistUpdates = await bridgeBankContract.getPastEvents("LogWhiteListUpdate", {
        fromBlock: 1,
        toBlock: 'latest'
    });

    const promises = [];
    for (let x of whitelistUpdates) {
        let token = x.returnValues["_token"];
        const promise = contractUtilites.buildContract(this, argv, logging, "BridgeToken", token)
            .then(async tokenContract => {
                const item = {
                    token,
                    value: x.returnValues["_value"],
                    symbol: await tokenContract.symbol(),
                    name: await tokenContract.name(),
                    decimals: (await tokenContract.decimals()).toString(10),
                }
                return item;
            });
        promises.push(promise)
    }
    const result = await Promise.all(promises);
    console.log(JSON.stringify(result));
    return cb();
};
