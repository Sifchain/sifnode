const {
  signHash,
  multiTokenSetup,
  deployTrollToken,
  getDigestNewProphecyClaim,
  getValidClaim
} = require('./helpers/testFixture');

const web3 = require("web3");
const { expect } = require('chai');
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

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

    // Add the token into white list
    await state.bridgeBank.connect(operator)
      .updateEthWhiteList(state.token1.address, true)
      .should.be.fulfilled;

    // Lock tokens on contract
    await state.bridgeBank.connect(userOne).lock(
      state.sender,
      state.token1.address,
      state.amount
    ).should.be.fulfilled;

    let TrollToken = await deployTrollToken();
    state.troll = TrollToken;
    await state.troll.mint(userOne.address, 100);
  });

  describe("should revert when", function () {
    it("no signatures are provided", async function () {
      state.recipient = userOne.address;
      state.nonce = 10;

      const { digest, claimData } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            []
          )
      ).to.be.revertedWith("INV_SIG_LEN");
    });
    
    it("hash digest doesn't match provided data", async function () {
      state.recipient = userOne.address;
      state.nonce = 10;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient,
        state.troll.address,
        state.amount,
        false,
        state.nonce + 1,
        state.networkDescriptor
      ]);

      const signatures = await signHash([userOne, userTwo, userFour], digest);
      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals
      };

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
          )
      ).to.be.revertedWith("INV_DATA");
    });

    it("there are duplicate signers", async function () {
      state.recipient = userOne.address;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour, userFour],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
          )
      ).to.be.revertedWith("DUP_SIGNER");
    });

    it("there is an invalid signer", async function () {
      state.recipient = userOne.address;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, operator],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
          )
      ).to.be.revertedWith("INV_SIGNER");
    });

    it("there is a signature that signs invalid data", async function () {
      state.recipient = userOne.address;
      state.nonce = 1;
      const digest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient,
        state.troll.address,
        state.amount,
        false,
        state.nonce,
        state.networkDescriptor
      ]);

      const invalidDigest = getDigestNewProphecyClaim([
        state.sender,
        state.senderSequence,
        state.recipient,
        state.troll.address,
        state.amount,
        false,
        state.nonce + 1,
        state.networkDescriptor
      ]);

      const signatures = await signHash([userOne, userTwo], digest);
      const invalidSig = await signHash([userFour], invalidDigest);

      // push this signature onto the valid signature array
      signatures.push(invalidSig[0]);

      let claimData = {
        cosmosSender: state.sender,
        cosmosSenderSequence: state.senderSequence,
        ethereumReceiver: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
      };

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
          )
      ).to.be.revertedWith("INV_SIG");
    });

    it("there is not enough power to complete prophecy", async function () {
      state.recipient = userOne.address;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
          )
      ).to.be.revertedWith("INV_POW");
    });

    it("prophecy is in an invalid order", async function () {
      state.recipient = userOne.address;
      state.nonce = 2;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
          )
      ).to.be.revertedWith("INV_ORD");
    });

    it("prophecy is already redeemed", async function () {
      state.recipient = userOne.address;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        state,
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient,
        tokenAddress: state.troll.address,
        amount: state.amount,
        isDoublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        validators: [userOne, userTwo, userFour],
      });

      state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
      );

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
          )
      ).to.be.revertedWith("INV_ORD");
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