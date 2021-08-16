const { deployProxy, silenceWarnings } = require('@openzeppelin/truffle-upgrades');

const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeToken = artifacts.require("BridgeToken");
const BridgeBank = artifacts.require("BridgeBank");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

const {
  BN,
  expectRevert, // Assertions for transactions that should fail
} = require('@openzeppelin/test-helpers');
const { expect } = require('chai');

const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";
const sifRecipient = web3.utils.utf8ToHex(
  "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
);

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Security Test", function (accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  // Consensus threshold of 70%
  const consensusThreshold = 70;

  describe("BridgeBank Security", function () {
    beforeEach(async function () {
      await silenceWarnings();

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        this.initialValidators,
        this.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        this.cosmosBridge.address,
        operator,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      this.token = await BridgeToken.new("erowan");

      await this.bridgeBank.addExistingBridgeToken(this.token.address, { from: operator });
    });

    it("should deploy the BridgeBank, correctly setting the operator and valset", async function () {
      this.bridgeBank.should.exist;

      const bridgeBankOperator = await this.bridgeBank.operator();
      bridgeBankOperator.should.be.equal(operator);
    });

    it("should be able to change the owner", async function () {
      expect(await this.bridgeBank.owner()).to.be.equal(operator);
      await this.bridgeBank.changeOwner(userTwo, { from: operator });
      expect(await this.bridgeBank.owner()).to.be.equal(userTwo);
    });

    it("should not be able to change the owner if the caller is not the owner", async function () {
      expect(await this.bridgeBank.owner()).to.be.equal(operator);
      await expectRevert(
        this.bridgeBank.changeOwner(userTwo, { from: userThree }),
        "!owner"
      );
      expect((await this.bridgeBank.owner())).to.be.equal(operator);
    });

    it("should be able to change the operator", async function () {
      expect((await this.bridgeBank.operator())).to.be.equal(operator);
      await this.bridgeBank.changeOperator(userTwo, { from: operator });
      expect((await this.bridgeBank.operator())).to.be.equal(userTwo);
    });

    it("should not be able to change the operator if the caller is not the operator", async function () {
      expect((await this.bridgeBank.operator())).to.be.equal(operator);
      await expectRevert(
        this.bridgeBank.changeOperator(userTwo, { from: userThree }),
        "!operator"
      );
      expect((await this.bridgeBank.operator())).to.be.equal(operator);
    });

    it("should correctly set initial values", async function () {
      // CosmosBank initial values
      const bridgeTokenCount = Number(await this.bridgeBank.bridgeTokenCount());
      bridgeTokenCount.should.be.bignumber.equal(1);
    });

    it("should be able to pause the contract", async function () {
      await this.bridgeBank.pause();
      expect(await this.bridgeBank.paused()).to.be.true;
    });

    it("should not be able to pause the contract if you are not the owner", async function () {
      await expectRevert(
        this.bridgeBank.pause({ from: userOne }),
        "PauserRole: caller does not have the Pauser role"
      );
      expect(await this.bridgeBank.paused()).to.be.false;
    });

    it("should be able to add a new pauser if you are a pauser", async function () {
      expect(await this.bridgeBank.pausers(operator)).to.be.true;
      expect(await this.bridgeBank.pausers(userOne)).to.be.false;
      await this.bridgeBank.addPauser(userOne, { from: operator })
      expect(await this.bridgeBank.pausers(operator)).to.be.true;
      expect(await this.bridgeBank.pausers(userOne)).to.be.true;
    });

    it("should be able to renounce yourself as pauser", async function () {
      expect(await this.bridgeBank.pausers(operator)).to.be.true;
      expect(await this.bridgeBank.pausers(userOne)).to.be.false;
      await this.bridgeBank.addPauser(userOne, { from: operator })
      expect(await this.bridgeBank.pausers(operator)).to.be.true;
      expect(await this.bridgeBank.pausers(userOne)).to.be.true;
      await this.bridgeBank.renouncePauser({ from: userOne });
      expect(await this.bridgeBank.pausers(userOne)).to.be.false;
    });

    it("should be able to pause and then unpause the contract", async function () {
      // CosmosBank initial values
      await expectRevert(
        this.bridgeBank.unpause(),
        "Pausable: not paused"
      );
      await this.bridgeBank.pause();
      await expectRevert(
        this.bridgeBank.pause(),
        "Pausable: paused"
      );
      expect(await this.bridgeBank.paused()).to.be.true;
      await this.bridgeBank.unpause();
      expect(await this.bridgeBank.paused()).to.be.false;
    });
    
    it("should not be able to lock when contract is paused", async function () {
      await this.bridgeBank.pause();
      expect(await this.bridgeBank.paused()).to.be.true;

      await expectRevert(
        this.bridgeBank.lock(sifRecipient, NULL_ADDRESS, 100),
        "Pausable: paused"
      );
    });
    
    it("should not be able to burn when contract is paused", async function () {
      await this.bridgeBank.pause();
      expect(await this.bridgeBank.paused()).to.be.true;

      await expectRevert(
        this.bridgeBank.burn(sifRecipient, this.token.address, 100),
        "Pausable: paused"
      );
    });

    it("should not allow a user to send ethereum directly to the contract", async function () {
      await this.bridgeBank
        .send(Web3Utils.toWei("0.25", "ether"), {
          from: userOne
        })
        .should.be.rejectedWith(EVMRevert);
    });
  });

  // This entire scenario is mimicking the mainnet scenario where there will be
  // cosmos assets on sifchain, and then we hook into an existing ERC20 contract on mainnet
  // that is eRowan. Then we will try to transfer rowan to eRowan to ensure that
  // everything is set up correctly.
  // We will do this by making a new prophecy claim, validating it with the validators
  // Then ensure that the prohpecy claim paid out the person that it was supposed to
  describe("Bridge token burning", function () {
    before(async function () {
      // this test needs to create a new token contract that will
      // effectively be able to be treated as if it was a cosmos native asset
      // even though it was created on top of ethereum

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [33, 33, 33];

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        this.initialValidators,
        this.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        this.cosmosBridge.address,
        operator,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {from: operator})
    });

    it("should not allow burning of non whitelisted token address", async function () {
      function convertToHex(str) {
        let hex = '';
        for (let i = 0; i < str.length; i++) {
            hex += '' + str.charCodeAt(i).toString(16);
        }
        return hex;
      }

      const symbol = 'eRowan'
      const amount = 100000;
      const sifAddress = "0x" + convertToHex("sif12qfvgsq76eghlagyfcfyt9md2s9nunsn40zu2h");
      
      // create new fake eRowan token
      const bridgeToken = await BridgeToken.new("eRowan");

      // Attempt to burn tokens
      await expectRevert(
        this.bridgeBank.burn(
            sifAddress,
            bridgeToken.address,
            amount, { from: operator }
        ),
        "Only token in whitelist can be transferred to cosmos"
      );
    });
  });

  describe("Consensus Threshold Limits", function () {
    beforeEach(async function () {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [33, 33, 33];
    });

    it("should not allow initialization of CosmosBridge with a consensus threshold over 100", async function () {
      this.bridge = await CosmosBridge.new();
      await expectRevert(
        this.bridge.initialize(
          operator,
          101,
          this.initialValidators,
          this.initialPowers
        ),
        "Invalid consensus threshold."
      );
    });

    it("should not allow initialization of oracle with a consensus threshold of 0", async function () {
      this.bridge = await CosmosBridge.new();
      await expectRevert(
        this.bridge.initialize(
          operator,
          0,
          this.initialValidators,
          this.initialPowers
        ),
        "Consensus threshold must be positive."
      );
    });
  });

  describe("Bulk whitelist and limit add", function () {
    before(async function () {
      // this test needs to create a new token contract that will
      // effectively be able to be treated as if it was a cosmos native asset
      // even though it was created on top of ethereum

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [33, 33, 33];

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        this.initialValidators,
        this.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        this.cosmosBridge.address,
        operator,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {from: operator});
    });

    it("should not allow a non operator to call the function", async function () {
      await expectRevert(
        this.bridgeBank.bulkWhitelistUpdateLimits([], {from: userOne}),
        "!operator"
      );
    });

    it("Should allow bulk whitelisting", async function () {
      const addresses = [];

      // create tokens and address array
      for (let i = 0; i < 10; i++) {
        const bridgeToken = await BridgeToken.new("eRowan" + i.toString());
        addresses.push(bridgeToken.address);
      }

      await this.bridgeBank.bulkWhitelistUpdateLimits(addresses, {from: operator});

      // query each token in the array and make sure that the limit is correct
      for (let i = 0; i < 10; i++) {
        const isWhitelisted = await this.bridgeBank.getTokenInEthWhiteList(addresses[i]);
        expect(isWhitelisted).to.be.equal(true);
      }
    });
  });
});
