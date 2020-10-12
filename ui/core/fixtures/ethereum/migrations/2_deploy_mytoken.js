const AliceToken = artifacts.require("AliceToken");
const BobToken = artifacts.require("BobToken");
const fs = require("fs");
const path = require("path");
module.exports = function(deployer) {
  deployer.deploy(AliceToken);
  deployer.deploy(BobToken);
  deployer.then(() => {
    const atk = AliceToken.address;
    const btk = BobToken.address;
    fs.writeFileSync(
      path.resolve(__dirname, "../ctx.json"),
      JSON.stringify({ atk, btk })
    );
  });
};

// Q: How do I know the address of these contracts
// in our dev frontend app?

// A: Write a JSON file that returns them
