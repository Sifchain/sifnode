const { deployProxy, silenceWarnings } = require('@openzeppelin/truffle-upgrades');

const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");
const BridgeToken = artifacts.require("BridgeToken");

var bigInt = require("big-integer");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(web3.BigNumber))
  .should();

const {
  expectRevert, // Assertions for transactions that should fail
} = require('@openzeppelin/test-helpers');
const { expect } = require('chai');

contract("End To End", function (accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];
  const userFour = accounts[4];

  // User account
  const userSeven = accounts[7];

  // Contract's enum ClaimType can be represented a sequence of integers
  const CLAIM_TYPE_BURN = 1;
  const CLAIM_TYPE_LOCK = 2;

  // Consensus threshold
  const consensusThreshold = 70;

  describe("CosmosBridge smart contract deployment", function () {
    beforeEach(async function () {
      await silenceWarnings();
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
    });
  });

  describe("Claim flow", function () {
    beforeEach(async function () {
      // Set up ProphecyClaim values
      this.cosmosSender = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      this.cosmosSenderSequence = 1;
      this.ethereumReceiver = userSeven;
      this.ethTokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "eth";
      this.nativeCosmosAssetDenom = "ATOM";
      this.prefixedNativeCosmosAssetDenom = "eATOM";
      this.amountWei = 100;
      this.amountNativeCosmos = 815;

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.secondValidators = [userOne, userTwo];
      this.secondPowers = [50, 50];
      this.thirdValidators = [userThree, userFour];
      this.thirdPowers = [50, 50];

      this.symbol.token
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
    });

    it("Burn prophecy claim flow", async function () {
      console.log("\t[Attempt burn -> unlock]");

      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await this.bridgeBank.lock(
        this.cosmosSender,
        this.ethTokenAddress,
        this.amountWei,
        {
          from: userOne,
          value: this.amountWei
        }
      ).should.be.fulfilled;

      const contractBalanceWei = await web3.eth.getBalance(
        this.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei).should.be.equal(this.amountWei);

      // --------------------------------------------------------
      //  Check receiver's account balance prior to the claims
      // --------------------------------------------------------
      const priorRecipientBalance = await web3.eth.getBalance(userSeven);

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amountWei,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amountWei,
        {
          from: userTwo
        }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.symbol,
        this.amountWei,
        {
          from: userFour
        }
      ).should.be.fulfilled;

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      const postRecipientBalance = bigInt(
        String(await web3.eth.getBalance(userSeven))
      );

      var expectedBalance = bigInt(String(priorRecipientBalance)).plus(
        String(this.amountWei)
      );

      const receivedFunds = expectedBalance.equals(postRecipientBalance);
      receivedFunds.should.be.equal(true);

      // Fail to create prophecy claim if from non validator
      await expectRevert(
          this.cosmosBridge.newProphecyClaim(
              CLAIM_TYPE_BURN,
              this.cosmosSender,
              this.cosmosSenderSequence,
              this.ethereumReceiver,
              this.symbol,
              this.amountWei,
              {
                from: userSeven
              }
          ),
          "Must be an active validator"
      );

      // Also make sure everything runs twice.

      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await this.bridgeBank.lock(
          this.cosmosSender,
          this.ethTokenAddress,
          this.amountWei,
          {
            from: userOne,
            value: this.amountWei
          }
      ).should.be.fulfilled;

      const contractBalanceWei2 = await web3.eth.getBalance(
          this.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei2).should.be.equal(this.amountWei);

      // --------------------------------------------------------
      //  Check receiver's account balance prior to the claims
      // --------------------------------------------------------
      const priorRecipientBalance2 = await web3.eth.getBalance(userSeven);

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          ++this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amountWei,
          {
            from: userOne
          }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amountWei,
          {
            from: userTwo
          }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amountWei,
          {
            from: userFour
          }
      ).should.be.fulfilled;

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      const postRecipientBalance2 = bigInt(
          String(await web3.eth.getBalance(userSeven))
      );

      var expectedBalance2 = bigInt(String(priorRecipientBalance)).plus(
          String(this.amountWei)
      );

      const receivedFunds2 = expectedBalance.equals(postRecipientBalance);
      receivedFunds2.should.be.equal(true);

      // Also make sure everything runs third time after switching validators.

      // Operator resets the valset
      await this.cosmosBridge.updateValset(
          this.secondValidators,
          this.secondPowers,
          {
            from: operator
          }
      ).should.be.fulfilled;

      // Confirm that both initial validators are now active validators
      const isUserOneValidator = await this.cosmosBridge.isActiveValidator.call(
          userOne
      );
      isUserOneValidator.should.be.equal(true);
      const isUserTwoValidator = await this.cosmosBridge.isActiveValidator.call(
          userTwo
      );
      isUserTwoValidator.should.be.equal(true);

      // Confirm that all both secondary validators are not active validators
      const isUserThreeValidator = await this.cosmosBridge.isActiveValidator.call(
          userThree
      );
      isUserThreeValidator.should.be.equal(false);
      const isUserFourValidator = await this.cosmosBridge.isActiveValidator.call(
          userFour
      );
      isUserFourValidator.should.be.equal(false);

      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await this.bridgeBank.lock(
          this.cosmosSender,
          this.ethTokenAddress,
          this.amountWei,
          {
            from: userOne,
            value: this.amountWei
          }
      ).should.be.fulfilled;

      const contractBalanceWei3 = await web3.eth.getBalance(
          this.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei3).should.be.equal(this.amountWei);

      // --------------------------------------------------------
      //  Check receiver's account balance prior to the claims
      // --------------------------------------------------------
      const priorRecipientBalance3 = await web3.eth.getBalance(userSeven);

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          ++this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amountWei,
          {
            from: userOne
          }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amountWei,
          {
            from: userTwo
          }
      ).should.be.fulfilled;

      // Fail to create prophecy claim if from non validator
      await expectRevert(
          this.cosmosBridge.newProphecyClaim(
              CLAIM_TYPE_BURN,
              this.cosmosSender,
              this.cosmosSenderSequence,
              this.ethereumReceiver,
              this.symbol,
              this.amountWei,
              {
                from: userThree
              }
          ),
          "Must be an active validator"
      );

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      const postRecipientBalance3 = bigInt(
          String(await web3.eth.getBalance(userSeven))
      );

      var expectedBalance3 = bigInt(String(priorRecipientBalance)).plus(
          String(this.amountWei)
      );

      const receivedFunds3 = expectedBalance.equals(postRecipientBalance);
      receivedFunds3.should.be.equal(true);

      // Also make sure everything runs fourth time after switching validators a second time.

      // Operator resets the valset
      await this.cosmosBridge.updateValset(
          this.thirdValidators,
          this.thirdPowers,
          {
            from: operator
          }
      ).should.be.fulfilled;

      // Confirm that both initial validators are no longer an active validators
      const isUserOneValidator2 = await this.cosmosBridge.isActiveValidator.call(
          userOne
      );
      isUserOneValidator2.should.be.equal(false);
      const isUserTwoValidator2 = await this.cosmosBridge.isActiveValidator.call(
          userTwo
      );
      isUserTwoValidator2.should.be.equal(false);

      // Confirm that both secondary validators are now active validators
      const isUserThreeValidator2 = await this.cosmosBridge.isActiveValidator.call(
          userThree
      );
      isUserThreeValidator2.should.be.equal(true);
      const isUserFourValidator2 = await this.cosmosBridge.isActiveValidator.call(
          userFour
      );
      isUserFourValidator2.should.be.equal(true);

      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await this.bridgeBank.lock(
          this.cosmosSender,
          this.ethTokenAddress,
          this.amountWei,
          {
            from: userOne,
            value: this.amountWei
          }
      ).should.be.fulfilled;

      const contractBalanceWei4 = await web3.eth.getBalance(
          this.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      Number(contractBalanceWei4).should.be.equal(this.amountWei);

      // --------------------------------------------------------
      //  Check receiver's account balance prior to the claims
      // --------------------------------------------------------
      const priorRecipientBalance4 = await web3.eth.getBalance(userSeven);

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          ++this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amountWei,
          {
            from: userThree
          }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
          CLAIM_TYPE_BURN,
          this.cosmosSender,
          this.cosmosSenderSequence,
          this.ethereumReceiver,
          this.symbol,
          this.amountWei,
          {
            from: userFour
          }
      ).should.be.fulfilled;

      // Fail to create prophecy claim if from non validator
      await expectRevert(
          this.cosmosBridge.newProphecyClaim(
              CLAIM_TYPE_BURN,
              this.cosmosSender,
              this.cosmosSenderSequence,
              this.ethereumReceiver,
              this.symbol,
              this.amountWei,
              {
                from: userOne
              }
          ),
          "Must be an active validator"
      );

      const contractBalanceWeiAfter = await web3.eth.getBalance(
        this.bridgeBank.address
      );
      contractBalanceWeiAfter.toString().should.be.equal("0")
      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      // const postRecipientBalance4 = (await web3.eth.getBalance(userSeven)).toString();

      // var expectedBalance4 = bigInt(String(priorRecipientBalance)).plus(
      //     String(this.amountWei)
      // );

      // const receivedFunds4 = Number(expectedBalance4).should.be.equal(postRecipientBalance4);
      // receivedFunds4.should.be.equal(true);
    });

    it("Lock prophecy claim flow", async function () {
      console.log("\t[Attempt lock -> mint] (new)");
      const priorRecipientBalance = 0;

      // --------------------------------------------------------
      //  Create a new lock prophecy claim on cosmos bridge
      // --------------------------------------------------------
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        ++this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.prefixedNativeCosmosAssetDenom.toLowerCase(),
        this.amountNativeCosmos,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Check that the bridge token is a controlled bridge token
      const bridgeTokenAddr = await this.bridgeBank.getBridgeToken(
        this.prefixedNativeCosmosAssetDenom.toLowerCase()
      );

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      this.bridgeToken = await BridgeToken.at(bridgeTokenAddr);

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.prefixedNativeCosmosAssetDenom.toLowerCase(),
        this.amountNativeCosmos,
        {
          from: userTwo
        }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.prefixedNativeCosmosAssetDenom.toLowerCase(),
        this.amountNativeCosmos,
        {
          from: userThree
        }
      ).should.be.fulfilled;

      (await this.bridgeToken.balanceOf(this.ethereumReceiver)).toString().should.be.equal(this.amountNativeCosmos.toString());
      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const claimProphecyId = Number(event.args._prophecyID);
      const claimCosmosSender = event.args._cosmosSender;
      const claimEthereumReceiver = event.args._ethereumReceiver;


      const postRecipientBalance = bigInt(
        String(await this.bridgeToken.balanceOf(claimEthereumReceiver))
      );

      // var expectedBalance = bigInt(String(priorRecipientBalance)).plus(
      //   String(this.amountNativeCosmos)
      // );

      // const receivedFunds = expectedBalance.equals(postRecipientBalance);
      // receivedFunds.should.be.equal(true);

      // --------------------------------------------------------
      //  Now we'll do a 2nd lock prophecy claim of the native cosmos asset
      // --------------------------------------------------------
      console.log("\t[Attempt lock -> mint] (existing)");

      const { logs: logs2 } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        ++this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.prefixedNativeCosmosAssetDenom.toLowerCase(),
        this.amountNativeCosmos,
        {
          from: userTwo
        }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.prefixedNativeCosmosAssetDenom.toLowerCase(),
        this.amountNativeCosmos,
        {
          from: userThree
        }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.prefixedNativeCosmosAssetDenom.toLowerCase(),
        this.amountNativeCosmos,
        {
          from: userFour
        }
      ).should.be.fulfilled;

      const event2 = logs2.find(e => e.event === "LogNewProphecyClaim");
      const claimProphecyId2 = Number(event2.args._prophecyID);
      const claimCosmosSender2 = event2.args._cosmosSender;
      const claimEthereumReceiver2 = event2.args._ethereumReceiver;
      const claimTokenAddress2 = event2.args._tokenAddress;
      const claimAmount2 = Number(event2.args._amount);
      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------

      const postRecipientBalance2 = bigInt(
        String(await this.bridgeToken.balanceOf(claimEthereumReceiver2))
      );

      var expectedBalance2 = bigInt(String(postRecipientBalance)).plus(
        String(this.amountNativeCosmos)
      );

      const receivedFunds2 = expectedBalance2.equals(postRecipientBalance2);
      receivedFunds2.should.be.equal(true);
    });
  });
});
