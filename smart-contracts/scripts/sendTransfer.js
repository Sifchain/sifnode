/*******************************************
 *** the script just used in local test to generate a new block via trivial amount transfer
  ******************************************/
module.exports = async (cb) => {
  // var sleep = require('sleep')
  const { promisify } = require('util')
  const sleep = promisify(setTimeout)

  let txNumber = 1;

  if (process.argv.length > 4) {
    txNumber = process.argv[4];
  }

  const Web3 = require("web3");

  let provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);

  const web3 = new Web3(provider);
  try {
    const accounts = await web3.eth.getAccounts();
    for (i = 0; i < txNumber; i++) {
      await web3.eth.sendTransaction({from: accounts[8], to: accounts[9], value: 1})
      await sleep(3000)
    }
    
    console.log("Sent transfer transaction...");

  } catch (error) {
    console.error({ error });
  }
  return cb();
};
