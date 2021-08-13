const {
  setup,
  getValidClaim
} = require('./helpers/testFixture');

const web3 = require("web3");
const { expect } = require('chai');
const BigNumber = web3.BigNumber;

// Set `use` to `true` to compare a new implementation with the previous gas costs;
// Please set the previous gas costs accordingly in this object
const gasProfiling = {
  use: true,
  lock: 173978,
  mint: 179749,
  newBb: 1162769
}

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

describe.only("Gas Cost Tests", function () {
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
      console.log("~~~~~~~~~~~~\nTotal: ", sum);
      if(gasProfiling.use) {
        console.log("Improvement: ", gasProfiling.lock - sum);
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
        validators: accounts.slice(1, 5),
      });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);
      const receipt = await tx.wait();

      const sum = Number(receipt.gasUsed);
      console.log("~~~~~~~~~~~~\nTotal: ", sum);
      if(gasProfiling.use) {
        console.log("Improvement: ", gasProfiling.mint - sum);
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
      console.log("~~~~~~~~~~~~\nTotal: ", Number(receipt.gasUsed));
      if(gasProfiling.use) {
        console.log("Improvement: ", gasProfiling.newBb - Number(receipt.gasUsed));
      }

      const newlyCreatedTokenAddress = await state.cosmosBridge.sourceAddressToDestinationAddress(state.token1.address);
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);
    });
  });
});

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