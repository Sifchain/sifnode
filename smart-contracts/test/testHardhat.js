const EVMRevert = "revert";
const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers, upgrades } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const {
  BN,           // Big Number support
  constants,    // Common constants, like the zero address and largest integers
  expectEvent,  // Assertions for emitted events
  expectRevert, // Assertions for transactions that should fail
} = require('@openzeppelin/test-helpers');
const { after } = require("lodash");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

use(solidity);

// Contract's enum ClaimType can be represented as a sequence of integers
const CLAIM_TYPE_BURN = 1;
const CLAIM_TYPE_LOCK = 2;

const getBalance = async function(address) {
  return await network.provider.send("eth_getBalance", [address]);
}

describe("Test Bridge Bank", function () {
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
    this.symbol = "TEST";
    this.ethereumToken = "0x0000000000000000000000000000000000000000";
    this.weiAmount = web3.utils.toWei("0.25", "ether");

    this.token = await BridgeToken.deploy(this.symbol);
    await this.token.deployed();
    this.amount = 100;
    //Load user account with ERC20 tokens for testing
    await this.token.mint(userOne.address, this.amount, {
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
  });

  describe("CosmosBridge", function () {
    it("should return true if a sifchain address prefix is correct", async function () {
      (await this.bridgeBank.verifySifPrefix(this.sender)).should.be.equal(true);
    });

    it("should return false if a sifchain address has an incorrect `sif` prefix", async function () {
      const incorrectSifAddress = web3.utils.utf8ToHex(
        "eif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      (await this.bridgeBank.verifySifPrefix(incorrectSifAddress)).should.be.equal(false);
    });

    it("Should deploy cosmos bridge and bridge bank", async function () {
      expect(
        (await this.cosmosBridge.consensusThreshold()).toString()
      ).to.equal(consensusThreshold.toString());

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

    it("should not allow users to lock ERC20 tokens if the sifaddress prefix is incorrect", async function () {
      const invalidSifAddress = web3.utils.utf8ToHex(
        "zif1gdnl9jj2xgy5n04r7heqxlqvvzcy24zc96ns2f"
      );
      // Attempt to lock tokens
      await expect(this.bridgeBank.connect(userOne).lock(
          invalidSifAddress,
          this.token.address,
          this.amount
        )
      ).to.be.revertedWith("Invalid sif address");
    });

    it("should mint bridge tokens upon the successful processing of a burn prophecy claim", async function () {
      const beforeUserBalance = Number(
        await this.token.balanceOf(this.recipient)
      );
      beforeUserBalance.should.be.bignumber.equal(Number(0));

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      this.nonce = 1;
      let receipt = await this.cosmosBridge.connect(userOne).newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.symbol,
        this.senderSequence,
        this.recipient,
        this.token.address,
        this.amount
      ).should.be.fulfilled;

      receipt = await this.cosmosBridge.connect(userTwo).newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.symbol,
        this.senderSequence,
        this.recipient,
        this.token.address,
        this.amount
      ).should.be.fulfilled;

      receipt = await this.cosmosBridge.connect(userFour).newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.symbol,
        this.senderSequence,
        this.recipient,
        this.token.address,
        this.amount
      ).should.be.fulfilled;

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await this.token.balanceOf(this.recipient)
      );
      afterUserBalance.should.be.bignumber.equal(this.amount);
    });

    it("should allow users to lock Ethereum in the bridge bank", async function () {
      const tx = await this.bridgeBank.connect(userOne).lock(
        this.sender,
        this.ethereumToken,
        this.weiAmount, {
          value: this.weiAmount
        }
      ).should.be.fulfilled;
      await tx.wait();

      const contractBalanceWei = await getBalance(this.bridgeBank.address);
      const contractBalance = Web3Utils.fromWei(contractBalanceWei, "ether");

      contractBalance.should.be.bignumber.equal(
        Web3Utils.fromWei(this.weiAmount, "ether")
      );
    });

    it("should not allow users to lock Ethereum in the bridge bank if the sent amount and amount param are different", async function () {
      await expect(
        this.bridgeBank.connect(userOne).lock(
          this.sender,
          this.ethereumToken,
          this.weiAmount + 1, {
            value: this.weiAmount
          },
        ),
      ).to.be.revertedWith("amount mismatch");
    });

    it("should not allow users to lock Ethereum in the bridge bank if sending tokens", async function () {
      await expect(
        this.bridgeBank.connect(userOne).lock(
          this.sender,
          this.token.address,
          this.weiAmount + 1, {
            value: this.weiAmount
          },
        ),
      ).to.be.revertedWith("do not send currency if locking tokens");
    });
  });
});
