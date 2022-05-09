const web3 = require("web3");
const BigNumber = web3.BigNumber;
const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { setup } = require("./helpers/testFixture");

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

describe("Many Decimals Token", function () {
  let userOne;
  let userTwo;
  let addresses;
  let state;
  let tokenFactory;
  let manyDecimalsToken;

  before(async function () {
    accounts = await ethers.getSigners();
    addresses = accounts.map((e) => {
      return e.address;
    });

    tokenFactory = await ethers.getContractFactory("ManyDecimalsToken");

    userOne = accounts[6];
    userTwo = accounts[7];
  });

  beforeEach(async function () {
    state = await setup({
      initialValidators: addresses.slice(2, 6),
      initialPowers: [25, 25, 25, 25],
      operator: accounts[0],
      consensusThreshold: 75,
      owner: accounts[1],
      user: userOne,
      recipient: userTwo,
      pauser: accounts[8],
      unpauser: accounts[9],
      networkDescriptor: 1,
    });
    state.amount = 1000;

    manyDecimalsToken = await tokenFactory.deploy(
      "Many Decimals Token",
      "DEC",
      userOne.address,
      state.amount
    );

    await manyDecimalsToken.deployed();
  });

  describe("Lock", function () {
    it("should allow user to lock a token that has 100 decimals, but not twice!", async function () {
      // approve bridgebank to spend rowan
      await manyDecimalsToken.connect(userOne).approve(state.bridgeBank.address, state.amount);

      // Lock & burn tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn([state.sender], [manyDecimalsToken.address], [state.amount / 2], [false])
      )
        .to.emit(state.bridgeBank, "LogLock")
        .withArgs(
          userOne.address,
          state.sender,
          manyDecimalsToken.address,
          state.amount / 2,
          1,
          255,
          "DEC",
          "Many Decimals Token",
          1
        );

      afterUserBalance = Number(await manyDecimalsToken.balanceOf(userOne.address));
      afterUserBalance.should.be.equal(state.amount / 2);

      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn([state.sender], [manyDecimalsToken.address], [state.amount / 2], [false])
      )
        .to.emit(state.bridgeBank, "LogLock")
        .withArgs(
          userOne.address,
          state.sender,
          manyDecimalsToken.address,
          state.amount / 2,
          2,
          255,
          "DEC",
          "Many Decimals Token",
          1
        );
    });
  });
});
