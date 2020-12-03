const { deployProxy } = require('@openzeppelin/truffle-upgrades');

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

contract("CosmosBridge", function (accounts) {
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
      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.valset = await deployProxy(Valset, [
        operator,
        this.initialValidators,
        this.initialPowers

      ],
      {unsafeAllowCustomTypes: true}
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [operator, this.valset.address],
        {unsafeAllowCustomTypes: true});

      // Deploy Oracle contract
      this.oracle = await deployProxy(Oracle, [
        operator,
        this.valset.address,
        this.cosmosBridge.address,
        consensusThreshold
      ],
      {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
        operator,
        this.oracle.address,
        this.cosmosBridge.address,
        operator
      ],
      {unsafeAllowCustomTypes: true}
      );
    });

    it("should deploy the CosmosBridge with the correct parameters", async function () {
      this.cosmosBridge.should.exist;

      const claimCount = await this.cosmosBridge.prophecyClaimCount();
      Number(claimCount).should.be.bignumber.equal(0);

      const cosmosBridgeValset = await this.cosmosBridge.valset();
      cosmosBridgeValset.should.be.equal(this.valset.address);
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
      this.symbol = "ETH";
      this.nativeCosmosAssetDenom = "ATOM";
      this.prefixedNativeCosmosAssetDenom = "eATOM";
      this.amountWei = 100;
      this.amountNativeCosmos = 815;

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.valset = await deployProxy(Valset, [
        operator,
        this.initialValidators,
        this.initialPowers
      ],
      {unsafeAllowCustomTypes: true}
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await deployProxy(CosmosBridge, [operator, this.valset.address],
        {unsafeAllowCustomTypes: true});

      // Deploy Oracle contract
      this.oracle = await deployProxy(Oracle, [
        operator,
        this.valset.address,
        this.cosmosBridge.address,
        consensusThreshold
      ],
      {unsafeAllowCustomTypes: true}
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await deployProxy(BridgeBank, [
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
      contractBalanceWei.should.be.bignumber.equal(this.amountWei);

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
        this.nativeCosmosAssetDenom,
        this.amountNativeCosmos,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.nativeCosmosAssetDenom,
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
        this.nativeCosmosAssetDenom,
        this.amountNativeCosmos,
        {
          from: userThree
        }
      ).should.be.fulfilled;

      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const claimProphecyId = Number(event.args._prophecyID);
      const claimCosmosSender = event.args._cosmosSender;
      const claimEthereumReceiver = event.args._ethereumReceiver;

      // Check that the bridge token is a controlled bridge token
      const bridgeTokenAddr = await this.bridgeBank.getBridgeToken(
        this.prefixedNativeCosmosAssetDenom
      );
      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      this.bridgeToken = await BridgeToken.at(bridgeTokenAddr);

      const postRecipientBalance = bigInt(
        String(await this.bridgeToken.balanceOf(claimEthereumReceiver))
      );

      var expectedBalance = bigInt(String(priorRecipientBalance)).plus(
        String(this.amountNativeCosmos)
      );

      const receivedFunds = expectedBalance.equals(postRecipientBalance);
      receivedFunds.should.be.equal(true);

      // --------------------------------------------------------
      //  Now we'll do a 2nd lock prophecy claim of the native cosmos asset
      // --------------------------------------------------------
      console.log("\t[Attempt lock -> mint] (existing)");

      const { logs: logs2 } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        ++this.cosmosSenderSequence,
        this.ethereumReceiver,
        this.nativeCosmosAssetDenom,
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
        this.nativeCosmosAssetDenom,
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
        this.nativeCosmosAssetDenom,
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
