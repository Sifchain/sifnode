const {
  signHash,
  setup,
  deployTrollToken,
  getDigestNewProphecyClaim,
  getValidClaim,
} = require("./helpers/testFixture");

const web3 = require("web3");
const { expect } = require("chai");
const BigNumber = web3.BigNumber;

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

describe("submitProphecyClaimAggregatedSigs Security", function () {
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

  before(async function () {
    accounts = await ethers.getSigners();

    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userThree = accounts[3];
    userFour = accounts[4];

    owner = accounts[5];
    pauser = accounts[6];

    initialPowers = [25, 25, 25, 25];
    initialValidators = [userOne.address, userTwo.address, userThree.address, userFour.address];

    networkDescriptor = 1;
  });

  beforeEach(async function () {
    // Deploy Valset contract
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
    });

    // Lock tokens on contract
    await state.bridgeBank
      .connect(userOne)
      .lock(state.sender, state.token1.address, state.amount).should.be.fulfilled;

    let TrollToken = await deployTrollToken();
    state.troll = TrollToken;
    await state.troll.mint(userOne.address, 100);
  });

  describe("should revert when", function () {
    it("no signatures are provided", async function () {
      state.recipient = userOne;
      state.nonce = 1;

      const { digest, claimData } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
        validators: [userOne, userTwo, userFour],
      });

      await expect(
        state.cosmosBridge.connect(userOne).submitProphecyClaimAggregatedSigs(digest, claimData, [])
      ).to.be.revertedWith("INV_SIG_LEN");
    });

    it("hash digest doesn't match provided data", async function () {
      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence + 1,
        state.recipient.address,
        state.troll.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce + 1,
        state.constants.denom.none,
      ]);

      const signatures = await signHash([userOne, userTwo, userFour], digest);
      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
      };

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_DATA");
    });

    it("there are duplicate signers", async function () {
      state.recipient = userOne;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
        validators: [userOne, userOne, userTwo, userFour],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("custom error 'DuplicateSigner(3, \"" + userOne.address + "\")'");
    });

    it("there is an invalid signer", async function () {
      state.recipient = userOne;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
        validators: [userOne, userTwo, operator],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_SIGNER");
    });

    it("there is a signature that signs invalid data", async function () {
      state.recipient = userOne;
      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.troll.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.none,
      ]);

      const invalidDigest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence + 1,
        state.recipient.address,
        state.troll.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce + 1,
        state.constants.denom.none,
      ]);

      const signatures = await signHash([userOne], digest);
      const invalidSig = await signHash([userFour], invalidDigest);

      // push this signature onto the valid signature array
      signatures.push(invalidSig[0]);

      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
      };

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("custom error 'OutOfOrderSigner(0)'");
    });

    it("there is not enough power to complete prophecy", async function () {
      state.recipient = userOne;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
        validators: [userOne, userTwo],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_POW");
    });

    it("prophecy is in an invalid order", async function () {
      state.recipient = userOne;
      state.nonce = 2;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
        validators: [userOne, userTwo, userFour],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_ORD");
    });

    it("prophecy is already redeemed", async function () {
      state.recipient = userOne;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.troll.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.none,
        validators: [userOne, userTwo, userFour],
      });

      state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_ORD");
    });

    it("one of the claims in a batch prophecy claim has the wrong nonce", async function () {
      // Lock token2 on contract
      await state.bridgeBank.connect(userOne).lock(state.sender, state.token2.address, state.amount)
        .should.be.fulfilled;

      // Lock token3 on contract
      await state.bridgeBank.connect(userOne).lock(state.sender, state.token3.address, state.amount)
        .should.be.fulfilled;

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token1.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: accounts.slice(1, 5),
      });

      const {
        digest: digest2,
        claimData: claimData2,
        signatures: signatures2,
      } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token2.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce + 2, // this should be rejected because the expected value is state.nonce + 1
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.two,
        validators: accounts.slice(1, 5),
      });

      const {
        digest: digest3,
        claimData: claimData3,
        signatures: signatures3,
      } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token3.address,
        amount: state.amount,
        bridgeToken: false,
        nonce: state.nonce + 2,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.three,
        validators: accounts.slice(1, 5),
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .batchSubmitProphecyClaimAggregatedSigs(
            [digest, digest2, digest3],
            [claimData, claimData2, claimData3],
            [signatures, signatures2, signatures3]
          )
      ).to.be.rejectedWith("INV_ORD");

      // global nonce should not have changed:
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);
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
