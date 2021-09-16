const {
  setup,
  getValidClaim,
  batchAddTokensToEthWhitelist,
} = require('./helpers/testFixture');

const { colorLog } = require('./helpers/helpers');

const web3 = require("web3");
const { expect } = require('chai');
const BigNumber = web3.BigNumber;

// Set `use` to `true` to compare a new implementation with the previous gas costs;
// Please set the previous gas costs accordingly in this object
const gasProfiling = {
  use: false,
  lock: 175867,
  mint: 181881,
  newBt: 1665627,
  multiLock: 346709,
  current: {}
}

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
    pauser = accounts[6];

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
    state = await setup({
      initialValidators,
      initialPowers,
      operator,
      consensusThreshold,
      owner,
      user: userOne,
      recipient: userThree,
      pauser,
      networkDescriptor
    });

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
  });

  describe("Gas Cost With 4 Validators", function () {
    it("should allow us to check the cost of submitting a prophecy claim lock", async function () {
      let balance = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(0);

      // Last nonce should now be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token1.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: accounts.slice(1, 5),
      });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        );
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);
      
      if(gasProfiling.use) {
        gasProfiling.current.lock = sum;
        logGasDiff('LOCK:', gasProfiling.lock, sum);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", sum);
      }

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

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.rowan.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.rowan,
        validators: accounts.slice(1, 5),
      });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);

      if(gasProfiling.use) {
        gasProfiling.current.mint = sum;
        logGasDiff('MINT:', gasProfiling.mint, sum);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", sum);
      }

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      // balance should have increased
      balance = Number(await state.rowan.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(state.amount);
    });

    it("should allow us to check the cost of creating a new BridgeToken", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token1.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: accounts.slice(1, 5),
      });

      const expectedAddress = ethers.utils.getContractAddress({ from: state.bridgeBank.address, nonce: 1 });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(
            digest,
            claimData,
            signatures
        );
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);
      
      if(gasProfiling.use) {
        gasProfiling.current.newBt = sum;
        logGasDiff('DoublePeg :: New BridgeToken:', gasProfiling.newBt, sum);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", Number(receipt.gasUsed));
      }

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token1.address);
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

      // expect the token to have a denom
      const registeredDenom = await state.bridgeBank.contractDenom(newlyCreatedTokenAddress);
      expect(registeredDenom).to.be.equal(state.constants.denom.one);
    });

    it("should allow us to check the cost of submitting a batch prophecy claim lock", async function () {
      // Add tokens 2 and 3 into white list
      await batchAddTokensToEthWhitelist(state, [state.token2.address, state.token3.address]);

      // Make sure the tokens were added to the whitelist
      const isToken2InWhitelist = await state.bridgeBank.getTokenInEthWhiteList(state.token2.address);
      expect(isToken2InWhitelist).to.be.equal(true);

      const isToken3InWhitelist = await state.bridgeBank.getTokenInEthWhiteList(state.token3.address);
      expect(isToken3InWhitelist).to.be.equal(true);
      
      // Lock token2 on contract
      await state.bridgeBank.connect(userOne).lock(
        state.sender,
        state.token2.address,
        state.amount
      ).should.be.fulfilled;

      // Lock token3 on contract
      await state.bridgeBank.connect(userOne).lock(
        state.sender,
        state.token3.address,
        state.amount
      ).should.be.fulfilled;

      let balanceToken1 = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balanceToken1).to.be.equal(0);

      let balanceToken2 = Number(await state.token2.balanceOf(state.recipient.address));
      expect(balanceToken2).to.be.equal(0);

      let balanceToken3 = Number(await state.token3.balanceOf(state.recipient.address));
      expect(balanceToken3).to.be.equal(0);

      // Last nonce should now be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token1.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: accounts.slice(1, 5),
      });

      const { digest: digest2, claimData: claimData2, signatures: signatures2 } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token2.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce + 1,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.two,
        validators: accounts.slice(1, 5),
      });

      const { digest: digest3, claimData: claimData3, signatures: signatures3 } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token3.address,
        amount: state.amount,
        doublePeg: false,
        nonce: state.nonce + 2,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.three,
        validators: accounts.slice(1, 5),
      });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .batchSubmitProphecyClaimAggregatedSigs(
            [digest, digest2, digest3],
            [claimData, claimData2, claimData3],
            [signatures, signatures2, signatures3]
        );
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);
      
      if(gasProfiling.use) {
        const numberOfClaimsInBatch = 3;
        logGasDiff('BATCH, regarding previous implementation:', gasProfiling.multiLock, sum);
        logGasDiff(`${numberOfClaimsInBatch} Batched claims VS single claim:`, gasProfiling.current.lock * numberOfClaimsInBatch, sum, true);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", sum);
      }

      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(3);

      balanceToken1 = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balanceToken1).to.be.equal(state.amount);

      balanceToken2 = Number(await state.token2.balanceOf(state.recipient.address));
      expect(balanceToken2).to.be.equal(state.amount);

      balanceToken3 = Number(await state.token3.balanceOf(state.recipient.address));
      expect(balanceToken3).to.be.equal(state.amount);
    });

    it("should allow us to check the cost of batch creating new BridgeTokens", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token1.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.one,
        validators: accounts.slice(1, 5),
      });

      const { digest: digest2, claimData: claimData2, signatures: signatures2 } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token2.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce + 1,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.two,
        validators: accounts.slice(1, 5),
      });

      const { digest: digest3, claimData: claimData3, signatures: signatures3 } = await getValidClaim({
        sender: state.sender,
        senderSequence: state.senderSequence,
        recipientAddress: state.recipient.address,
        tokenAddress: state.token3.address,
        amount: state.amount,
        doublePeg: true,
        nonce: state.nonce + 2,
        networkDescriptor: state.networkDescriptor,
        tokenName: state.name,
        tokenSymbol: state.symbol,
        tokenDecimals: state.decimals,
        cosmosDenom: state.constants.denom.three,
        validators: accounts.slice(1, 5),
      });

      const expectedAddress1 = ethers.utils.getContractAddress({ from: state.bridgeBank.address, nonce: 1 });
      const expectedAddress2 = ethers.utils.getContractAddress({ from: state.bridgeBank.address, nonce: 2 });
      const expectedAddress3 = ethers.utils.getContractAddress({ from: state.bridgeBank.address, nonce: 3 });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .batchSubmitProphecyClaimAggregatedSigs(
          [digest, digest2, digest3],
          [claimData, claimData2, claimData3],
          [signatures, signatures2, signatures3]
        );
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);
      
      if(gasProfiling.use) {
        const numberOfClaimsInBatch = 3;
        logGasDiff(`DoublePeg :: ${numberOfClaimsInBatch} Batched claims VS single claim:`, gasProfiling.current.newBt * numberOfClaimsInBatch, sum);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", Number(receipt.gasUsed));
      }

      const newlyCreatedTokenAddress1 = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token1.address);
      expect(newlyCreatedTokenAddress1).to.be.equal(expectedAddress1);

      const newlyCreatedTokenAddress2 = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token2.address);
      expect(newlyCreatedTokenAddress2).to.be.equal(expectedAddress2);

      const newlyCreatedTokenAddress3 = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token3.address);
      expect(newlyCreatedTokenAddress3).to.be.equal(expectedAddress3);
    });

    it("should allow us to check the cost of adding a token to Eth Whitelist", async function () {
      // First, remove token1 from the whitelist so that we can add it again
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, false)
        .should.be.fulfilled;

      // Measure gas costs of adding a single token into Eth WhiteList
      const tx = await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true);
      
      const receipt = await tx.wait();
      gasProfiling.current.addToEthWhiteList = Number(receipt.gasUsed);

      // Now, we remove token1 from the whitelist to be able to add it again in a batch
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, false)
        .should.be.fulfilled;
      
      const batchTx = await state.bridgeBank.connect(operator)
        .batchUpdateEthWhiteList(
          [state.token1.address, state.token2.address, state.token3.address],
          [true, true, true]
        );

      const batchReceipt = await batchTx.wait();
      const sum = Number(batchReceipt.gasUsed);

      if(gasProfiling.use) {
        const numberOfTokensInBatch = 3;
        logGasDiff(
          `Batch Whitelist :: ${numberOfTokensInBatch} Batched VS single:`,
          gasProfiling.current.addToEthWhiteList * numberOfTokensInBatch, sum
        );
      }
    });
  });
});

// Helper function to aid comparing implementations wrt gas costs
function logGasDiff(title, original, current, useShortTitle) {
  const separator = useShortTitle ? '---' : '~~~~~~~~~~~~';
  colorLog('cyan', `${separator}\n${title}`);
  console.log('Original:', original);
  console.log('Current :', current);
  const pct = Math.abs(((1 - current / original) * 100)).toFixed(2);
  const diff = current - original;
  colorLog(getColorName(diff), `Diff    : ${diff} (${pct}%)`);
}

function getColorName(value) {
  if(value > 0) {
    return 'red';
  } else if(value < 0) {
    return 'green';
  } else {
    return 'white';
  }
}

/**
 * 
 * 
Unlock Gas Cost With 4 Validators
tx0  173978
~~~~~~~~~~~~
Total:  173978

Mint Gas Cost With 4 Validators
tx0  179749
~~~~~~~~~~~~
Total:  179749

Create new BridgeToken Gas Cost With 4 Validators
tx0  1162769
~~~~~~~~~~~~
Total:  1162769
 * 
 */