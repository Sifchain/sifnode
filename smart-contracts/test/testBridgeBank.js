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

  describe("BridgeBank single lock burn transactions", function () {
    it("should allow user to lock ERC20 tokens", async function () {
      // approve and lock tokens
      await this.token.connect(userOne).approve(
        this.bridgeBank.address,
        this.amount
      );

      // Attempt to lock tokens
      await this.bridgeBank.connect(userOne).lock(
        this.sender,
        this.token.address,
        this.amount
      );

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await this.token.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(0);
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
        Web3Utils.fromWei((+this.weiAmount + +this.amount).toString(), "ether")
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
  describe("BridgeBank single lock burn transactions", function () {
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
