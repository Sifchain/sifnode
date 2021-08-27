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
    const isMinter = await ibcToken.minters(owner.address);

    expect(_name).to.be.equal(name);
    expect(_symbol).to.be.equal(symbol);
    expect(_decimals).to.be.equal(decimals);
    expect(_denom).to.be.equal(denom);
    expect(isMinter).to.be.true;
  });

  it("should allow owner (who is a minter) to mint ERC20 tokens", async function () {
    const amount = 1000000;

    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    await ibcToken.connect(owner).mint(userOne.address, amount);

    // check if the user received the minted tokens
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(amount);
  });

  it("should NOT allow a non-minter user to mint ERC20 tokens", async function () {
    const amount = 1000000;

    // Try to mint some tokens
    await expect(ibcToken.connect(userOne).mint(userOne.address, amount))
      .to.be.revertedWith("MinterRole: caller does not have the Minter role");

    // check if the user received the minted tokens
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);
  });

  it("should allow owner to add a new minter", async function () {
    // Add a new minter
    await ibcToken.connect(owner).addMinter(userOne.address);

    // check if the user received the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.true;
  });

  it("should emit en event when the owner adds a new minter", async function () {
    // Add a new minter
    await expect(ibcToken.connect(owner).addMinter(userOne.address))
      .to.emit(ibcToken, 'MinterUpdate')
      .withArgs(userOne.address, true);

    // check if the user received the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.true;
  });

  it("should NOT allow a user to add a new minter", async function () {
    // Add a new minter
    await expect(ibcToken.connect(userOne).addMinter(userOne.address))
      .to.be.revertedWith('Ownable: caller is not the owner');

    // check if the user received the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.false;
  });

  it("should allow a new minter to mint tokens", async function () {
    // Add a new minter
    await ibcToken.connect(owner).addMinter(userOne.address);

    // check if the user received the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.true;

    const amount = 1000000;

    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    await ibcToken.connect(userOne).mint(userOne.address, amount);

    // check if the user received the minted tokens
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(amount);
  });

  it("should allow owner to remove a minter", async function () {
    // Add a new minter
    await ibcToken.connect(owner).addMinter(userOne.address);

    // check if the user received the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.true;

    // Remove the new minter
    await ibcToken.connect(owner).removeMinter(userOne.address);

    // check if the user lost the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.false;

    // Try to mint tokens (should fail)
    const amount = 1000000;

    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    await expect(ibcToken.connect(userOne).mint(userOne.address, amount))
      .to.be.revertedWith("MinterRole: caller does not have the Minter role");

    // check if the user received the minted tokens (should not have)
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);
  });

  it("should emit en event when owner removes a minter", async function () {
    // Add a new minter
    await ibcToken.connect(owner).addMinter(userOne.address);

    // check if the user received the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.true;

    // Remove the new minter
    await expect(ibcToken.connect(owner).removeMinter(userOne.address))
      .to.emit(ibcToken, 'MinterUpdate')
      .withArgs(userOne.address, false);

    // check if the user lost the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.false;
  });

  it("should NOT allow a user to remove a minter", async function () {
    // Add a new minter
    await ibcToken.connect(owner).addMinter(userOne.address);

    // check if the user received the minter role
    isMinter = await ibcToken.minters(userOne.address);
    expect(isMinter).to.be.true;

    // Try to remove the new minter (should fail)
    await expect(ibcToken.connect(userOne).removeMinter(owner.address))
      .to.be.revertedWith("Ownable: caller is not the owner")

    // check if the owner kept the minter role
    isMinter = await ibcToken.minters(owner.address);
    expect(isMinter).to.be.true;
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