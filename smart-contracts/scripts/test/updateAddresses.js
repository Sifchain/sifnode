// Given two files, pull the token addressess from the first file and update the token addresses
// in the second file.
//
// Used to update ui/core/src/assets.ethereum.ropsten.json for the front end.  Use the
// output of yarn integrationtest:whitelistedTokens to get the current addresses.
// (this will be obsolete when the frontend just gets it from the smart contracts
// directly)
//
// For example:
// 
//   node scripts/test/updateAddresses.js scripts/test/updateAddresses.js $BASEDIR/ui/core/src/tokenwhitelist.sandpit.json $BASEDIR/ui/core/src/assets.ethereum.ropsten.json

const fs = require('fs')

const addressFileContents = fs.readFileSync(process.argv[3], 'utf8')
const targetFileContents = fs.readFileSync(process.argv[4], 'utf8')

const addresses = JSON.parse(addressFileContents);

const symbolToToken = {};
for (let x of addresses) {
    symbolToToken[x["symbol"]] = x["token"];
}
const targets = JSON.parse(targetFileContents);
const assets = [];
for (let t of targets["assets"]) {
    t.address = symbolToToken[t["symbol"]];
    assets.push(t);
}
console.log(JSON.stringify({assets: assets}))
