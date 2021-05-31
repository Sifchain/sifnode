const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers, upgrades } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

use(solidity);

const getBalance = async function(address) {
  return await network.provider.send("eth_getBalance", [address]);
}

describe("Test Cosmos Bridge", function () {
  let userOne;
  let userTwo;
  let userThree;
  let userFour;
  let accounts;
  let signerAccounts;
  let operator;
  let owner;
  const consensusThreshold = 75;
  let initialPowers;
  let initialValidators;
  let CosmosBridge;
  let BridgeBank;
  let BridgeToken;

  before(async function() {
    CosmosBridge = await ethers.getContractFactory("CosmosBridge");
    BridgeBank = await ethers.getContractFactory("BridgeBank");
    BridgeToken = await ethers.getContractFactory("BridgeToken");
    
    accounts = await ethers.getSigners();
    
    signerAccounts = accounts.map((e) => { return e.address });

    operator = accounts[0].address;
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[7].address;

    owner = accounts[5].address;
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);
  });

  beforeEach(async function () {
    // Deploy Valset contract
    this.initialValidators = initialValidators;
    this.initialPowers = initialPowers;

    // Deploy CosmosBridge contract
    this.cosmosBridge = await upgrades.deployProxy(CosmosBridge, [
      operator,
      consensusThreshold,
      initialValidators,
      initialPowers
    ]);
    await this.cosmosBridge.deployed();

    // Deploy BridgeBank contract
    this.bridgeBank = await upgrades.deployProxy(BridgeBank, [
      this.cosmosBridge.address,
      owner,
      pauser
    ]);
    await this.bridgeBank.deployed();

    // Operator sets Bridge Bank
    await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
      from: operator
    });

    // This is for ERC20 deposits
    this.sender = web3.utils.utf8ToHex(
      "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
    );
    this.senderSequence = 1;
    this.recipient = userThree;
    this.name = "TEST COIN";
    this.symbol = "TEST";
    this.ethereumToken = "0x0000000000000000000000000000000000000000";
    this.weiAmount = web3.utils.toWei("0.25", "ether");

    this.token = await BridgeToken.deploy(
      this.name,
      this.symbol,
      18
    );

    await this.token.deployed();
    this.amount = 100;
    //Load user account with ERC20 tokens for testing
    await this.token.mint(userOne.address, this.amount * 2, {
      from: operator
    }).should.be.fulfilled;

    // Approve tokens to contract
    await this.token.connect(userOne).approve(this.bridgeBank.address, this.amount).should.be.fulfilled;
      
    // Lock tokens on contract
    await this.bridgeBank.connect(userOne).lock(
      this.sender,
      this.token.address,
      this.amount
    ).should.be.fulfilled;

    // Lock tokens on contract
    await this.bridgeBank.connect(userOne).lock(
      this.sender,
      this.ethereumToken,
      this.amount, {
        value: this.amount
      }
    ).should.be.fulfilled;
  });

  describe("CosmosBridge", function () {
    it("should return true if a sifchain address prefix is correct", async function () {
      (await this.bridgeBank.VSA(this.sender)).should.be.equal(true);
    });

    it("should return false if a sifchain address length is incorrect", async function () {
      const incorrectSifAddress = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpaceee"
      );
      (await this.bridgeBank.VSA(incorrectSifAddress)).should.be.equal(false);
    });

    it("should return false if a sifchain address has an incorrect `sif` prefix", async function () {
      const incorrectSifAddress = web3.utils.utf8ToHex(
        "eif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      (await this.bridgeBank.VSA(incorrectSifAddress)).should.be.equal(false);
    });

    it("Should deploy cosmos bridge and bridge bank", async function () {
      expect(
        (await this.cosmosBridge.consensusThreshold()).toString()
      ).to.equal(consensusThreshold.toString());

      // iterate over all validators and ensure they have the proper
      // powers and that they have been succcessfully whitelisted
      for (let i = 0; i < initialValidators.length; i++) {
        const address = initialValidators[i];

        expect(
          await this.cosmosBridge.isActiveValidator(address)
        ).to.be.true;
        
        expect(
          (await this.cosmosBridge.getValidatorPower(address)).toString()
        ).to.equal("25");
      }

      expect(await this.bridgeBank.cosmosBridge()).to.be.equal(this.cosmosBridge.address);
      expect(await this.bridgeBank.owner()).to.be.equal(owner);
      expect(await this.bridgeBank.pausers(pauser)).to.be.true;
    });

    it("should mint bridge tokens upon the successful processing of a burn prophecy claim", async function () {
      const beforeUserBalance = Number(
        await this.token.balanceOf(this.recipient)
      );
      beforeUserBalance.should.be.bignumber.equal(Number(0));

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      this.nonce = 1;
      let receipt = (
        await this.cosmosBridge.connect(userOne).newProphecyClaim(
          this.sender,
          this.senderSequence,
          this.recipient,
          this.token.address,
          this.amount,
          false
        ).should.be.fulfilled
      );

      receipt = await this.cosmosBridge.connect(userTwo).newProphecyClaim(
        this.sender,
        this.senderSequence,
        this.recipient,
        this.token.address,
        this.amount,
        false
      ).should.be.fulfilled;

      receipt = await this.cosmosBridge.connect(userFour).newProphecyClaim(
        this.sender,
        this.senderSequence,
        this.recipient,
        this.token.address,
        this.amount,
        false
      ).should.be.fulfilled;

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await this.token.balanceOf(this.recipient)
      );
      afterUserBalance.should.be.bignumber.equal(this.amount);
    });

    it("should unlock eth upon the successful processing of a burn prophecy claim", async function () {
      // assert recipient balance before receiving proceeds from newProphecyClaim is correct
      const recipientStartingBalance = await getBalance(this.recipient);
      const recipientCurrentBalance = Web3Utils.fromWei(recipientStartingBalance);

      expect(recipientCurrentBalance).to.be.equal(
        "10000"
      );

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      this.nonce = 1;
      let receipt = await this.cosmosBridge.connect(userOne).newProphecyClaim(
        this.sender,
        this.senderSequence,
        this.recipient,
        this.ethereumToken,
        this.amount,
        false
      ).should.be.fulfilled;
      
      receipt = await this.cosmosBridge.connect(userTwo).newProphecyClaim(
        this.sender,
        this.senderSequence,
        this.recipient,
        this.ethereumToken,
        this.amount,
        false
      ).should.be.fulfilled;

      receipt = await this.cosmosBridge.connect(userFour).newProphecyClaim(
        this.sender,
        this.senderSequence,
        this.recipient,
        this.ethereumToken,
        this.amount,
        false
      ).should.be.fulfilled;

      const recipientEndingBalance = await getBalance(this.recipient);
      const recipientBalance = Web3Utils.fromWei(recipientEndingBalance);

      expect(recipientBalance).to.be.equal(
        "10000.0000000000000001"
      );
    });
  });
});
