const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { ROWAN_DENOM } = require("./helpers/denoms");

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

async function setAllowance(user, erowan, rowan) {
  // Makes sure user has Erowans to migrate
  const erowanBalance = await erowan.balanceOf(user.address);
  expect(erowanBalance).to.not.be.equal(0);

  // Makes sure user has no Rowans yet:
  const rowanBalance = await rowan.balanceOf(user.address);
  expect(rowanBalance).to.be.equal(0);

  // Provides allowance
  await expect(erowan.connect(user).approve(rowan.address, erowanBalance)).to.be.fulfilled;

  // Make sure the allowance is set
  const allowance = await erowan.allowance(user.address, rowan.address);
  expect(allowance).to.be.equal(erowanBalance);

  return { erowanBalance, rowanBalance };
}

describe("Test Erowan migration", function () {
  let accounts;
  let userOne;
  let owner;
  let rowanTokenFactory;
  let rowanToken;
  let erowanTokenFactory;
  let erowanToken;

  const state = {
    erowan: {
      name: "SifChain",
      symbol: "erowan",
      decimals: 18,
      denom: "",
    },
    rowan: {
      name: "Rowan",
      symbol: "Rowan",
      decimals: 18,
      denom: ROWAN_DENOM,
    },
    amountToMint: 1000000,
  };

  before(async function () {
    accounts = await ethers.getSigners();

    rowanTokenFactory = await ethers.getContractFactory("Rowan");
    erowanTokenFactory = await ethers.getContractFactory("Erowan");

    owner = accounts[0];
    userOne = accounts[1];
  });

  beforeEach(async function () {
    // Deploy the old Erowan token
    erowanToken = await erowanTokenFactory.deploy(state.erowan.symbol);
    await erowanToken.deployed();

    // Deploy the new Rowan token
    rowanToken = await rowanTokenFactory.deploy(
      state.rowan.name,
      state.rowan.symbol,
      state.rowan.decimals,
      state.rowan.denom,
      erowanToken.address
    );
    await rowanToken.deployed();

    // Mint Erowans to userOne
    await expect(erowanToken.mint(userOne.address, state.amountToMint)).to.be.fulfilled;
  });

  it("should allow a user to migrate their Erowans to the new Rowan token after a correct allowance", async function () {
    let { erowanBalance } = await setAllowance(userOne, erowanToken, rowanToken);

    // Calls the migrate function on Rowan
    await expect(rowanToken.connect(userOne).migrate())
      .to.emit(rowanToken, "MigrationComplete")
      .withArgs(userOne.address, state.amountToMint);

    // Check if userOne received the tokens
    const rowanBalance = await rowanToken.balanceOf(userOne.address);
    expect(rowanBalance).to.be.equal(erowanBalance);

    // Check if userOne has no  more erowans
    erowanBalance = await erowanToken.balanceOf(userOne.address);
    expect(erowanBalance).to.be.equal(0);
  });

  it("should NOT allow a user to migrate their Erowans to the new Rowan token without allowance", async function () {
    // Calls the migrate function on Rowan
    await expect(rowanToken.connect(userOne).migrate()).to.be.rejectedWith(
      "ERC20: burn amount exceeds allowance"
    );

    // Check if userOne received the tokens (should not have)
    const rowanBalance = await rowanToken.balanceOf(userOne.address);
    expect(rowanBalance).to.be.equal(0);

    // Check if userOne has no more erowans (should have)
    const erowanBalance = await erowanToken.balanceOf(userOne.address);
    expect(erowanBalance).to.not.be.equal(0);
  });

  it("should allow a user to migrate 0 Erowans to the new Rowan token, but receive 0 Rowans", async function () {
    const { erowanBalance: erowanBalanceBefore } = await setAllowance(
      userOne,
      erowanToken,
      rowanToken
    );

    // Calls the migrate function on Rowan
    await expect(rowanToken.connect(userOne).migrate())
      .to.emit(rowanToken, "MigrationComplete")
      .withArgs(userOne.address, state.amountToMint);

    // Check if userOne received the tokens
    let rowanBalance = await rowanToken.balanceOf(userOne.address);
    expect(rowanBalance).to.be.equal(erowanBalanceBefore);

    // Check if userOne has no  more erowans
    let erowanBalance = await erowanToken.balanceOf(userOne.address);
    expect(erowanBalance).to.be.equal(0);

    // Calls the migrate function on Rowan AGAIN (event should inform that 0 tokens have been migrated)
    await expect(rowanToken.connect(userOne).migrate())
      .to.emit(rowanToken, "MigrationComplete")
      .withArgs(userOne.address, 0);

    // Check if userOne's Rowan balance remains the same
    rowanBalance = await rowanToken.balanceOf(userOne.address);
    expect(rowanBalance).to.be.equal(erowanBalanceBefore);

    // Check if userOne's Erowan balance remains the same
    erowanBalance = await erowanToken.balanceOf(userOne.address);
    expect(erowanBalance).to.be.equal(0);
  });
});
