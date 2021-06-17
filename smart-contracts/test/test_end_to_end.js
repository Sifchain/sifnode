var bigInt = require("big-integer");
const web3 = require("web3");
const Web3Utils = require("web3-utils");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(web3.BigNumber))
  .should();

const { expect } = require('chai');
const { multiTokenSetup } = require('./helpers/testFixture');

const getBalance = async function(address) {
  return await network.provider.send("eth_getBalance", [address]);
}

describe("End To End", function () {
  let userOne;
  let userTwo;
  let userThree;
  let userFour;
  let userSeven;
  let accounts;
  let operator;
  let owner;
  let pauser;

  // Consensus threshold of 70%
  const consensusThreshold = 70;
  let initialPowers;
  let initialValidators;
  let state;

  before(async function() {
    accounts = await ethers.getSigners();
    
    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userThree = accounts[3];
    userFour = accounts[4];
    userSeven = accounts[7];

    owner = accounts[5];
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = [
      userOne.address,
      userTwo.address,
      userThree.address,
      userFour.address
    ];

      // Deploy Valset contract
    state = await multiTokenSetup(
      initialValidators,
      initialPowers,
      operator,
      consensusThreshold,
      owner,
      userOne,
      userThree,
      pauser
    );

    state.amountWei = 100;
    state.amountNativeCosmos = 815;
    state.secondValidators = [userOne.address, userTwo.address];
    state.secondPowers = [50, 50];
    state.thirdValidators = [userThree.address, userFour.address];
    state.thirdPowers = [50, 50];
    state.ethTokenAddress = state.ethereumToken;
  });

  describe("Claim flow", function () {
    it("Burn prophecy claim flow lock", async function () {
      console.log("\t[Attempt burn -> unlock]");

      const startContractBalanceWei = await getBalance(
        state.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(startContractBalanceWei).should.be.equal(0);

      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await state.bridgeBank.connect(userOne).lock(
        state.cosmosSender,
        state.ethTokenAddress,
        state.amountWei,
        {
          value: state.amountWei
        }
      ).should.be.fulfilled;

      const contractBalanceWei = await getBalance(
        state.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei).should.be.equal(state.amountWei);
    });

    it("New prophecy claim eth unlock", async function () {
      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      state.nonce = 1;
      await state.cosmosBridge.connect(userOne).newProphecyClaim(
        state.sender,
        state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        state.nonce
      ).should.be.fulfilled;

      await state.cosmosBridge.connect(userTwo).newProphecyClaim(
        state.sender,
        state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        state.nonce
      ).should.be.fulfilled;

      await state.cosmosBridge.connect(userFour).newProphecyClaim(
        state.sender,
        state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        state.nonce
      ).should.be.fulfilled;

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      let postRecipientBalance = Web3Utils.fromWei(await getBalance(userSeven.address));

      // assert user received their funds
      let expectedBalance = "10000.0000000000000001";
      expect(expectedBalance).to.be.equal(postRecipientBalance);
    });

    it("Fail to create new prophecy claim eth unlock", async function () {
      // Fail to create prophecy claim if from non validator
      await expect(
        state.cosmosBridge.connect(owner).newProphecyClaim(
          state.sender,
          state.senderSequence,
          state.recipient.address,
          state.token1.address,
          state.amount,
          false,
          state.nonce
        ),
      ).to.be.revertedWith("Must be an active validator");
    });

    it("Create new prophecy claim eth unlock after eth lock", async function () {
      // Also make sure everything runs twice.

      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await state.bridgeBank.connect(userOne).lock(
        state.cosmosSender,
        state.ethTokenAddress,
        state.amountWei,
        {
          value: state.amountWei
        }
      ).should.be.fulfilled;

      const contractBalanceWei2 = await getBalance(
          state.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei2).should.be.equal(state.amountWei);

      // if nonce is not incremented, things should revert
      await expect(
        state.cosmosBridge.connect(userOne).newProphecyClaim(
          state.sender,
          state.senderSequence + 1,
          userSeven.address,
          state.ethTokenAddress,
          state.amount,
          false,
          state.nonce
        ),
      ).to.be.revertedWith("INV_ORD");

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      await state.cosmosBridge.connect(userOne).newProphecyClaim(
        state.sender,
        ++state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        ++state.nonce
      ).should.be.fulfilled;

      await state.cosmosBridge.connect(userTwo).newProphecyClaim(
        state.sender,
        state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        state.nonce
      ).should.be.fulfilled;

      await state.cosmosBridge.connect(userFour).newProphecyClaim(
        state.sender,
        state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        state.nonce
      ).should.be.fulfilled;

      postRecipientBalance = Web3Utils.fromWei(await getBalance(userSeven.address));
      expectedBalance = "10000.0000000000000002";
      expect(postRecipientBalance).to.be.equal(expectedBalance);
    });

    it("Can update the valset", async function () {
      // Operator resets the valset
      await state.cosmosBridge.connect(operator).updateValset(
        state.secondValidators,
        state.secondPowers,
      ).should.be.fulfilled;

      // Confirm that both initial validators are now active validators
      const isUserOneValidator = await state.cosmosBridge.isActiveValidator(
        userOne.address
      );
      isUserOneValidator.should.be.equal(true);
      const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(
        userTwo.address
      );
      isUserTwoValidator.should.be.equal(true);

      // Confirm that all both secondary validators are not active validators
      const isUserThreeValidator = await state.cosmosBridge.isActiveValidator(
        userThree.address
      );
      isUserThreeValidator.should.be.equal(false);
      const isUserFourValidator = await state.cosmosBridge.isActiveValidator(
        userFour.address
      );
      isUserFourValidator.should.be.equal(false);
    });

    it("Create new prophecy claim eth unlock and lock with new validator set", async function () {
      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await state.bridgeBank.connect(userOne).lock(
        state.cosmosSender,
        state.ethTokenAddress,
        state.amountWei,
        {
          value: state.amountWei
        }
      ).should.be.fulfilled;

      const contractBalanceWei3 = await getBalance(
          state.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei3).should.be.equal(state.amountWei);

      // --------------------------------------------------------
      //  Check receiver's account balance prior to the claims
      // --------------------------------------------------------

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      await state.cosmosBridge.connect(userOne).newProphecyClaim(
        state.sender,
        ++state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        ++state.nonce
      ).should.be.fulfilled;

      await state.cosmosBridge.connect(userTwo).newProphecyClaim(
        state.sender,
        state.senderSequence,
        userSeven.address,
        state.ethTokenAddress,
        state.amount,
        false,
        state.nonce
      ).should.be.fulfilled;

      postRecipientBalance = Web3Utils.fromWei(await getBalance(userSeven.address));
      expectedBalance = "10000.0000000000000003";
      expect(postRecipientBalance).to.be.equal(expectedBalance);

      // Fail to create prophecy claim if from non validator
      await expect(
        state.cosmosBridge.connect(userFour).newProphecyClaim(
          state.sender,
          state.senderSequence,
          userSeven.address,
          state.ethTokenAddress,
          state.amount,
          false,
          state.nonce
        ),
      ).to.be.revertedWith("Must be an active validator");

      // Fail to create prophecy claim if out of order
      await expect(
        state.cosmosBridge.connect(userTwo).newProphecyClaim(
          state.sender,
          state.senderSequence + 10,
          userSeven.address,
          state.ethTokenAddress,
          state.amount,
          false,
          state.nonce + 10
        ),
      ).to.be.revertedWith("INV_ORD");
    });

    it("Can update the validator set", async function () {
      // Also make sure everything runs fourth time after switching validators a second time.
      // Operator resets the valset
      await state.cosmosBridge.connect(operator).updateValset(
          state.thirdValidators,
          state.thirdPowers,
      ).should.be.fulfilled;

      // Confirm that both initial validators are no longer an active validators
      const isUserOneValidator2 = await state.cosmosBridge.isActiveValidator(
          userOne.address
      );
      isUserOneValidator2.should.be.equal(false);
      const isUserTwoValidator2 = await state.cosmosBridge.isActiveValidator(
          userTwo.address
      );
      isUserTwoValidator2.should.be.equal(false);

      // Confirm that both secondary validators are now active validators
      const isUserThreeValidator2 = await state.cosmosBridge.isActiveValidator(
          userThree.address
      );
      isUserThreeValidator2.should.be.equal(true);
      const isUserFourValidator2 = await state.cosmosBridge.isActiveValidator(
          userFour.address
      );
      isUserFourValidator2.should.be.equal(true);
    });

    it("Create new prophecy claim eth unlock and lock with new validator set", async function () {
      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await state.bridgeBank.connect(userOne).lock(
        state.cosmosSender,
        state.ethTokenAddress,
        state.amountWei,
        {
          value: state.amountWei
        }
      ).should.be.fulfilled;

      const contractBalanceWei4 = await getBalance(
          state.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei4).should.be.equal(state.amountWei);

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      
      // should have two validators
      expect(2).to.be.equal(Number(await state.cosmosBridge.validatorCount()));

      let currentNonceUserThree = Number(await state.cosmosBridge.lastNonceSubmitted(userThree.address)) + 1;
      let currentNonceUserFour = Number(await state.cosmosBridge.lastNonceSubmitted(userFour.address)) + 1;

      await state.cosmosBridge.connect(userThree).newProphecyClaim(
        state.sender,
        ++state.senderSequence,
        state.recipient.address,
        state.ethTokenAddress,
        state.amount,
        false,
        currentNonceUserThree
      ).should.be.fulfilled;

      await state.cosmosBridge.connect(userFour).newProphecyClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.ethTokenAddress,
        state.amount,
        false,
        currentNonceUserFour
      ).should.be.fulfilled;

      // Fail to create prophecy claim if from non validator
      await expect(
        state.cosmosBridge.connect(userOne).newProphecyClaim(
          state.sender,
          state.senderSequence,
          state.recipient.address,
          state.ethTokenAddress,
          state.amount,
          false,
          state.nonce
        ),
      ).to.be.revertedWith("Must be an active validator");

      const contractBalanceWeiAfter = Web3Utils.fromWei(
        await getBalance(
          state.bridgeBank.address
        )
      );
      contractBalanceWeiAfter.toString().should.be.equal("0")
    });
  });
});
