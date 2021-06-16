const AliceToken = artifacts.require("AliceToken");
const BobToken = artifacts.require("BobToken");
const UsdCoin = artifacts.require("UsdCoin");
const LinkCoin = artifacts.require("LinkCoin");

module.exports = function (deployer) {
  deployer.deploy(AliceToken);
  deployer.deploy(BobToken);
  deployer.deploy(UsdCoin);
  deployer.deploy(LinkCoin);
};
