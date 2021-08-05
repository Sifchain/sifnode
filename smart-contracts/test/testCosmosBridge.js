const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const {
  singleSetup,
  getValidClaim
} = require("./helpers/testFixture");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

use(solidity);

const getBalance = async function(address) {
  return await network.provider.send("eth_getBalance", [address]);
}

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

    networkDescriptor = 1;
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
      networkDescriptor
    );
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
      await state.cosmosBridge.connect(operator).updateValset(
        [userOne.address, userTwo.address],
        [50, 50],
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
        userThree
      );
      isUserThreeValidator.should.be.equal(false);
      const isUserFourValidator = await state.cosmosBridge.isActiveValidator(
        userFour.address
      );
      isUserFourValidator.should.be.equal(false);
    });

    it("Can update the validator set", async function () {
      // Also make sure everything runs fourth time after switching validators a second time.
      // Operator resets the valset
      await state.cosmosBridge.connect(operator).updateValset(
        [userThree, userFour.address],
        [50, 50],
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
          userThree
      );
      isUserThreeValidator2.should.be.equal(true);
      const isUserFourValidator2 = await state.cosmosBridge.isActiveValidator(
          userFour.address
      );
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
      expect(
        (await state.cosmosBridge.consensusThreshold()).toString()
      ).to.equal(consensusThreshold.toString());

      // iterate over all validators and ensure they have the proper
      // powers and that they have been succcessfully whitelisted
      for (let i = 0; i < initialValidators.length; i++) {
        const address = initialValidators[i];

        expect(
          await state.cosmosBridge.isActiveValidator(address)
        ).to.be.true;

        expect(
          (await state.cosmosBridge.getValidatorPower(address)).toString()
        ).to.equal("25");
      }

      expect(await state.bridgeBank.cosmosBridge()).to.be.equal(state.cosmosBridge.address);
      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);
      expect(await state.bridgeBank.pausers(pauser)).to.be.true;
    });

    it("Should deploy cosmos bridge and bridge bank, correctly setting the networkDescriptor", async function () {
        expect(await state.cosmosBridge.networkDescriptor()).to.equal(state.networkDescriptor);
        expect(await state.bridgeBank.networkDescriptor()).to.equal(state.networkDescriptor);
    });

    it("should unlock tokens upon the successful processing of a burn prophecy claim", async function () {
      const beforeUserBalance = Number(
        await state.token.balanceOf(state.recipient)
      );
      beforeUserBalance.should.be.bignumber.equal(Number(0));

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.token.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });
     
      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        );

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      balance = Number(await state.token.balanceOf(state.recipient));
      expect(balance).to.be.equal(state.amount);
    });

    it("should unlock eth upon the successful processing of a burn prophecy claim", async function () {
      // assert recipient balance before receiving proceeds from newProphecyClaim is correct
      const recipientStartingBalance = await getBalance(state.recipient);
      const recipientCurrentBalance = Web3Utils.fromWei(recipientStartingBalance);

      expect(recipientCurrentBalance).to.be.equal(
        "10000"
      );
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.ethereumToken,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: "Ether",
        tokenSymbol: "ETH",
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      
      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        );

      const recipientEndingBalance = await getBalance(state.recipient);
      const recipientBalance = Web3Utils.fromWei(recipientEndingBalance);

      expect(recipientBalance).to.be.equal(
        "10000.0000000000000001"
      );
    });

    it("should deploy a new token upon the successful processing of a double-pegged burn prophecy claim", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.token.address,
        amount: state.amount,
        isDoublePeg: true,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      const expectedAddress = ethers.utils.getContractAddress({ from: state.bridgeBank.address, nonce: 1 });

      await expect(state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        )).to.emit(state.cosmosBridge, 'LogNewBridgeTokenCreated')
        .withArgs(18, state.networkDescriptor, state.name, state.symbol, state.token.address, expectedAddress);

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token.address);
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);
    });

    it("should NOT deploy a new token upon the successful processing of a normal burn prophecy claim", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.token.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token.address);
      expect(newlyCreatedTokenAddress).to.be.equal(state.constants.zeroAddress);
    });

    it("should NOT deploy a new token upon the successful processing of a double-pegged burn prophecy claim for an already managed token", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.token.address,
        amount: state.amount,
        isDoublePeg: true,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      const expectedAddress = ethers.utils.getContractAddress({ from: state.bridgeBank.address, nonce: 1 });
      await expect(state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        )).to.emit(state.cosmosBridge, 'LogNewBridgeTokenCreated')
        .withArgs(18, state.networkDescriptor, state.name, state.symbol, state.token.address, expectedAddress);

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token.address);
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

      // Everything again, but this time submitProphecyClaimAggregatedSigs should NOT emit the event
      const { digest: digest2, claimData: claimData2, signatures: signatures2 } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence + 1,
        recipientAddress: state.recipient,
        tokenAddress: state.token.address,
        amount: state.amount,
        isDoublePeg: true,
        nonce: state.nonce + 1,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      await expect(state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest2,
            claimData2,
            signatures2
        )).to.not.emit(state.cosmosBridge, 'LogNewBridgeTokenCreated');
    });
  });
});
