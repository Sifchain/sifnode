const { time } = require("@openzeppelin/test-helpers");
require('@openzeppelin/test-helpers/configure')({
  provider: process.env.LOCAL_PROVIDER,
});

/*******************************************
 *** the script just used in local test to generate a new block via trivial amount transfer
  ******************************************/
module.exports = async (cb) => {
  // default is to advance 5 blocks
  let txNumber = 5;

  if (process.argv.length > 4) {
    txNumber = process.argv[4];
  }

  try {
    for (i = 0; i < txNumber; i++) {
      await time.advanceBlock();
    }
    
    console.log(`Advanced ${txNumber} blocks`);

  } catch (error) {
    console.error({ error });
  }
  return cb();
};
