const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { setup, getValidClaim } = require("./helpers/testFixture");

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

describe("Test Cosmos Bridge", function () {
  let userOne;
  let userTwo;
  let userThree;
  let userFour;
  let accounts;
  let signerAccounts;
  let operator;
  let owner;
  let pauser;
  let unpauser;
  const consensusThreshold = 75;
  let initialPowers;
  let initialValidators;
  let networkDescriptor;
  let state;

  before(async function () {
    accounts = await ethers.getSigners();

    signerAccounts = accounts.map((e) => {
      return e.address;
    });

    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[9];

    owner = accounts[5];
    pauser = accounts[6];
    unpauser = accounts[8];

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);

    networkDescriptor = 1;
  });

  beforeEach(async function () {
    state = await setup({
      initialValidators,
      initialPowers,
      operator,
      consensusThreshold,
      owner,
      user: userOne,
      recipient: userThree,
      pauser,
      unpauser,
      networkDescriptor,
      lockTokensOnBridgeBank: true,
    });
  });


  it("should deploy a new token upon the successful processing of a ibc token burn prophecy claim just for first time", async function () {
    state.nonce = 1;

    const { digest, claimData, signatures } = await getValidClaim({
      sender: state.sender,
      senderSequence: state.senderSequence,
      recipientAddress: state.recipient.address,
      tokenAddress: ZERO_ADDRESS,
      amount: state.amount,
      bridgeToken: true,
      nonce: state.nonce,
      networkDescriptor: state.networkDescriptor,
      tokenName: state.name,
      tokenSymbol: state.symbol,
      tokenDecimals: state.decimals,
      cosmosDenom: state.constants.denom.ibc,
      validators: [userOne, userTwo, userFour],
    });

    const expectedAddress = ethers.utils.getContractAddress({
      from: state.bridgeBank.address,
      nonce: 1,
    });
    await expect(
      state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
    )
      .to.emit(state.cosmosBridge, "LogNewBridgeTokenCreated")
      .withArgs(
        state.decimals,
        state.networkDescriptor,
        state.name,
        state.symbol,
        ZERO_ADDRESS,
        expectedAddress,
        state.constants.denom.ibc
      );

    const newlyCreatedTokenAddress = await state.cosmosBridge.cosmosDenomToDestinationAddress(
      state.constants.denom.ibc
    );
    expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

    // Everything again, but this time submitProphecyClaimAggregatedSigs should NOT emit the event
    const {
      digest: digest2,
      claimData: claimData2,
      signatures: signatures2,
    } = await getValidClaim({
      sender: state.sender,
      senderSequence: state.senderSequence + 1,
      recipientAddress: state.recipient.address,
      tokenAddress: ZERO_ADDRESS,
      amount: state.amount,
      bridgeToken: true,
      nonce: state.nonce + 1,
      networkDescriptor: state.networkDescriptor,
      tokenName: state.name,
      tokenSymbol: state.symbol,
      tokenDecimals: state.decimals,
      cosmosDenom: state.constants.denom.ibc,
      validators: [userOne, userTwo, userFour],
    });

    await expect(
      state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest2, claimData2, signatures2)
    ).to.not.emit(state.cosmosBridge, "LogNewBridgeTokenCreated");

    // But should have minted:
    const deployedTokenFactory = await ethers.getContractFactory("BridgeToken");
    const deployedToken = await deployedTokenFactory.attach(newlyCreatedTokenAddress);
    const endingBalance = await deployedToken.balanceOf(state.recipient.address);
    expect(endingBalance).to.be.equal(state.amount * 2);
  });
});
