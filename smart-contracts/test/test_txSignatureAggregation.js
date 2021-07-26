const {
  multiTokenSetup,
  signHash,
  getDigestNewProphecyClaim
} = require('./helpers/testFixture');

const web3 = require("web3");
const { expect } = require('chai');
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

describe("Gas Cost Tests", function () {
  let userOne;
  let userTwo;
  let userThree;
  let userFour;
  let accounts;
  let operator;
  let owner;
  let pauser;

  // Consensus threshold of 70%
  const consensusThreshold = 70;
  let initialPowers;
  let initialValidators;
  let networkDescriptor;
  let state;

  before(async function() {
    accounts = await ethers.getSigners();
    
    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userThree = accounts[3];
    userFour = accounts[4];

    owner = accounts[5];
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = [
      userOne.address,
      userTwo.address,
      userThree.address,
      userFour.address
    ];

    networkDescriptor = 1;
  });

  beforeEach(async function () {
    // Deploy Valset contract
    state = await multiTokenSetup(
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

    // Lock tokens on contract
    await state.bridgeBank.connect(userOne).lock(
      state.sender,
      state.token1.address,
      state.amount
    ).should.be.fulfilled;
  });

  describe("Unlock Gas Cost With 4 Validators", function () {
    it("should allow us to check the cost of submitting a prophecy claim lock", async function () {
      let balance = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(0);

      // Last nonce should now be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token1.address,
        state.amount,
        false,
        state.nonce,
        state.networkDescriptor
      ]);

      let validators = accounts.slice(1, 5);
      const signatures = await signHash(validators, digest);
      let sum = 0;

      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient.address,
        tokenAddress: state.token1.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor
      };

      let tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        );
      let receipt = await tx.wait();
      sum += Number(receipt.gasUsed);

      console.log("~~~~~~~~~~~~\nTotal: ", sum);

      // Bridge claim should be completed
      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      balance = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(state.amount);
    });

    it("should allow us to check the cost of submitting a prophecy claim mint", async function () {
      let balance = Number(await state.rowan.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(0);

      // Last nonce should now be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.rowan.address,
        state.amount,
        false,
        state.nonce,
        state.networkDescriptor
      ]);

      let validators = accounts.slice(1, 5);
      const signatures = await signHash(validators, digest);
      let sum = 0;

      let claimData = {
          cosmosSender: state.sender,
          cosmosSenderSequence: state.senderSequence,
          ethereumReceiver: state.recipient.address,
          tokenAddress: state.rowan.address,
          amount: state.amount,
          doublePeg: false,
          nonce: state.nonce,
          networkDescriptor: state.networkDescriptor
      };

      let tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        );
      let receipt = await tx.wait();
      sum += Number(receipt.gasUsed);

      console.log("~~~~~~~~~~~~\nTotal: ", sum);

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      // balance should have increased
      balance = Number(await state.rowan.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(state.amount);
    });
  });
});

/**
 * 
 * 
Unlock Gas Cost With 4 Validators
tx0  182434
~~~~~~~~~~~~
Total:  182434

Mint Gas Cost With 4 Validators
tx0  198100
~~~~~~~~~~~~
Total:  198100
 * 
 */