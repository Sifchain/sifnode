// Given two files, pull the token addressess from the first file and update the token addresses
// in the second file.
//
// Used to update ui/core/src/assets.ethereum.ropsten.json for the front end.  Use the
// output of yarn integrationtest:whitelistedTokens to get the current addresses.
// (this will be obsolete when the frontend just gets it from the smart contracts
// directly)
//
// For example, to get the sifchain version:
// 
//   node scripts/test/updateAddresses.js $BASEDIR/ui/core/src/tokenwhitelist.sandpit.json $BASEDIR/ui/core/src/assets.ethereum.ropsten.json | jq .sifchain
//
//

const fs = require('fs')

const addressFileContents = fs.readFileSync(process.argv[2], 'utf8')
const targetFileContents = fs.readFileSync(process.argv[3], 'utf8')

// Build a map of symbols (like usdt) to hex addresses
const addresses = JSON.parse(addressFileContents);

const symbolToHexAddress = {};
for (let x of addresses) {
    symbolToHexAddress[x["symbol"]] = x["token"];
}

const targets = JSON.parse(targetFileContents);

const assets = targets["assets"].map(t => {
    const newElement = {
        ...t,
        address: symbolToHexAddress[t["symbol"].toLowerCase()],
    }
    return newElement;
});

const sifchainAssets = assets.map(t => {
    return {
        ...t,
        symbol: ((t.symbol === "erowan") ? "rowan" : `c${t.symbol}`).toLowerCase(),
        network: "sifchain",
    }
});

const result = {
    ethereum: {assets},
    sifchain: {assets: sifchainAssets},
}

console.log(JSON.stringify(result))
