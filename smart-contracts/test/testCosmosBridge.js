const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { setup, getValidClaim, assertTokensMinted } = require("./helpers/testFixture");

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

const getBalance = async function (address, fromWei = false) {
  let balance = await network.provider.send("eth_getBalance", [address]);

  if (fromWei) {
    balance = Web3Utils.fromWei(balance);
  }

  return balance;
};

describe("Test Cosmos Bridge", function () {
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
      networkDescriptor,
      lockTokensOnBridgeBank: true,
    });
  });

  describe("CosmosBridge:Oracle", function () {
    it("should be able to getProphecyStatus with 1/3 of power", async function () {
      const status = await state.cosmosBridge.connect(userOne).getProphecyStatus(25);
      expect(status).to.equal(false);
    });

    it("should be able to getProphecyStatus with 2/3 of power", async function () {
      const status = await state.cosmosBridge.getProphecyStatus(50);
      expect(status).to.equal(false);
    });

    it("should be able to getProphecyStatus with 3/3 of power", async function () {
      const status = await state.cosmosBridge.getProphecyStatus(75);
      expect(status).to.equal(true);
    });

    it("should be able to getProphecyStatus with 100% of power", async function () {
      const status = await state.cosmosBridge.getProphecyStatus(100);
      expect(status).to.equal(true);
    });
  });

  describe("CosmosBridge", function () {
    it("Can update the valset", async function () {
      // Operator resets the valset
      await state.cosmosBridge
        .connect(operator)
        .updateValset([userOne.address, userTwo.address], [50, 50]).should.be.fulfilled;

      // Confirm that both initial validators are now active validators
      const isUserOneValidator = await state.cosmosBridge.isActiveValidator(userOne.address);
      isUserOneValidator.should.be.equal(true);
      const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(userTwo.address);
      isUserTwoValidator.should.be.equal(true);

      // Confirm that all both secondary validators are not active validators
      const isUserThreeValidator = await state.cosmosBridge.isActiveValidator(userThree.address);
      isUserThreeValidator.should.be.equal(false);
      const isUserFourValidator = await state.cosmosBridge.isActiveValidator(userFour.address);
      isUserFourValidator.should.be.equal(false);
    });

    it("Can change the operator", async function () {
      // Confirm that the operator has changed
      const originalOperator = await state.cosmosBridge.operator();
      expect(originalOperator).to.be.equal(operator.address);

      // Operator resets the valset
      await state.cosmosBridge.connect(operator).changeOperator(userOne.address).should.be
        .fulfilled;

      // Confirm that the operator has changed
      const newOperator = await state.cosmosBridge.operator();
      expect(newOperator).to.be.equal(userOne.address);
    });

    it("should NOT allow to change the operator to the zero address", async function () {
      // Confirm that the operator has changed
      const originalOperator = await state.cosmosBridge.operator();
      expect(originalOperator).to.be.equal(operator.address);

      // Operator resets the valset
      await expect(
        state.cosmosBridge.connect(operator).changeOperator(state.constants.zeroAddress)
      ).to.be.rejectedWith("invalid address");

      // Confirm that the operator has NOT changed
      const newOperator = await state.cosmosBridge.operator();
      expect(newOperator).to.be.equal(operator.address);
    });

    it("Can update the validator set", async function () {
      // Also make sure everything runs fourth time after switching validators a second time.
      // Operator resets the valset
      await state.cosmosBridge
        .connect(operator)
        .updateValset([userThree.address, userFour.address], [50, 50]).should.be.fulfilled;

      // Confirm that both initial validators are no longer an active validators
      const isUserOneValidator2 = await state.cosmosBridge.isActiveValidator(userOne.address);
      isUserOneValidator2.should.be.equal(false);
      const isUserTwoValidator2 = await state.cosmosBridge.isActiveValidator(userTwo.address);
      isUserTwoValidator2.should.be.equal(false);

      // Confirm that both secondary validators are now active validators
      const isUserThreeValidator2 = await state.cosmosBridge.isActiveValidator(userThree.address);
      isUserThreeValidator2.should.be.equal(true);
      const isUserFourValidator2 = await state.cosmosBridge.isActiveValidator(userFour.address);
      isUserFourValidator2.should.be.equal(true);
    });

    it("should return true if a sifchain address prefix is correct", async function () {
      (await state.bridgeBank.VSA(state.sender)).should.be.equal(true);
    });

    it("should return false if a sifchain address length is incorrect", async function () {
      const incorrectSifAddress = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpaceee"
      );
      (await state.bridgeBank.VSA(incorrectSifAddress)).should.be.equal(false);
    });

    it("should return false if a sifchain address has an incorrect `sif` prefix", async function () {
      const incorrectSifAddress = web3.utils.utf8ToHex(
        "eif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      (await state.bridgeBank.VSA(incorrectSifAddress)).should.be.equal(false);
    });

    it("Should deploy cosmos bridge and bridge bank", async function () {
      expect((await state.cosmosBridge.consensusThreshold()).toString()).to.equal(
        consensusThreshold.toString()
      );

      // iterate over all validators and ensure they have the proper
      // powers and that they have been succcessfully whitelisted
      for (let i = 0; i < initialValidators.length; i++) {
        const address = initialValidators[i];

        expect(await state.cosmosBridge.isActiveValidator(address)).to.be.true;

        expect((await state.cosmosBridge.getValidatorPower(address)).toString()).to.equal("25");
      }

      expect(await state.bridgeBank.cosmosBridge()).to.be.equal(state.cosmosBridge.address);
      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
    });

    it("Should deploy cosmos bridge and bridge bank, correctly setting the networkDescriptor", async function () {
      expect(await state.cosmosBridge.networkDescriptor()).to.equal(state.networkDescriptor);
      expect(await state.bridgeBank.networkDescriptor()).to.equal(state.networkDescriptor);
    });

    it("should unlock tokens upon the successful processing of a burn prophecy claim", async function () {
      const beforeUserBalance = Number(await state.token.balanceOf(state.recipient.address));
      beforeUserBalance.should.be.bignumber.equal(Number(0));

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: [userOne, userTwo, userFour],
      });

      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      balance = Number(await state.token.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(state.amount);
    });

    it("should unlock eth upon the successful processing of a burn prophecy claim", async function () {
      // assert recipient balance before receiving proceeds from newProphecyClaim is correct
      const startingBalance = await getBalance(state.recipient.address, true);
      expect(startingBalance).to.be.equal("10000");
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.constants.zeroAddress,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: "Ether",
        tokenSymbol: "ETH",
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.ether,
        validators: [userOne, userTwo, userFour],
      });

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      const endingBalance = await getBalance(state.recipient.address, true);
      expectedEndingBalance = "10000.0000000000000001"; // added 100 weis
      expect(endingBalance).to.be.equal(expectedEndingBalance);
    });

    it("should NOT unlock to a blocklisted address eth upon the processing of a burn prophecy claim", async function () {
      // Add recipient to the blocklist
      await expect(state.blocklist.connect(operator).addToBlocklist(state.recipient.address)).to.be
        .fulfilled;

      // assert recipient balance before receiving proceeds from newProphecyClaim is correct
      const startingBalance = await getBalance(state.recipient.address, true);
      expect(startingBalance).to.be.equal("10000.0000000000000001");
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.constants.zeroAddress,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: "Ether",
        tokenSymbol: "ETH",
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.ether,
        validators: [userOne, userTwo, userFour],
      });

      // get the prophecyID:
      const prophecyID = await state.cosmosBridge.getProphecyID(
        state.sender, // cosmosSender
        state.senderSequence, // cosmosSenderSequence
        state.recipient.address, // ethereumReceiver
        state.constants.zeroAddress, // tokenAddress
        state.amount, // amount
        "Ether", // tokenName
        "ETH", // tokenSymbol
        state.decimals, // tokenDecimals
        state.networkDescriptor, // networkDescriptor
        false, // doublePeg
        state.nonce, // nonce
        state.constants.denom.ether // cosmosDenom
      );

      // Doesn't revert, but user doesn't get the funds either
      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      )
        .to.emit(state.cosmosBridge, "LogProphecyCompleted")
        .withArgs(prophecyID, false); // the second argument here is 'success'

      // Make sure the balance didn't change
      const endingBalance = await getBalance(state.recipient.address, true);
      expect(endingBalance).to.be.equal(startingBalance);
    });

    it("should deploy a new token upon the successful processing of a double-pegged burn prophecy claim", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
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
          state.token.address,
          expectedAddress,
          state.constants.denom.one
        );

      // check if the new token has been correctly deployed
      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(
        state.token.address
      );
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

      // check if the user received minted tokens
      const deployedToken = await state.factories.BridgeToken.attach(newlyCreatedTokenAddress);
      const mintedTokensToRecipient = await deployedToken.balanceOf(state.recipient.address);
      const totalSupply = await deployedToken.totalSupply();
      expect(totalSupply).to.be.equal(state.amount);
      expect(mintedTokensToRecipient).to.be.equal(state.amount);
    });

    it("should deploy a new token upon the successful processing of a double-pegged burn prophecy claim if user is blocklisted, but should not mint", async function () {
      // Add recipient to the blocklist:
      await expect(state.blocklist.connect(operator).addToBlocklist(state.recipient.address)).to.be
        .fulfilled;

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: [userOne, userTwo, userFour],
      });

      // get the prophecyID:
      const prophecyID = await state.cosmosBridge.getProphecyID(
        state.sender, // cosmosSender
        state.senderSequence, // cosmosSenderSequence
        state.recipient.address, // ethereumReceiver
        state.token.address, // tokenAddress
        state.amount, // amount
        state.name, // tokenName
        state.symbol, // tokenSymbol
        state.decimals, // tokenDecimals
        state.networkDescriptor, // networkDescriptor
        true, // doublePeg
        state.nonce, // nonce
        state.constants.denom.one // cosmosDenom
      );

      const expectedAddress = ethers.utils.getContractAddress({
        from: state.bridgeBank.address,
        nonce: 1,
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      )
        .to.emit(state.cosmosBridge, "LogProphecyCompleted")
        .withArgs(prophecyID, false); // the second argument here is 'success'

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(
        state.token.address
      );
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

      // check if tokens were minted (should not have)
      const deployedToken = await state.factories.BridgeToken.attach(newlyCreatedTokenAddress);
      const mintedTokensToRecipient = await deployedToken.balanceOf(state.recipient.address);
      const totalSupply = await deployedToken.totalSupply();
      expect(totalSupply).to.be.equal(0);
      expect(mintedTokensToRecipient).to.be.equal(0);
    });

    it("should NOT deploy a new token upon the successful processing of a normal burn prophecy claim", async function () {
      state.nonce = 1;

      const beforeUserBalance = Number(await state.token.balanceOf(state.recipient.address));
      beforeUserBalance.should.be.bignumber.equal(Number(0));

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: [userOne, userTwo, userFour],
      });

      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(
        state.token.address
      );
      expect(newlyCreatedTokenAddress).to.be.equal(state.constants.zeroAddress);

      // assert that the recipient's balance of the token went up by the amount we specified in the claim
      balance = Number(await state.token.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(state.amount);
    });

    it("should NOT deploy a new token upon the successful processing of a double-pegged burn prophecy claim for an already managed token", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
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
          state.token.address,
          expectedAddress,
          state.constants.denom.one
        );

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(
        state.token.address
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
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce + 1,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
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
});
