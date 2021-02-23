const { deployProxy, upgradeProxy, silenceWarnings } = require('@openzeppelin/truffle-upgrades');
const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");
const MockCosmosBridgeUpgrade = artifacts.require("MockCosmosBridgeUpgrade");

const { expectRevert } = require('@openzeppelin/test-helpers');

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("CosmosBridge Upgrade", function (accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];
  const userFour = accounts[4];

  // Consensus threshold of 70%
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

      this.cosmosBridge = await upgradeProxy(
          this.cosmosBridge.address,
          MockCosmosBridgeUpgrade,
          {unsafeAllowCustomTypes: true}
      )
    });

    it("should be able to mint tokens for a user", async function () {
      const amount = 100000000000;
      this.cosmosBridge.should.exist;
  
      await this.cosmosBridge.tokenFaucet({ from: operator});
      const operatorBalance = await this.cosmosBridge.balanceOf(operator);
      Number(operatorBalance).should.be.bignumber.equal(amount);
    });
    
    it("should be able to transfer tokens from the operator", async function () {
      const startingOperatorBalance = await this.cosmosBridge.balanceOf(operator);
      Number(startingOperatorBalance).should.be.bignumber.equal(0);

      const amount = 100000000000;
      this.cosmosBridge.should.exist;

      await this.cosmosBridge.tokenFaucet({ from: operator});

      await this.cosmosBridge.transfer(userOne, amount, { from: operator});
      const operatorBalance = await this.cosmosBridge.balanceOf(operator);
      const userOneBalance = await this.cosmosBridge.balanceOf(userOne);
      
      Number(operatorBalance).should.be.bignumber.equal(0);
      Number(userOneBalance).should.be.bignumber.equal(amount);
    });

    it("should not be able to initialize cosmos bridge a second time", async function () {
      this.cosmosBridge.should.exist;

      await expectRevert(
        this.cosmosBridge.initialize(userFour, 50, this.initialValidators, this.initialPowers),
        "Initialized"
      )
    });

    describe("CosmosBridge has all previous functionality", function () {

      it("should allow the operator to set the Bridge Bank", async function () {
        this.bridgeBank.should.exist;

        await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
          from: operator
        }).should.be.fulfilled;

        const bridgeBank = await this.cosmosBridge.bridgeBank();
        bridgeBank.should.be.equal(this.bridgeBank.address);
      });

      it("should not allow the operator to update the Bridge Bank once it has been set", async function () {
        await this.cosmosBridge.setBridgeBank(operator, {
          from: operator
        }).should.be.fulfilled;

        await this.cosmosBridge
        .setBridgeBank(operator, {
          from: operator
        })
        .should.be.rejectedWith(EVMRevert);
      });
    });
  });
});