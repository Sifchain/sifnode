const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

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
  let state;

  before(async function() {
    accounts = await ethers.getSigners();
    
    signerAccounts = accounts.map((e) => { return e.address });

    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[7].address;

    owner = accounts[5];
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);
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
      pauser
    );
  });

  describe("CosmosBridge", function () {
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

    it("should unlock tokens upon the successful processing of a burn prophecy claim", async function () {
      const beforeUserBalance = Number(
        await state.token.balanceOf(state.recipient)
      );
      beforeUserBalance.should.be.bignumber.equal(Number(0));

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient,
        state.token.address,
        state.amount,
        false,
        state.nonce
      ]);

      const signatures = await signHash([userOne, userTwo, userFour], digest);
      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient,
        tokenAddress: state.token.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce
      };

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

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient,
        state.ethereumToken,
        state.amount,
        false,
        state.nonce
      ]);
      const signatures = await signHash([userOne, userTwo, userFour], digest);

      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient,
        tokenAddress: state.ethereumToken,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce
      };

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
  });
});
