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

  describe("Unlock Gas Cost With 4 Validators", function () {
    it("should allow us to check the cost of submitting a prophecy claim", async function () {
        
      // Lock tokens on contract
      await state.bridgeBank.connect(userOne).lock(
        state.sender,
        state.token1.address,
        state.amount
      ).should.be.fulfilled;
      
      state.cosmosSenderSequence = 10;
      state.nonce = 1;

      await state.cosmosBridge
        .connect(userOne)
        .newProphecyClaim(
          state.sender,
          state.senderSequence,
          state.recipient.address,
          state.token1.address,
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
          state.token1.address,
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
          state.token1.address,
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
          state.token1.address,
          state.amount,
          false,
          state.nonce
      );

      let prophecyID = await state.cosmosBridge.getProphecyID(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token1.address,
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


// Cost to unlock ethereum
/*

run: 1
tx:  399966
tx2:  151915
tx3:  217354
~~~~~~~~~~~~
Total: 769235

run: 2
tx:  368936
tx2:  103245
tx3:  151044
~~~~~~~~~~~~
Total: 623225

run: 2

tx:  355313
tx2:  89622
tx3:  137421
~~~~~~~~~~~~
Total: 582356

run: 3

tx:  355079
tx2:  89388
tx3:  137187
~~~~~~~~~~~~
Total: 581654

run: 4 (make newProphecyClaim external)

tx:  353990
tx2:  88705
tx3:  136503
~~~~~~~~~~~~
Total: 579198

run: 5 (combine oracle, valset and cosmosBridge together)
tx:  334064
tx2:  68773
tx3:  116571
~~~~~~~~~~~~
Total: 519408


run: 6 (cut down on storage used when creating prophecy claim)
tx:  230957
tx2:  68763
tx3:  112208
~~~~~~~~~~~~
Total: 411928

run: 7 (use 1 less storage slot when creating prophecy claim)
tx:  221869
tx2:  68763
tx3:  118444
~~~~~~~~~~~~
Total: 409076

run 8: (do not make call to BridgeBank to check if we have enough funds)
tx:  213875
tx2:  68763
tx3:  118444

~~~~Total Gas Used~~~~~
401082

run: 9 (use 2 less storage slots for the propheyClaim)
tx:  194043
tx2:  71652
tx3:  111847
~~~~~~~~~~~~
Total: 377542

run: 10 (remove prophecyClaim Count)
tx:  173135
tx2:  71652
tx3:  111847
~~~~~~~~~~~~
Total: 356634

run: 11 (remove usedNonce mapping)
tx:  152245
tx2:  71652
tx3:  111847
~~~~~~~~~~~~
Total: 335744

run: 12 (remove branching before calling newOracleClaim)
tx:  152241
tx2:  71638
tx3:  111833
~~~~~~~~~~~~
Total: 335712

run: 13 (add balance check back in)
tx:  160235
tx2:  71638
tx3:  111833
~~~~~~~~~~~~
Total: 343706

run: 14 (remove all use of ProphecyClaim stored in the struct inside of cosmos bridge and 100% leverage data in oracle contract)
tx:  97855
tx2:  71588
tx3:  108160
~~~~~~~~~~~~
Total: 277603

run: 15 (more EVM wizardry)
tx:  88797
tx2:  65469
tx3:  94453
~~~~~~~~~~~~
Total: 248719

*/


// Cost to mint erowan
/*
run: 1
tx:  89888
tx2:  65597
tx3:  290227
~~~~~~~~~~~~
Total: 445712

run: 2 (remove cosmos deposit stored in storage)
tx:  89888
tx2:  65597
tx3:  127339
~~~~~~~~~~~~
Total: 282824

run: 3 (remove function params)
tx:  89866
tx2:  65597
tx3:  126573
~~~~~~~~~~~~
Total: 282036

run: 4 (remove more function params)
tx:  89866
tx2:  65597
tx3:  126568
~~~~~~~~~~~~
Total: 282031
*/