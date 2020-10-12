const AliceToken = artifacts.require("AliceToken");
const BobToken = artifacts.require("BobToken");

module.exports = function(deployer) {
  deployer.deploy(AliceToken);
  deployer.deploy(BobToken);
};
