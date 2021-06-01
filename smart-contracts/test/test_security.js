const { deployProxy } = require('@openzeppelin/truffle-upgrades');

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const web3 = require("web3");
const BigNumber = web3.BigNumber;
const { expect } = require('chai');
const { multiTokenSetup } = require("./helpers/testFixture");

const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";
const sifRecipient = web3.utils.utf8ToHex(
  "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
);

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

describe("Security Test", function () {
  let userOne;
  let userTwo;
  let userThree;
  let userFour;
  let accounts;
  let signerAccounts;
  let operator;
  let owner;
  const consensusThreshold = 70;
  let initialPowers;
  let initialValidators;
  // track the state of the deployed contracts
  let state;
  let CosmosBridge;
  let BridgeToken;

  before(async function() {
    CosmosBridge = await ethers.getContractFactory("CosmosBridge");
    BridgeToken = await ethers.getContractFactory("BridgeToken");

    accounts = await ethers.getSigners();
    
    signerAccounts = accounts.map((e) => { return e.address });

    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[7].address;

    owner = accounts[5];
    pauser = accounts[6];

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);
  });


  describe("BridgeBank Security", function () {
    beforeEach(async function () {
      state = await multiTokenSetup(
        initialValidators,
        initialPowers,
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser.address
      );
    });

    it("should be able to change the owner", async function () {
      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);
      await state.bridgeBank.connect(owner).changeOwner(userTwo.address);
      expect(await state.bridgeBank.owner()).to.be.equal(userTwo.address);
    });

    it("should not be able to change the owner if the caller is not the owner", async function () {
      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);

      await expect(
        state.bridgeBank.connect(accounts[7])
          .changeOwner(userTwo.address),
      ).to.be.revertedWith("!owner");

      expect((await state.bridgeBank.owner())).to.be.equal(owner.address);
    });

    it("should not be able to change the operator if the caller is not the operator", async function () {
      expect((await state.cosmosBridge.operator())).to.be.equal(operator.address);
      await expect(
        state.cosmosBridge.connect(userOne)
          .changeOperator(userTwo.address),
      ).to.be.revertedWith("Must be the operator.");

      expect((await state.cosmosBridge.operator())).to.be.equal(operator.address);
    });

    it("should correctly set initial values", async function () {
      // CosmosBank initial values
      // bridgeTokenCount is deprecated
      const bridgeTokenCount = Number(await state.bridgeBank.bridgeTokenCount());
      bridgeTokenCount.should.be.bignumber.equal(0);
    });

    it("should be able to pause the contract", async function () {
      await state.bridgeBank.connect(pauser).pause();
      expect(await state.bridgeBank.paused()).to.be.true;
    });

    it("should not be able to pause the contract if you are not the owner", async function () {
      await expect(
        state.bridgeBank.connect(userOne).pause(),
      ).to.be.revertedWith("PauserRole: caller does not have the Pauser role");

      expect(await state.bridgeBank.paused()).to.be.false;
    });

    it("should be able to add a new pauser if you are a pauser", async function () {
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.false;

      await state.bridgeBank.connect(pauser).addPauser(userOne.address);

      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.true;
    });

    it("should be able to renounce yourself as pauser", async function () {
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.false;

      await state.bridgeBank.connect(pauser).addPauser(userOne.address);
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.true;

      await state.bridgeBank.connect(userOne).renouncePauser();
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.false;
    });

    it("should be able to pause and then unpause the contract", async function () {
      // CosmosBank initial values
      await expect(
        state.bridgeBank.connect(pauser).unpause(),
      ).to.be.revertedWith("Pausable: not paused");

      await state.bridgeBank.connect(pauser).pause();
      await expect(
        state.bridgeBank.connect(pauser).pause(),
      ).to.be.revertedWith("Pausable: paused");

      expect(await state.bridgeBank.paused()).to.be.true;
      await state.bridgeBank.connect(pauser).unpause();

      expect(await state.bridgeBank.paused()).to.be.false;
    });
    
    it("should not be able to lock when contract is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      expect(await state.bridgeBank.paused()).to.be.true;

      await expect(
        state.bridgeBank.connect(userOne)
          .lock(sifRecipient, NULL_ADDRESS, 100),
      ).to.be.revertedWith("Pausable: paused");
    });
    
    it("should not be able to burn when contract is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      expect(await state.bridgeBank.paused()).to.be.true;

      await expect(
        state.bridgeBank.connect(userOne)
          .burn(sifRecipient, state.rowan.address, 100),
      ).to.be.revertedWith("Pausable: paused");
    });
  });

  // state entire scenario is mimicking the mainnet scenario where there will be
  // cosmos assets on sifchain, and then we hook into an existing ERC20 contract on mainnet
  // that is eRowan. Then we will try to transfer rowan to eRowan to ensure that
  // everything is set up correctly.
  // We will do state by making a new prophecy claim, validating it with the validators
  // Then ensure that the prohpecy claim paid out the person that it was supposed to
  describe("Bridge token burning", function () {
    before(async function () {
      // state test needs to create a new token contract that will
      // effectively be able to be treated as if it was a cosmos native asset
      // even though it was created on top of ethereum

      // Deploy Valset contract
      state.initialValidators = [userOne.address, userTwo.address, userThree];
      state.initialPowers = [33, 33, 33];

      state = await multiTokenSetup(
        state.initialValidators,
        state.initialPowers,
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser.address
      );
    });

    it("should not allow burning of non whitelisted token address", async function () {
      function convertToHex(str) {
        let hex = '';
        for (let i = 0; i < str.length; i++) {
            hex += '' + str.charCodeAt(i).toString(16);
        }
        return hex;
      }

      const amount = 100000;
      const sifAddress = "0x" + convertToHex("sif12qfvgsq76eghlagyfcfyt9md2s9nunsn40zu2h");
      
      // create new fake eRowan token
      const bridgeToken = await BridgeToken.deploy("rowan", "rowan", 18);

      // Attempt to burn tokens
      await expect(
        state.bridgeBank.connect(operator).burn(
          sifAddress,
          bridgeToken.address,
          amount
        ),
      ).to.be.revertedWith("Only token in whitelist can be transferred to cosmos");
    });
  });

  describe("Consensus Threshold Limits", function () {
    beforeEach(async function () {
      state.initialValidators = [userOne.address, userTwo.address, userThree];
      state.initialPowers = [33, 33, 33];

      state = await multiTokenSetup(
        state.initialValidators,
        state.initialPowers,
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser.address
      );
    });

    it("should not allow initialization of CosmosBridge with a consensus threshold over 100", async function () {
      state.bridge = await CosmosBridge.deploy();

      await expect(
        state.bridge.connect(operator).initialize(
          operator.address,
          101,
          state.initialValidators,
          state.initialPowers
        ),
      ).to.be.revertedWith("Invalid consensus threshold.");
    });

    it("should not allow initialization of oracle with a consensus threshold of 0", async function () {
      state.bridge = await CosmosBridge.deploy();
      await expect(
        state.bridge.connect(operator).initialize(
          operator.address,
          0,
          state.initialValidators,
          state.initialPowers
        ),
      ).to.be.revertedWith("Consensus threshold must be positive.");
    });
  });

  describe("Bulk whitelist and limit add", function () {
    before(async function () {
      // state test needs to create a new token contract that will
      // effectively be able to be treated as if it was a cosmos native asset
      // even though it was created on top of ethereum

      // Deploy Valset contract
      state.initialValidators = [userOne, userTwo, userThree];
      state.initialPowers = [33, 33, 33];

      // Deploy CosmosBridge contract
      state.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        state.initialValidators,
        state.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      state.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        state.cosmosBridge.address,
        operator,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      await state.cosmosBridge.setBridgeBank(state.bridgeBank.address, {from: operator});
    });
  });
});
