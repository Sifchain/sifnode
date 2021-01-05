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
        operator
      ],
        {unsafeAllowCustomTypes: true}
      );

      console.log("Here: 1")
      // Deploy Oracle contract
      this.oracle = await deployProxy(Oracle,
        [
          operator,
          this.cosmosBridge.address,
          consensusThreshold,
          this.initialValidators,
          this.initialPowers
        ],
        {unsafeAllowCustomTypes: true}
        );
      console.log("Here: 2")

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank,[
        operator,
        this.oracle.address,
        this.cosmosBridge.address,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      // Operator sets Oracle
      await this.cosmosBridge.setOracle(this.oracle.address, {
        from: operator
      });

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
      


      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const prophecyClaimCount = event.args._prophecyID;
      // Get the ProphecyClaim's status
      let status = await this.cosmosBridge.isProphecyClaimActive(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(true);

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

      status = await this.cosmosBridge.isProphecyClaimActive(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(true);

      console.log("tx2: ", tx.receipt.gasUsed);
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

      console.log("tx3: ", tx.receipt.gasUsed);
      status = await this.cosmosBridge.isProphecyClaimActive(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(false);

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
      const status = await this.cosmosBridge.isProphecyClaimValidatorActive(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(true);
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

*/