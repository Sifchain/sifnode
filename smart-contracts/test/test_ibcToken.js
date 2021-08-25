const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

use(solidity);

describe("Test IBC Token", function () {
  let userOne;
  let accounts;
  let owner;
  let ibcTokenFactory;
  let ibcToken;

  const name = "Test Ibc Token";
  const symbol = "TST";
  const decimals = 6;
  const denom = "ibc51b91cb1c1b98e88e4651a654b6541a65464846e6565b161651bb4aa84c654dd";
  const anotherDenom = "sif789de8f7997bd47c4a0928a001e916b5c68f1f33fef33d6588b868b93b6dcde6";

  before(async function () {
    accounts = await ethers.getSigners();
  
    ibcTokenFactory = await ethers.getContractFactory("IbcToken");
    signerAccounts = accounts.map((e) => { return e.address });
  
    owner = accounts[0];
    userOne = accounts[1];
  });
  
  beforeEach(async function () {
    ibcToken = await ibcTokenFactory.deploy(
      name,
      symbol,
      decimals,
      denom
    );
    await ibcToken.deployed();
  });

  it("should deploy and assign the correct values to variables", async function () {
    const _name = await ibcToken.name();
    const _symbol = await ibcToken.symbol();
    const _decimals = await ibcToken.decimals();
    const _denom = await ibcToken.cosmosDenom();

    expect(_name).to.be.equal(name);
    expect(_symbol).to.be.equal(symbol);
    expect(_decimals).to.be.equal(decimals);
    expect(_denom).to.be.equal(denom);
  });

  it("should allow owner to mint ERC20 tokens", async function () {
    const amount = 1000000;

    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    await ibcToken.connect(owner).mint(userOne.address, amount)

    // check if the user received the minted tokens
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(amount);
  });

  it("should NOT allow a user to mint ERC20 tokens", async function () {
    const amount = 1000000;

    // Try to mint some tokens
    await expect(ibcToken.connect(userOne).mint(userOne.address, amount))
      .to.be.revertedWith("Ownable: caller is not the owner");

    // check if the user received the minted tokens
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);
  });

  it("should allow owner to set the cosmosDenom", async function () {
    await expect(ibcToken.connect(owner).setDenom(anotherDenom))
      .to.be.fulfilled;

    // check if the denom changed
    const newDenom = await ibcToken.cosmosDenom();
    expect(newDenom).to.be.equal(anotherDenom);
  });

  it("should NOT allow a user to set the cosmosDenom", async function () {
    await expect(ibcToken.connect(userOne).setDenom(anotherDenom))
      .to.be.revertedWith("Ownable: caller is not the owner");
  });
});