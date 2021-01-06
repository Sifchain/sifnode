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

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("BridgeBank", function (accounts) {
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
      this.valset = await deployProxy(Valset,
        [
          operator,
          this.initialValidators,
          this.initialPowers
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [operator, this.valset.address], {unsafeAllowCustomTypes: true});

      // Deploy Oracle contract
      this.oracle = await deployProxy(Oracle,
        [
          operator,
          this.valset.address,
          this.cosmosBridge.address,
          consensusThreshold
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank,
        [
          operator,
          this.oracle.address,
          this.cosmosBridge.address,
          operator
        ],
        {unsafeAllowCustomTypes: true}
      );
    });

    it("should deploy the BridgeBank, correctly setting the operator and valset", async function () {
      this.bridgeBank.should.exist;

      const bridgeBankOperator = await this.bridgeBank.operator();
      bridgeBankOperator.should.be.equal(operator);

      const bridgeBankOracle = await this.bridgeBank.oracle();
      bridgeBankOracle.should.be.equal(this.oracle.address);
    });

    it("should correctly set initial values", async function () {
      // EthereumBank initial values
      const bridgeLockBurnNonce = Number(await this.bridgeBank.lockBurnNonce());
      bridgeLockBurnNonce.should.be.bignumber.equal(0);

      // CosmosBank initial values
      const bridgeTokenCount = Number(await this.bridgeBank.bridgeTokenCount());
      bridgeTokenCount.should.be.bignumber.equal(0);
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
      this.valset = await deployProxy(Valset,
        [
          operator,
          this.initialValidators,
          this.initialPowers
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [operator, this.valset.address], {unsafeAllowCustomTypes: true});

      // Deploy Oracle contract
      this.oracle = await deployProxy(Oracle,
        [
          operator,
          this.valset.address,
          this.cosmosBridge.address,
          consensusThreshold
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank,
        [
          operator,
          this.oracle.address,
          this.cosmosBridge.address,
          operator
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Set oracle and bridge bank for the cosmos bridge
      await this.cosmosBridge.setOracle(this.oracle.address, {from: operator})
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
    it("should not allow initialization of oracle with a consensus threshold over 100", async function () {
      this.oracle = await Oracle.new();
      await expectRevert(
        this.oracle.initialize(
          accounts[0],
          accounts[0],
          accounts[0],
          101
        ),
        "Invalid consensus threshold."
      );
    });

    it("should not allow initialization of oracle with a consensus threshold of 0", async function () {
      this.oracle = await Oracle.new();
      await expectRevert(
        this.oracle.initialize(
          accounts[0],
          accounts[0],
          accounts[0],
          0
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
      this.valset = await deployProxy(Valset,
        [
          operator,
          this.initialValidators,
          this.initialPowers
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [operator, this.valset.address], {unsafeAllowCustomTypes: true});

      // Deploy Oracle contract
      this.oracle = await deployProxy(Oracle,
        [
          operator,
          this.valset.address,
          this.cosmosBridge.address,
          consensusThreshold
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank,
        [
          operator,
          this.oracle.address,
          this.cosmosBridge.address,
          operator
        ],
        {unsafeAllowCustomTypes: true}
      );

      // Set oracle and bridge bank for the cosmos bridge
      await this.cosmosBridge.setOracle(this.oracle.address, {from: operator})
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {from: operator})
    });

    it("should not allow a non operator to call the function", async function () {
      await expectRevert(
        this.bridgeBank.bulkWhitelistUpdateLimits([], [], {from: userOne}),
        "Must be BridgeBank operator."
      );
    });

    it("should not allow arrays of different sizes", async function () {
      await expectRevert(
        this.bridgeBank.bulkWhitelistUpdateLimits([], [1], {from: operator}),
        "!same length"
      );
    });

    it("Should allow bulk whitelisting", async function () {
      const addresses = [];
      const nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

      // create tokens and address array
      for (let i = 0; i < 10; i++) {
        const bridgeToken = await BridgeToken.new("eRowan" + i.toString());
        addresses.push(bridgeToken.address);
      }

      await this.bridgeBank.bulkWhitelistUpdateLimits(addresses, nums, {from: operator});

      // query each token in the array and make sure that the limit is correct
      for (let i = 0; i < 10; i++) {
        const limit = Number(await this.bridgeBank.maxTokenAmount("eRowan" + i.toString()));
        expect(limit).to.be.equal(nums[i]);
      }
    });
  });
});
