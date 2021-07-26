const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const {
  singleSetup,
  getDigestNewProphecyClaim,
  signHash
} = require("./helpers/testFixture");

require("chai")
  .use(require("chai-as-promised"))
  .should();

use(solidity);

describe("Test Chain Id Mismatch", function () {
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
  let chainId;
  let state;

  before(async function() {
    accounts = await ethers.getSigners();

    signerAccounts = accounts.map((e) => { return e.address });

    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[9].address;

    owner = accounts[5];
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);

    chainId = 1;
  });

  beforeEach(async function () {
    state = await singleSetup(
      initialValidators,
      initialPowers,
      operator,
      consensusThreshold,
      owner,
      userOne,
      userThree,
      pauser,
      chainId,
      true // force chain id mismatch
    );
  });

  describe("CosmosBridge", function () {
    it("should not allow unlocking tokens upon the processing of a burn prophecy claim with the wrong chain id", async function () {
      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient,
        state.token.address,
        state.amount,
        false,
        state.nonce,
        state.chainId
      ]);

      const signatures = await signHash([userOne, userTwo, userFour], digest);
      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient,
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        chainId: state.chainId
      };

      await expect(state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        )).to.be.revertedWith("INV_CHAIN_ID");
    });

    it("should not allow unlocking native tokens upon the processing of a burn prophecy claim with the wrong chain id", async function () {
      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient,
        state.ethereumToken,
        state.amount,
        false,
        state.nonce,
        state.chainId
      ]);
      const signatures = await signHash([userOne, userTwo, userFour], digest);

      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient,
        tokenAddress: state.ethereumToken,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        chainId: state.chainId
      };

      await expect(state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        )).to.be.revertedWith("INV_CHAIN_ID");
    });
  });
});
