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

// Bytes32 representation of Roles, according to OpenZeppelin's docs
const MINTER_ROLE = web3.utils.soliditySha3('MINTER_ROLE');
const DEFAULT_ADMIN_ROLE = '0x0000000000000000000000000000000000000000000000000000000000000000';

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
    const isAdmin = await ibcToken.hasRole(DEFAULT_ADMIN_ROLE, owner.address);
    const isMinter = await ibcToken.hasRole(MINTER_ROLE, owner.address);

    expect(_name).to.be.equal(name);
    expect(_symbol).to.be.equal(symbol);
    expect(_decimals).to.be.equal(decimals);
    expect(_denom).to.be.equal(denom);
    expect(isAdmin).to.be.true;
    expect(isMinter).to.be.false;
  });

  it("should allow owner to add a new minter", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    const isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;
  });

  it("should allow a minter to mint ERC20 tokens", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    const isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    const amount = 1000000;
    await ibcToken.connect(userOne).mint(userOne.address, amount);

    // check if the user received the minted tokens
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(amount);
  });

  it("should NOT allow a non-minter user to mint ERC20 tokens", async function () {
    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Try to mint some tokens (should fail)
    const amount = 1000000;
    await expect(ibcToken.connect(userOne).mint(userOne.address, amount))
      .to.be.revertedWith(
        `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${MINTER_ROLE}`
      );

    // check if the user received the minted tokens (should not have)
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);
  });

  it("should emit an event when the owner adds a new minter", async function () {
    // Add a new minter
    await expect(ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address))
      .to.emit(ibcToken, 'RoleGranted')
      .withArgs(MINTER_ROLE, userOne.address, owner.address);

    // check if the user received the minter role
    const isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;
  });

  it("should NOT allow a user to add a new minter", async function () {
    // Add a new minter
    await expect(ibcToken.connect(userOne).grantRole(MINTER_ROLE, userOne.address))
      .to.be.revertedWith(
        `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${DEFAULT_ADMIN_ROLE}`
      );

    // check if the user received the minter role
    isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.false;
  });

  it("should allow a new minter to mint tokens", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    const isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    const amount = 1000000;
    await ibcToken.connect(userOne).mint(userOne.address, amount);

    // check if the user received the minted tokens
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(amount);
  });

  it("should allow owner to revoke minter role", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    let isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Remove the new minter
    await ibcToken.connect(owner).revokeRole(MINTER_ROLE, userOne.address);

    // check if the user lost the minter role
    isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.false;

    // User should have no tokens yet
    let userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Try to mint some tokens (should fail)
    const amount = 1000000;
    await expect(ibcToken.connect(userOne).mint(userOne.address, amount))
      .to.be.revertedWith(
        `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${MINTER_ROLE}`
      );

    // check if the user received the minted tokens (should not have)
    userBalance = Number(await ibcToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);
  });

  it("should emit en event when owner revokes minter role", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);
      
    // check if the user received the minter role
    let isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Revoke minter role
    await expect(ibcToken.connect(owner).revokeRole(MINTER_ROLE, userOne.address))
      .to.emit(ibcToken, 'RoleRevoked')
      .withArgs(MINTER_ROLE, userOne.address, owner.address);

    // check if the user received the minter role
    isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.false;
  });

  it("should NOT allow a user to revoke minter role", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);
      
    // check if the user received the minter role
    let isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Try to revoke minter role (should fail)
    await expect(ibcToken.connect(userOne).revokeRole(MINTER_ROLE, userOne.address))
      .to.be.revertedWith(
        `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${DEFAULT_ADMIN_ROLE}`
      );

    // check if the user kept the minter role
    isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;
  });

  it("should allow a minter to renounce it's own minter role", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);
      
    // check if the user received the minter role
    let isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Renounces the minter role
    await ibcToken.connect(userOne).renounceRole(MINTER_ROLE, userOne.address);

    // check if the user lost the minter role
    isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.false;
  });

  it("should emit an event when a minter renounces the minter role", async function () {
    // Add a new minter
    await ibcToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);
      
    // check if the user received the minter role
    let isMinter = await ibcToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Renounces the minter role
    await expect(ibcToken.connect(userOne).renounceRole(MINTER_ROLE, userOne.address))
      .to.emit(ibcToken, 'RoleRevoked')
      .withArgs(MINTER_ROLE, userOne.address, userOne.address);
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