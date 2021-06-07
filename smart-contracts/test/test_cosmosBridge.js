const { deployProxy, silenceWarnings } = require('@openzeppelin/truffle-upgrades');
const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");
const BridgeToken = artifacts.require("BridgeToken");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

const {
  expectRevert, // Assertions for transactions that should fail
} = require('@openzeppelin/test-helpers');

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

  // Default Peggy token prefix
  const defaultTokenPrefix = "e"
  describe("CosmosBridge smart contract deployment", function () {
    beforeEach(async function () {
      // Deploy Valset contract
      this.initialValidators = [userOne.address, userTwo.address, userThree.address, userFour.address];
      this.initialPowers = [30, 20, 21, 29];
      
      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        this.initialValidators,
        this.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );
      
      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        this.cosmosBridge.address,
        operator,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );
    });

    it("should not allow the operator to update the Bridge Bank once it has been set", async function () {
      await this.cosmosBridge
        .setBridgeBank(operator, {
          from: operator
        })
        .should.be.rejectedWith(EVMRevert);
    });
  });

  describe("Creation of prophecy claims", function () {
    beforeEach(async function () {
      // Set up ProphecyClaim values
      this.cosmosSender = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      this.cosmosSenderSequence = 1;
      this.ethereumReceiver = userThree;

      // Deploy Valset contract
      this.initialValidators = [
        userOne.address,
        userTwo.address,
        userThree.address,
        userFour.address
      ];
      this.initialPowers = [30, 20, 21, 29];
      
      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        this.initialValidators,
        this.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        this.cosmosBridge.address,
        operator,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      // Fail to set BridgeBank if not the operator.
      await expectRevert(
          this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
            from: userOne
          }),
          "Must be the operator."
      );

      // Operator sets Bridge Bank
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      });

      // Fail to set BridgeBank a second time.
      await expectRevert(
          this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
            from: operator
          }),
          "The Bridge Bank cannot be updated once it has been set"
      );
    });

    it("should allow for the creation of new burn prophecy claims", async function () {
      // Load user account with ERC20 tokens
      await this.token.mint(userOne, 2000, {
        from: operator
      }).should.be.fulfilled;

      // Approve tokens to contract
      await this.token.approve(this.bridgeBank.address, this.amount, {
        from: userOne
      }).should.be.fulfilled;

      const tx = await this.bridgeBank.lock(
        this.cosmosRecipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      ).should.be.fulfilled;

      const logs = await tx.wait();

      const event = logs.find(e => e.event === "LogLock");
      event.args._token.should.be.equal(this.token.address);
      event.args._symbol.should.be.equal(this.actualSymbol);
      Number(event.args._value).should.be.equal(Number(this.amount));

      const nonce = 0;

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.cosmosSenderSequence,
        userFour,
        this.actualSymbol.toLowerCase(),
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    it("should not allow for the creation of a new burn prophecy claim over current amount locked", async function () {
      await expectRevert(
          this.cosmosBridge.newProphecyClaim(
              CLAIM_TYPE_BURN,
              this.cosmosSender,
              ++this.cosmosSenderSequence,
              this.ethereumReceiver,
              this.symbol,
              1,
              {
                from: userOne
              }
          ),
          "Not enough locked assets to complete the proposed prophecy"
      );
    });

    it("should allow correct operator to change the operator", async function () {
      await this.cosmosBridge.changeOperator(userTwo, { from: operator })
        .should.be.fulfilled;
      (await this.cosmosBridge.operator()).should.be.equal(userTwo);
    });

    it("should not allow incorrect operator to change the operator", async function () {
      await expectRevert(
        this.cosmosBridge.changeOperator(
            userTwo,
            {
              from: userOne
            }
        ),
        "Must be the operator."
      );
      (await this.cosmosBridge.operator()).should.be.equal(operator);
    });

    it("should not allow for anything other than BURN/LOCK (1 or 2)", async function () {
      await this.cosmosBridge.newProphecyClaim(
          3,
          this.cosmosSender,
          ++this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amount,
          {
            from: userOne
          }
      ).should.be.rejectedWith(EVMRevert);
    });

    it("should allow for the creation of new lock prophecy claims", async function () {
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        ++this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    it("should log an event containing the new prophecy claim's information", async function () {
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        ++this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol.toLowerCase(),
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      const event = logs.find(e => e.event === "LogNewProphecyClaim");

      Number(event.args._claimType).should.be.equal(CLAIM_TYPE_LOCK);

      event.args._ethereumReceiver.should.be.equal(this.ethereumReceiver);
      event.args._symbol.should.be.equal(defaultTokenPrefix + this.symbol);
      Number(event.args._amount).should.be.equal(this.amount);
    });

    it("should be able to create a new prophecy claim", async function () {
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        ++this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

    });

    it("should not allow a eth to be locked if the amount is over the limit", async function () {
      const maxLockAmount = Number(await this.bridgeBank.maxTokenAmount("ETH"));
      // Calculate and check expected max lock amount
      maxLockAmount.should.be.equal(Number(0));
      
      await expectRevert(
        this.bridgeBank.lock(
          this.cosmosRecipient,
          this.ethereumToken,
          this.amount, {
            from: userOne,
            value: this.amount
          }
        ),
        "Amount being transferred is over the limit"
      );
    });

    it("should not allow a token to be locked if the amount is over the limit", async function () {
      // Update the lock/burn limit for this token
      await this.bridgeBank.updateTokenLockBurnLimit(this.token.address, 0, {
        from: operator
      }).should.be.fulfilled;
      
      const maxLockAmount = Number(await this.bridgeBank.maxTokenAmount(await this.token.symbol()));
      // Calculate and check expected balances
      maxLockAmount.should.be.equal(Number(0));
      
      // Approve tokens to bridge bank contract
      await this.token.approve(this.bridgeBank.address, this.amount, {
        from: userOne
      }).should.be.fulfilled;

      // mint user tokens
      await this.token.mint(userOne, 2000, {
        from: operator
      }).should.be.fulfilled;

      await expectRevert(
        this.bridgeBank.lock(
          this.cosmosRecipient,
          this.token.address,
          100,
          {
            from: userOne,
            value: 0
          }
        ),
        "Amount being transferred is over the limit"
      );

    });
  });

  describe("Bridge claim status", function () {
    beforeEach(async function () {
      // Set up ProphecyClaim values
      this.cosmosSender = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      this.cosmosSenderSequence = 1;
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "TEST";
      this.actualSymbol = "eTEST"
      this.token = await BridgeToken.new(this.actualSymbol);
      this.amount = 100;

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [
        operator,
        consensusThreshold,
        this.initialValidators,
        this.initialPowers
      ],
        {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        this.cosmosBridge.address,
        operator,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );

      // Operator sets Bridge Bank
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      });

      // Update the lock/burn limit for this token
      await this.bridgeBank.updateTokenLockBurnLimit(this.token.address, this.amount, {
        from: operator
      }).should.be.fulfilled;
      await this.token.addMinter(this.bridgeBank.address);

    });

    it("should allow users to check if a prophecy claim is currently active", async function () {
      // Create the prophecy claim
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
      const status = await this.cosmosBridge.getProphecyThreshold(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active. False means it has not been 100% confirmed yet
      (status['0']).should.be.equal(false);
    });

    it("should allow us to check the cost of submitting a prophecy claim", async function () {
      // Create the prophecy claim
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        ++this.cosmosSenderSequence,
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
      const status = await this.cosmosBridge.getProphecyThreshold(prophecyClaimCount);

      // Bridge claim should be active
      (status[0]).should.be.equal(false);
    });

    it("should revert when a prophecy is resubmitted after payout", async function () {
      // Create the ProphecyClaim

      for (let i = 0; i < this.initialValidators.length - 1; i++) {
        await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_LOCK,
          this.cosmosSender,
          this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol.toLowerCase(),
          this.amount,
          {
            from: this.initialValidators[i]
          }
        );
      }

      await expectRevert(
        this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_LOCK,
          this.cosmosSender,
          this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol.toLowerCase(),
          this.amount,
          {
            from: this.initialValidators[ (this.initialValidators.length - 1) ]
          }
        ),
        "prophecyCompleted"
      );
      const claimID = (await this.cosmosBridge.getProphecyID(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol.toLowerCase(),
        this.amount,
      )).toString();

      const status = await this.cosmosBridge.getProphecyThreshold(claimID);

      // Bridge claim should be finished
      (status[0]).should.be.equal(true);
    });
  });
});
