const {web3} = require("@openzeppelin/test-helpers/src/setup");
const { time } = require("@openzeppelin/test-helpers");
require('@openzeppelin/test-helpers/configure')({
  provider: process.env.LOCAL_PROVIDER,
});

/*******************************************
 *** the script just used in local test to generate a new block via trivial amount transfer
 ******************************************/
console.log("Expected usage: \n node scripts/advanceBlock.js <block_number>");
console.log("Provider: ", process.env.LOCAL_PROVIDER);

async function main() {
  let txNumber = 5;

  // [node, <script_name>, <args>]
  if (process.argv.length >= 3) {
    txNumber = process.argv[2];
  }

  try {
    for (let i = 0; i < txNumber; i++) {
      await time.advanceBlock();
    }

    console.log(`Advanced ${txNumber} blocks`);

    let bn = await web3.eth.getBlockNumber();

    console.log(`current block number is ${bn}`)

    console.log(JSON.stringify({nBlocks: txNumber, currentBlockNumber: bn}))
  } catch (error) {
    console.error({ error });
  }
}

main()