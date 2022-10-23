const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

// Bytes32 representation of Roles, according to OpenZeppelin's docs
const MINTER_ROLE = web3.utils.soliditySha3("MINTER_ROLE");
const DEFAULT_ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000";

describe("Test Bridge Token", function () {
  let userOne;
  let userTwo;
  let accounts;
  let owner;
  let bridgeTokenFactory;
  let bridgeToken;

  const name = "Test Bridge Token";
  const symbol = "TST";
  const decimals = 6;
  const denom = "ibc51b91cb1c1b98e88e4651a654b6541a65464846e6565b161651bb4aa84c654dd";
  const anotherDenom = "sif789de8f7997bd47c4a0928a001e916b5c68f1f33fef33d6588b868b93b6dcde6";

  before(async function () {
    accounts = await ethers.getSigners();

    bridgeTokenFactory = await ethers.getContractFactory("BridgeToken");

    owner = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
  });

  beforeEach(async function () {
    bridgeToken = await bridgeTokenFactory.deploy(name, symbol, decimals, denom);
    await bridgeToken.deployed();
  });

  it("should deploy and assign the correct values to variables", async function () {
    const _name = await bridgeToken.name();
    const _symbol = await bridgeToken.symbol();
    const _decimals = await bridgeToken.decimals();
    const _denom = await bridgeToken.cosmosDenom();
    const isAdmin = await bridgeToken.hasRole(DEFAULT_ADMIN_ROLE, owner.address);
    const isMinter = await bridgeToken.hasRole(MINTER_ROLE, owner.address);

    expect(_name).to.be.equal(name);
    expect(_symbol).to.be.equal(symbol);
    expect(_decimals).to.be.equal(decimals);
    expect(_denom).to.be.equal(denom);
    expect(isAdmin).to.be.true;
    expect(isMinter).to.be.true;
  });

  it("should allow owner to add a new minter", async function () {
    await expect(bridgeToken.connect(owner).grantRole(MINTER_ROLE, userOne.address))
      .to.emit(bridgeToken, "RoleGranted")
      .withArgs(MINTER_ROLE, userOne.address, owner.address);

    // check if the user received the minter role
    const isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;
  });

  it("should allow a minter to mint ERC20 tokens", async function () {
    // Add a new minter
    await bridgeToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    const isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // User should have no tokens yet
    let userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    const amount = 1000000;
    await bridgeToken.connect(userOne).mint(userOne.address, amount);

    // check if the user received the minted tokens
    userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(amount);
  });

  it("should NOT allow a non-minter user to mint ERC20 tokens", async function () {
    // User should have no tokens yet
    let userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Try to mint some tokens (should fail)
    const amount = 1000000;
    await expect(bridgeToken.connect(userOne).mint(userOne.address, amount)).to.be.revertedWith(
      `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${MINTER_ROLE}`
    );

    // check if the user received the minted tokens (should not have)
    userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);
  });

  it("should NOT allow a user to add a new minter", async function () {
    // Add a new minter
    await expect(
      bridgeToken.connect(userOne).grantRole(MINTER_ROLE, userOne.address)
    ).to.be.revertedWith(
      `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${DEFAULT_ADMIN_ROLE}`
    );

    // check if the user received the minter role
    isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.false;
  });

  it("should allow a new minter to mint tokens", async function () {
    // Add a new minter
    await bridgeToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    const isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // User should have no tokens yet
    let userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Mint some tokens
    const amount = 1000000;
    await bridgeToken.connect(userOne).mint(userOne.address, amount);

    // check if the user received the minted tokens
    userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(amount);
  });

  it("should allow owner to revoke minter role", async function () {
    // Add a new minter
    await bridgeToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    let isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Revoke minter role
    await expect(bridgeToken.connect(owner).revokeRole(MINTER_ROLE, userOne.address))
      .to.emit(bridgeToken, "RoleRevoked")
      .withArgs(MINTER_ROLE, userOne.address, owner.address);

    // check if the user lost the minter role
    isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.false;

    // User should have no tokens yet
    let userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);

    // Try to mint some tokens (should fail)
    const amount = 1000000;
    await expect(bridgeToken.connect(userOne).mint(userOne.address, amount)).to.be.revertedWith(
      `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${MINTER_ROLE}`
    );

    // check if the user received the minted tokens (should not have)
    userBalance = Number(await bridgeToken.balanceOf(userOne.address));
    userBalance.should.be.bignumber.equal(0);
  });

  it("should NOT allow a user to revoke minter role", async function () {
    // Add a new minter
    await bridgeToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    let isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Try to revoke minter role (should fail)
    await expect(
      bridgeToken.connect(userOne).revokeRole(MINTER_ROLE, userOne.address)
    ).to.be.revertedWith(
      `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${DEFAULT_ADMIN_ROLE}`
    );

    // check if the user kept the minter role
    isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;
  });

  it("should allow a minter to renounce it's own minter role", async function () {
    // Add a new minter
    await bridgeToken.connect(owner).grantRole(MINTER_ROLE, userOne.address);

    // check if the user received the minter role
    let isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.true;

    // Renounces the minter role
    await expect(bridgeToken.connect(userOne).renounceRole(MINTER_ROLE, userOne.address))
      .to.emit(bridgeToken, "RoleRevoked")
      .withArgs(MINTER_ROLE, userOne.address, userOne.address);

    // check if the user lost the minter role
    isMinter = await bridgeToken.hasRole(MINTER_ROLE, userOne.address);
    expect(isMinter).to.be.false;
  });

  it("should allow admin to transfer adminship of roles", async function () {
    // Grants the Admin role to userOne
    await expect(bridgeToken.connect(owner).grantRole(DEFAULT_ADMIN_ROLE, userOne.address))
      .to.emit(bridgeToken, "RoleGranted")
      .withArgs(DEFAULT_ADMIN_ROLE, userOne.address, owner.address);

    // check if the user received the minter role
    let hasAdminRole = await bridgeToken.hasRole(DEFAULT_ADMIN_ROLE, userOne.address);
    expect(hasAdminRole).to.be.true;

    // Onwer renounces admin role
    await expect(bridgeToken.connect(owner).renounceRole(DEFAULT_ADMIN_ROLE, owner.address))
      .to.emit(bridgeToken, "RoleRevoked")
      .withArgs(DEFAULT_ADMIN_ROLE, owner.address, owner.address);

    // check if owner lost the admin role
    hasAdminRole = await bridgeToken.hasRole(DEFAULT_ADMIN_ROLE, owner.address);
    expect(hasAdminRole).to.be.false;

    // Owner now tries to manage roles (should fail)
    await expect(
      bridgeToken.connect(owner).grantRole(DEFAULT_ADMIN_ROLE, owner.address)
    ).to.be.revertedWith(
      `AccessControl: account ${owner.address.toLowerCase()} is missing role ${DEFAULT_ADMIN_ROLE}`
    );

    // check if the owner received the admin role (should not have)
    hasAdminRole = await bridgeToken.hasRole(DEFAULT_ADMIN_ROLE, owner.address);
    expect(hasAdminRole).to.be.false;

    // guarantees userTwo has no minter rights
    let hasMinterRole = await bridgeToken.hasRole(MINTER_ROLE, userTwo.address);
    expect(hasMinterRole).to.be.false;

    // new admin grants the minter role to userTwo
    await expect(bridgeToken.connect(userOne).grantRole(MINTER_ROLE, userTwo.address))
      .to.emit(bridgeToken, "RoleGranted")
      .withArgs(MINTER_ROLE, userTwo.address, userOne.address);

    // check if owner received the minter role
    hasMinterRole = await bridgeToken.hasRole(MINTER_ROLE, userTwo.address);
    expect(hasMinterRole).to.be.true;

    // new admin changes the cosmosDenom:
    await expect(bridgeToken.connect(userOne).setDenom(anotherDenom)).to.be.fulfilled;

    // check if the denom changed
    const newDenom = await bridgeToken.cosmosDenom();
    expect(newDenom).to.be.equal(anotherDenom);
  });

  it("should allow owner to set the cosmosDenom", async function () {
    await expect(bridgeToken.connect(owner).setDenom(anotherDenom)).to.be.fulfilled;

    // check if the denom changed
    const newDenom = await bridgeToken.cosmosDenom();
    expect(newDenom).to.be.equal(anotherDenom);
  });

  it("should NOT allow a user to set the cosmosDenom", async function () {
    await expect(bridgeToken.connect(userOne).setDenom(anotherDenom)).to.be.revertedWith(
      `AccessControl: account ${userOne.address.toLowerCase()} is missing role ${DEFAULT_ADMIN_ROLE}`
    );
  });
});
