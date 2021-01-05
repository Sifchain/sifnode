const { deployProxy, silenceWarnings } = require('@openzeppelin/truffle-upgrades');
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");

const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("CosmosBridge", function (accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];
  const userFour = accounts[4];

  // Contract's enum ClaimType can be represented a sequence of integers
  const CLAIM_TYPE_BURN = 1;
  const CLAIM_TYPE_LOCK = 2;

  // Consensus threshold of 70%
  const consensusThreshold = 70;

  describe("Bridge claim status", function () {
    beforeEach(async function () {
      // Set up ProphecyClaim values
      this.cosmosSender = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      this.cosmosSenderSequence = 1;
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "ETH";
      this.amount = 100;

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      // this.valset = await deployProxy(Valset,
      //   [
      //     operator,
      //     this.initialValidators,
      //     this.initialPowers
      //   ],
      //   {unsafeAllowCustomTypes: true}
      // );

      // Deploy CosmosBridge contract
      console.log("Here: 0")

      this.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        this.initialValidators,
        this.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );

      console.log("Here: 1")
      // Deploy Oracle contract
      // this.oracle = await deployProxy(Oracle,
      //   [
      //     operator,
      //     this.cosmosBridge.address,
      //     consensusThreshold,
      //     this.initialValidators,
      //     this.initialPowers
      //   ],
      //   {unsafeAllowCustomTypes: true}
      //   );
      // console.log("Here: 2")

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank,[
        operator,
        this.cosmosBridge.address,
        this.cosmosBridge.address,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      // Operator sets Bridge Bank
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      });

      this.recipient = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );

      this.weiAmount = web3.utils.toWei("0.25", "ether");

      // Update the lock/burn limit for this token
      await this.bridgeBank.updateTokenLockBurnLimit(this.tokenAddress, this.weiAmount, {
        from: operator
      }).should.be.fulfilled;

      await this.bridgeBank.lock(
        this.recipient,
        this.tokenAddress,
        this.weiAmount, {
          from: userOne,
          value: this.weiAmount
        }
      ).should.be.fulfilled;
    });

    it("should allow us to check the cost of submitting a prophecy claim", async function () {
        let sum = 0;
        this.cosmosSenderSequence = 10;
        const estimatedGas = await this.cosmosBridge.newProphecyClaim.estimateGas(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amount,
          {
              from: userOne
          }
        );
      // console.log("Params: ", CLAIM_TYPE_LOCK, this.cosmosSender, this.cosmosSenderSequence, this.ethereumReceiver, this.symbol, this.amount)
        // Create the prophecy claim
      let {receipt, logs} = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userOne,
          gasPrice: "1"
        }
      );
      console.log("Estimated Gas: ", estimatedGas)
      console.log("Gas price: ", await web3.eth.getGasPrice())
      sum += receipt.gasUsed


      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const prophecyClaimCount = event.args._prophecyID;
      // Get the ProphecyClaim's status
      // let status = await this.cosmosBridge.isProphecyClaimActive(
      //   prophecyClaimCount,
      //   {
      //     from: accounts[7]
      //   }
      // );

      // // Bridge claim should be active
      // status.should.be.equal(true);

      console.log("tx: ", receipt.gasUsed)

      let tx = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userTwo,
          gasPrice: "1"
        }
      );

      // status = await this.cosmosBridge.isProphecyClaimActive(
      //   prophecyClaimCount,
      //   {
      //     from: accounts[7]
      //   }
      // );

      // // Bridge claim should be active
      // status.should.be.equal(true);

      console.log("tx2: ", tx.receipt.gasUsed);
      sum += tx.receipt.gasUsed
      tx = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userThree,
          gasPrice: "1"
        }
      );
      sum += tx.receipt.gasUsed

      console.log("tx3: ", tx.receipt.gasUsed);
      // status = await this.cosmosBridge.isProphecyClaimActive(
      //   prophecyClaimCount,
      //   {
      //     from: accounts[7]
      //   }
      // );

      // Bridge claim should be active
      // status.should.be.equal(false);
      console.log(`~~~~~~~~~~~~\nTotal: ${sum}`);

    });

    it("should allow users to check if a prophecy claim's original validator is currently an active validator", async function () {
      // Create the ProphecyClaim
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const prophecyClaimCount = event.args._prophecyID;

      // Get the ProphecyClaim's status
      // const status = await this.cosmosBridge.isProphecyClaimActive(
      //   prophecyClaimCount,
      //   {
      //     from: accounts[7]
      //   }
      // );

      // // Bridge claim should be active
      // status.should.be.equal(true);
    });
  });
});

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

*/