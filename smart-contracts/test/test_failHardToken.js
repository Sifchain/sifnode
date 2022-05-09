const web3 = require("web3");
const BigNumber = web3.BigNumber;
const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { setup, getValidClaim } = require("./helpers/testFixture");

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

describe("Fail Hard Token", function () {
  let userOne;
  let userTwo;
  let addresses;
  let state;
  let failFactory;
  let failHardToken;

  before(async function () {
    accounts = await ethers.getSigners();
    addresses = accounts.map((e) => {
      return e.address;
    });

    failFactory = await ethers.getContractFactory("FailHardToken");

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

    failHardToken = await failFactory.deploy(
      "Fail Hard Token",
      "FAIL",
      userOne.address,
      state.amount
    );

    await failHardToken.deployed();
  });

  describe("Burn, Unlock", function () {
    it("should successfully process an unlock claim when token reverts transfer(), but user does not receive any tokens back", async function () {
      // Get balances before locking
      const beforeBridgeBankBalance = Number(
        await failHardToken.balanceOf(state.bridgeBank.address)
      );
      beforeBridgeBankBalance.should.be.equal(0);

      const beforeUserBalance = Number(await failHardToken.balanceOf(userOne.address));
      beforeUserBalance.should.be.equal(state.amount);

      // Attempt to lock tokens (will work without a previous approval - only for FailHardToken)
      await state.bridgeBank
        .connect(userOne)
        .lock(state.sender, failHardToken.address, state.amount);

      // Confirm that the tokens have left the user's wallet
      const afterUserBalance = Number(await failHardToken.balanceOf(userOne.address));
      afterUserBalance.should.be.bignumber.equal(0);

      // Confirm that bridgeBank now owns the tokens:
      const afterBridgeBankBalance = Number(
        await failHardToken.balanceOf(state.bridgeBank.address)
      );
      afterBridgeBankBalance.should.be.bignumber.equal(state.amount);

      // Now, try to unlock those tokens...

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: userOne.address,
        tokenAddress: failHardToken.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: "Fail Hard Token",
        tokenSymbol: "FAIL",
        tokenDecimals: 18,
        cosmosDenom: state.constants.denom.none,
        validators: accounts.slice(2, 6),
      });

      await state.cosmosBridge
        .connect(accounts[2])
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      // The prophecy should be completed without reverting, but the user shouldn't receive anything
      balance = Number(await failHardToken.balanceOf(userOne.address));
      expect(balance).to.be.equal(0);
    });
  });

  it("should revert on burn if user is blocklisted", async function () {
    await state.bridgeBank.connect(state.owner).addExistingBridgeToken(failHardToken.address);

    beforeUserBalance = Number(await failHardToken.balanceOf(userOne.address));

    // Lock & burn tokens
    const tx = await expect(
      state.bridgeBank
        .connect(userOne)
        .multiLockBurn([state.sender], [failHardToken.address], [state.amount], [true])
    ).to.be.rejected;

    // assert that there was no change
    afterUserBalance = Number(await failHardToken.balanceOf(userOne.address));
    afterUserBalance.should.be.equal(beforeUserBalance);
  });
});
