const { ethers, upgrades } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");

use(solidity);

describe("Test Bridge Bank", function () {
  let accounts;
  let signerAccounts;
  let operator;
  let owner;
  const consensusThreshold = 75;
  let initialPowers;
  let initialValidators;
  let cosmosBridge;

  before(async function() {
    accounts = await ethers.getSigners();

    signerAccounts = accounts.map((e) => { return e.address });
    operator = accounts[0].address;

    owner = accounts[5].address;
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);
  });

  describe("CosmosBridge", function () {
    it("Should deploy cosmos bridge and bridge bank", async function () {
      const CosmosBridge = await ethers.getContractFactory("CosmosBridge");

      cosmosbridge = await upgrades.deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        initialValidators,
        initialPowers
      ]);
      await cosmosbridge.deployed();

      expect(
        (await cosmosbridge.consensusThreshold()).toString()
      ).to.equal(consensusThreshold.toString());

      for (let i = 0; i < initialValidators.length; i++) {
        const address = initialValidators[i];

        expect(
          await cosmosbridge.isActiveValidator(address)
        ).to.be.true;
        
        expect(
          (await cosmosbridge.getValidatorPower(address)).toString()
        ).to.equal("25");
      }

      const BridgeBank = await ethers.getContractFactory("BridgeBank");

      const bridgeBank = await upgrades.deployProxy(BridgeBank, [
        operator,
        cosmosbridge.address,
        owner,
        pauser
      ]);

      await bridgeBank.deployed();

      expect(await bridgeBank.cosmosBridge()).to.be.equal(cosmosbridge.address);
      expect(await bridgeBank.operator()).to.be.equal(operator);
      expect(await bridgeBank.owner()).to.be.equal(owner);
      expect(await bridgeBank.pausers(pauser)).to.be.true;
    });
  });

    // describe("BridgeBank", function () {
    //   it("Should be able to deploy BridgeBank", async function () {
        // const BridgeBank = await ethers.getContractFactory("BridgeBank");

        // const bridgeBank = await upgrades.deployProxy(BridgeBank, [
        //   operator,
        //   cosmosBridge.address,
        //   owner,
        //   pauser
        // ]);

        // expect(await bridgeBank.cosmosBridge()).to.be.equal(cosmosBridge.address);
        // expect(await bridgeBank.operator()).to.be.equal(operator);
        // expect(await bridgeBank.isOwner(owner)).to.be.true;
        // expect(await bridgeBank.pausers(pauser)).to.be.true;
    //   });
    // });
});