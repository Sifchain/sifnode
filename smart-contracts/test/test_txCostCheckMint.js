const { multiTokenSetup } = require('./helpers/testFixture');

const web3 = require("web3");
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
  let state;

  before(async function() {
    accounts = await ethers.getSigners();
    
    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[7];

    owner = accounts[5];
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = [
      userOne.address,
      userTwo.address,
      userThree.address,
      userFour.address
    ];
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
      pauser
    );
  });

  describe("Mint Bridge Token Gas Cost With 4 Validators", function () {
    it("should allow us to check the cost of submitting a prophecy claim", async function () {
      state.cosmosSenderSequence = 10;
      state.nonce = 1;

      await state.cosmosBridge
        .connect(userOne)
        .newProphecyClaim(
          state.sender,
          state.senderSequence,
          state.recipient.address,
          state.rowan.address,
          state.amount,
          false,
          state.nonce
      );

      // Create the prophecy claim
      await state.cosmosBridge
        .connect(userTwo)
        .newProphecyClaim(
          state.sender,
          state.senderSequence,
          state.recipient.address,
          state.rowan.address,
          state.amount,
          false,
          state.nonce
      );

      await state.cosmosBridge
        .connect(userThree)
        .newProphecyClaim(
          state.sender,
          state.senderSequence,
          state.recipient.address,
          state.rowan.address,
          state.amount,
          false,
          state.nonce
      );

      await state.cosmosBridge
        .connect(userFour)
        .newProphecyClaim(
          state.sender,
          state.senderSequence,
          state.recipient.address,
          state.rowan.address,
          state.amount,
          false,
          state.nonce
      );

      let prophecyID = await state.cosmosBridge.getProphecyID(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.rowan.address,
        state.amount
      );

      status = await state.cosmosBridge.prophecyRedeemed(
        prophecyID.toString(),
      );

      // Bridge claim should be completed
      status.should.be.equal(true);
    });
  });
});
