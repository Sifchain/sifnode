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

    this.token1 = await BridgeToken.deploy(
      this.name,
      this.symbol,
      18
    );

    this.token2 = await BridgeToken.deploy(
      this.name,
      this.symbol,
      18
    );

    this.token3 = await BridgeToken.deploy(
      this.name,
      this.symbol,
      18
    );

    await this.token1.deployed();
    await this.token2.deployed();
    await this.token3.deployed();

    this.amount = 100;
    //Load user account with ERC20 tokens for testing
    await this.token1.mint(userOne.address, this.amount * 2, {
      from: operator
    }).should.be.fulfilled;

    await this.token2.mint(userOne.address, this.amount * 2, {
      from: operator
    }).should.be.fulfilled;

    await this.token3.mint(userOne.address, this.amount * 2, {
      from: operator
    }).should.be.fulfilled;
  });

  describe("BridgeBank", function () {
    it("should allow user to lock ERC20 tokens", async function () {
      await this.token1.connect(userOne).approve(
        this.bridgeBank.address,
        this.amount
      );

      // Attempt to lock tokens
      await this.bridgeBank.connect(userOne).lock(
        this.sender,
        this.token1.address,
        this.amount
      );

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await this.token1.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(this.amount);
    });

    it("should allow user to multi-lock ERC20 tokens", async function () {
      await this.token1.connect(userOne).approve(
        this.bridgeBank.address,
        this.amount
      );

      await this.token2.connect(userOne).approve(
        this.bridgeBank.address,
        this.amount
      );

      await this.token3.connect(userOne).approve(
        this.bridgeBank.address,
        this.amount
      );

      // Attempt to lock tokens
      await this.bridgeBank.connect(userOne).multiLock(
        [this.sender, this.sender, this.sender],
        [this.token1.address,this.token2.address,this.token3.address],
        [this.amount, this.amount, this.amount]
      );

      // Confirm that the user has been minted the correct token
      let afterUserBalance = Number(
        await this.token1.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(this.amount);

      afterUserBalance = Number(
        await this.token2.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(this.amount);

      afterUserBalance = Number(
        await this.token3.balanceOf(userOne.address)
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
        Web3Utils.fromWei((+this.weiAmount).toString(), "ether")
      );
    });
  });
});
