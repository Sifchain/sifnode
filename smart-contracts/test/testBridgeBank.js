const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { setup, getValidClaim } = require("./helpers/testFixture");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

use(solidity);

const getBalance = async function(address) {
  return await network.provider.send("eth_getBalance", [address]);
}

describe("Test Bridge Bank", function () {
  let userOne;
  let userTwo;
  let userThree;
  let userFour;
  let accounts;
  let signerAccounts;
  let operator;
  let owner;
  const consensusThreshold = 75;
  let initialPowers;
  let initialValidators;
  let networkDescriptor;
  let state;

  before(async function() {
    accounts = await ethers.getSigners();
    
    signerAccounts = accounts.map((e) => { return e.address });

    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[7];

    owner = accounts[5];
    pauser = accounts[6];

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);

    networkDescriptor = 1;
  });

  beforeEach(async function () {
    state = await setup({
        initialValidators,
        initialPowers,
        operator,
        consensusThreshold,
        owner,
        user: userOne,
        recipient: userThree,
        pauser,
        networkDescriptor,
        lockTokensOnBridgeBank: true
    });
  });

  describe("BridgeBank single lock burn transactions", function () {
    it("should allow user to lock ERC20 tokens", async function () {
      const bridgeBankBalanceBefore = await state.token.balanceOf(state.bridgeBank.address);

      // approve and lock tokens
      await state.token.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      // Attempt to lock tokens
      await state.bridgeBank.connect(userOne).lock(
        state.sender,
        state.token.address,
        state.amount
      );

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await state.token.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(0);

      // check if BridgeBank now owns the tokens
      const bridgeBankBalanceAfter = await state.token.balanceOf(state.bridgeBank.address);
      const expectedBalance = Number(bridgeBankBalanceBefore) + Number(state.amount);
      expect(bridgeBankBalanceAfter).to.be.equal(expectedBalance);
    });

    it("should allow user to lock fake ERC20 tokens", async function () {
      const FakeTokenFactory = await ethers.getContractFactory("FakeERC20");
      const fakeToken = await FakeTokenFactory.deploy();

      // Approve and lock tokens
      await expect(state.bridgeBank.connect(userOne).lock(state.sender, fakeToken.address, state.amount))
        .to.emit(state.bridgeBank, 'LogLock')
        .withArgs(userOne.address, state.sender, fakeToken.address, state.amount, "3", 18, "", "", state.networkDescriptor);
    });

    it("should allow users to lock Ethereum in the bridge bank", async function () {
      const tx = await state.bridgeBank.connect(userOne).lock(
        state.sender,
        state.constants.zeroAddress,
        state.weiAmount, {
          value: state.weiAmount
        }
      ).should.be.fulfilled;
      await tx.wait();

      const contractBalanceWei = await getBalance(state.bridgeBank.address);
      const contractBalance = Web3Utils.fromWei(contractBalanceWei, "ether");

      contractBalance.should.be.bignumber.equal(
        Web3Utils.fromWei((+state.weiAmount + +state.amount).toString(), "ether")
      );
    });

    it("should not allow users to lock Ethereum in the bridge bank if the sent amount and amount param are different", async function () {
      await expect(
        state.bridgeBank.connect(userOne).lock(
          state.sender,
          state.constants.zeroAddress,
          state.weiAmount + 1, {
            value: state.weiAmount
          },
        ),
      ).to.be.revertedWith("amount mismatch");
    });

    it("should not allow users to lock Ethereum in the bridge bank if sending tokens", async function () {
      await expect(
        state.bridgeBank.connect(userOne).lock(
          state.sender,
          state.token.address,
          state.weiAmount + 1, {
            value: state.weiAmount
          },
        ),
      ).to.be.revertedWith("INV_NATIVE_SEND");
    });
  });

  describe("BridgeBank single lock burn transactions", function () {
    it("should allow a user to burn tokens from the bridge bank", async function () {
      const BridgeToken = await ethers.getContractFactory("BridgeToken");
      const bridgeToken = await BridgeToken.deploy("rowan", "rowan", 18, state.constants.denom.rowan);

      await bridgeToken.connect(operator).grantRole(state.constants.roles.minter, operator.address)
      await bridgeToken.connect(operator).mint(userOne.address, state.amount);
      await bridgeToken.connect(userOne).approve(state.bridgeBank.address, state.amount);
      await state.bridgeBank.connect(owner).addExistingBridgeToken(bridgeToken.address);
  
      await state.bridgeBank.connect(userOne).burn(
        state.sender,
        bridgeToken.address,
        state.amount
      );

      const afterUserBalance = Number(
        await bridgeToken.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(0);
    });

    it("should allow a user to burn tokens twice from the bridge bank", async function () {
      const BridgeTokenFactory = await ethers.getContractFactory("BridgeToken");
      const bridgeToken = await BridgeTokenFactory.deploy("rowan", "rowan", 18, state.constants.denom.rowan);

      const doubleAmount = Number(state.amount) * 2;

      await bridgeToken.connect(operator).grantRole(state.constants.roles.minter, operator.address)
      await bridgeToken.connect(operator).mint(userOne.address, doubleAmount);
      await bridgeToken.connect(userOne).approve(state.bridgeBank.address, doubleAmount);
      await state.bridgeBank.connect(owner).addExistingBridgeToken(bridgeToken.address);
  
      await state.bridgeBank.connect(userOne).burn(
        state.sender,
        bridgeToken.address,
        state.amount
      );

      let afterUserBalance = Number(await bridgeToken.balanceOf(userOne.address));
      afterUserBalance.should.be.bignumber.equal(state.amount);

      // Do it again
      await state.bridgeBank.connect(userOne).burn(
        state.sender,
        bridgeToken.address,
        state.amount
      );

      afterUserBalance = Number(await bridgeToken.balanceOf(userOne.address));
      afterUserBalance.should.be.bignumber.equal(0);
    });

    it("should allow a user to burn a token that doesn't have a denom", async function () {
      const oldTokenFactory = await ethers.getContractFactory("Erowan");
      const oldToken = await oldTokenFactory.deploy("OLD");

      await oldToken.mint(userOne.address, state.amount);
      await oldToken.connect(userOne).approve(state.bridgeBank.address, state.amount);
      await state.bridgeBank.connect(owner).addExistingBridgeToken(oldToken.address);
  
      await state.bridgeBank.connect(userOne).burn(
        state.sender,
        oldToken.address,
        state.amount
      );

      let afterUserBalance = Number(await oldToken.balanceOf(userOne.address));
      afterUserBalance.should.be.bignumber.equal(0);
    });
  });

  describe("BridgeBank administration of Bridgetokens", function () {
    it("should allow the operator to set a BridgeToken's denom", async function () {
      // expect the token to NOT have a defined denom on BridgeBank
      let registeredDenom = await state.bridgeBank.contractDenom(state.rowan.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.none);

      // expect the token itself to have a denom
      let registeredDenomInBridgeToken = await state.rowan.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.rowan);

      // set a new denom
      await expect(state.bridgeBank.connect(owner)
        .setBridgeTokenDenom(state.rowan.address, state.constants.denom.one))
        .to.be.fulfilled;

      // check the denom saved on BridgeBank
      registeredDenom = await state.bridgeBank.contractDenom(state.rowan.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.one);

      // check the denom saved on the BridgeToken itself
      registeredDenomInBridgeToken = await state.rowan.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.one);
    });

    it("should not allow a user to set a BridgeToken's denom", async function () {
      // set a new denom
      await expect(state.bridgeBank.connect(userOne)
        .setBridgeTokenDenom(state.rowan.address, state.constants.denom.one))
        .to.be.revertedWith('!owner');
    });

    it("should revert when trying to set the denom of the old Erowan token", async function () {
      // Deploy the old Erowan token
      const erowanTokenFactory = await ethers.getContractFactory('Erowan');
      const erowanToken = await erowanTokenFactory.deploy('erowan');
      await erowanToken.deployed();

      // expect the token to NOT have a defined denom on BridgeBank
      let registeredDenom = await state.bridgeBank.contractDenom(erowanToken.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.none);

      // try to set a new denom
      await expect(state.bridgeBank.connect(owner)
        .setBridgeTokenDenom(erowanToken.address, state.constants.denom.one))
        .to.be.rejectedWith("Transaction reverted: function selector was not recognized and there's no fallback function");

      // check if the denom was saved on BridgeBank (should not)
      registeredDenom = await state.bridgeBank.contractDenom(erowanToken.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.none);
    });

    it("should allow the owner to set many BridgeTokens' denom in a batch", async function () {
      // expect rowan to NOT have a defined denom on BridgeBank
      let registeredDenom = await state.bridgeBank.contractDenom(state.rowan.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.none);

      // expect rowan itself to have a denom
      let registeredDenomInBridgeToken = await state.rowan.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.rowan);

      // transfer ownership of state.token_noDenom to the BridgeBank
      await state.token_noDenom.connect(operator).grantRole(state.constants.roles.admin, state.bridgeBank.address);

      // expect the noDenom token to NOT have a defined denom on BridgeBank
      let registeredDenom2 = await state.bridgeBank.contractDenom(state.token_noDenom.address);
      expect(registeredDenom2).to.be.equal(state.constants.denom.none);

      // expect the noDenom token itself to NOT have a denom either
      let registeredDenomInBridgeToken2 = await state.token_noDenom.cosmosDenom();
      expect(registeredDenomInBridgeToken2).to.be.equal(state.constants.denom.none);

      // set the new denom for both of them
      await expect(state.bridgeBank.connect(owner)
        .batchSetBridgeTokenDenom(
          [state.rowan.address, state.token_noDenom.address],
          [state.constants.denom.one, state.constants.denom.two]
        )).to.be.fulfilled;

      // check the denom saved on BridgeBank
      registeredDenom = await state.bridgeBank.contractDenom(state.rowan.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.one);

      // check the denom saved on Rowan itself
      registeredDenomInBridgeToken = await state.rowan.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.one);

      // check the denom saved on BridgeBank
      registeredDenom2 = await state.bridgeBank.contractDenom(state.token_noDenom.address);
      expect(registeredDenom2).to.be.equal(state.constants.denom.two);

      // check the denom saved on the noDenom BridgeToken itself
      registeredDenomInBridgeToken2 = await state.token_noDenom.cosmosDenom();
      expect(registeredDenomInBridgeToken2).to.be.equal(state.constants.denom.two);
    });

    it("should NOT allow a user to set many BridgeTokens' denom in a batch", async function () {
      // expect rowan to NOT have a defined denom on BridgeBank
      let registeredDenom = await state.bridgeBank.contractDenom(state.rowan.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.none);

      // expect rowan itself to have a denom
      let registeredDenomInBridgeToken = await state.rowan.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.rowan);

      // transfer ownership of state.token_noDenom to the BridgeBank
      await state.token_noDenom.connect(operator).grantRole(state.constants.roles.admin, state.bridgeBank.address);

      // expect the noDenom token to NOT have a defined denom on BridgeBank
      let registeredDenom2 = await state.bridgeBank.contractDenom(state.token_noDenom.address);
      expect(registeredDenom2).to.be.equal(state.constants.denom.none);

      // expect the noDenom token itself to NOT have a denom either
      let registeredDenomInBridgeToken2 = await state.token_noDenom.cosmosDenom();
      expect(registeredDenomInBridgeToken2).to.be.equal(state.constants.denom.none);

      // try to set the new denom for both of them
      await expect(state.bridgeBank.connect(userOne)
        .batchSetBridgeTokenDenom(
          [state.rowan.address, state.token_noDenom.address],
          [state.constants.denom.one, state.constants.denom.two]
        )).to.be.revertedWith("!owner");

      // check the denom saved on BridgeBank (shouldn't have changed)
      registeredDenom = await state.bridgeBank.contractDenom(state.rowan.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.none);

      // check the denom saved on Rowan itself (shouldn't have changed)
      registeredDenomInBridgeToken = await state.rowan.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.rowan);

      // check the denom saved on BridgeBank (shouldn't have changed)
      registeredDenom2 = await state.bridgeBank.contractDenom(state.token_noDenom.address);
      expect(registeredDenom2).to.be.equal(state.constants.denom.none);

      // check the denom saved on the noDenom BridgeToken itself (shouldn't have changed)
      registeredDenomInBridgeToken2 = await state.token_noDenom.cosmosDenom();
      expect(registeredDenomInBridgeToken2).to.be.equal(state.constants.denom.none);
    });

    it("should allow the operator to add many BridgeTokens in a batch", async function () {
      // expect token1 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist1 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token1.address);
      expect(isInCosmosWhitelist1).to.be.false;

      // expect token2 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist2 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token2.address);
      expect(isInCosmosWhitelist2).to.be.false;

      // expect token3 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist3 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token3.address);
      expect(isInCosmosWhitelist3).to.be.false;

      // add tokens as BridgeTokens
      await expect(state.bridgeBank.connect(owner)
        .batchAddExistingBridgeTokens([
          state.token1.address,
          state.token2.address,
          state.token3.address
        ])
      ).to.be.fulfilled;

      // check if the tokens are now correctly registered
      // expect token1 to be registered as a BridgeToken
      isInCosmosWhitelist1 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token1.address);
      expect(isInCosmosWhitelist1).to.be.true;

      // expect token2 to be registered as a BridgeToken
      isInCosmosWhitelist2 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token2.address);
      expect(isInCosmosWhitelist2).to.be.true;

      // expect token3 to be registered as a BridgeToken
      isInCosmosWhitelist3 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token3.address);
      expect(isInCosmosWhitelist3).to.be.true;
    });

    it("should allow the owner to add many BridgeTokens in a batch and then set the cosmosDenom", async function () {
      // add bridgebank as admin of the tokens
      await state.token1.connect(state.operator).grantRole(state.constants.roles.admin, state.bridgeBank.address);
      await state.token2.connect(state.operator).grantRole(state.constants.roles.admin, state.bridgeBank.address);
      await state.token3.connect(state.operator).grantRole(state.constants.roles.admin, state.bridgeBank.address);

      // expect token1 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist1 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token1.address);
      expect(isInCosmosWhitelist1).to.be.false;

      // expect token2 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist2 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token2.address);
      expect(isInCosmosWhitelist2).to.be.false;

      // expect token3 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist3 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token3.address);
      expect(isInCosmosWhitelist3).to.be.false;

      // add tokens as BridgeTokens
      await expect(state.bridgeBank.connect(owner)
        .batchAddExistingBridgeTokens([
          state.token1.address,
          state.token2.address,
          state.token3.address
        ])
      ).to.be.fulfilled;

      // check if the tokens are now correctly registered
      // expect token1 to be registered as a BridgeToken
      isInCosmosWhitelist1 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token1.address);
      expect(isInCosmosWhitelist1).to.be.true;

      // expect token2 to be registered as a BridgeToken
      isInCosmosWhitelist2 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token2.address);
      expect(isInCosmosWhitelist2).to.be.true;

      // expect token3 to be registered as a BridgeToken
      isInCosmosWhitelist3 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token3.address);
      expect(isInCosmosWhitelist3).to.be.true;

      // Check the current token denoms in each token:
      let registeredDenomInBridgeToken = await state.token1.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.two);

      let registeredDenomInBridgeToken2 = await state.token2.cosmosDenom();
      expect(registeredDenomInBridgeToken2).to.be.equal(state.constants.denom.three);

      let registeredDenomInBridgeToken3 = await state.token3.cosmosDenom();
      expect(registeredDenomInBridgeToken3).to.be.equal(state.constants.denom.four);

      // Now, set the denom for all those tokens
      await expect(state.bridgeBank.connect(owner)
        .batchSetBridgeTokenDenom(
          [state.token1.address, state.token2.address, state.token3.address],
          [state.constants.denom.one, state.constants.denom.two, state.constants.denom.three]
        )).to.be.fulfilled;

      // check the denom saved on BridgeBank
      registeredDenom = await state.bridgeBank.contractDenom(state.token1.address);
      expect(registeredDenom).to.be.equal(state.constants.denom.one);

      // check the denom saved on token1 itself
      registeredDenomInBridgeToken = await state.token1.cosmosDenom();
      expect(registeredDenomInBridgeToken).to.be.equal(state.constants.denom.one);

      // check the denom saved on BridgeBank
      registeredDenom2 = await state.bridgeBank.contractDenom(state.token2.address);
      expect(registeredDenom2).to.be.equal(state.constants.denom.two);

      // check the denom saved on token2 itself
      registeredDenomInBridgeToken2 = await state.token2.cosmosDenom();
      expect(registeredDenomInBridgeToken2).to.be.equal(state.constants.denom.two);

      // check the denom saved on BridgeBank
      registeredDenom3 = await state.bridgeBank.contractDenom(state.token3.address);
      expect(registeredDenom3).to.be.equal(state.constants.denom.three);

      // check the denom saved on token3 itself
      registeredDenomInBridgeToken3 = await state.token3.cosmosDenom();
      expect(registeredDenomInBridgeToken3).to.be.equal(state.constants.denom.three);
    });

    it("should NOT allow a user to add many BridgeTokens in a batch", async function () {
      // expect token1 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist1 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token1.address);
      expect(isInCosmosWhitelist1).to.be.false;

      // expect token2 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist2 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token2.address);
      expect(isInCosmosWhitelist2).to.be.false;

      // expect token3 to NOT be registered as a BridgeToken
      let isInCosmosWhitelist3 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token3.address);
      expect(isInCosmosWhitelist3).to.be.false;

      // add tokens as BridgeTokens
      await expect(state.bridgeBank.connect(userOne)
        .batchAddExistingBridgeTokens([
          state.token1.address,
          state.token2.address,
          state.token3.address
        ])
      ).to.be.revertedWith('!owner');

      // check if the tokens are now registered (should not be)
      // expect token1 to NOT be registered as a BridgeToken
      isInCosmosWhitelist1 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token1.address);
      expect(isInCosmosWhitelist1).to.be.false;

      // expect token2 to NOT be registered as a BridgeToken
      isInCosmosWhitelist2 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token2.address);
      expect(isInCosmosWhitelist2).to.be.false;

      // expect token3 to NOT be registered as a BridgeToken
      isInCosmosWhitelist3 = await state.bridgeBank.getCosmosTokenInWhiteList(state.token3.address);
      expect(isInCosmosWhitelist3).to.be.false;
    });

    it("should allow anyone to forceSetBridgeTokenDenom", async function () {
      // expect token2's denom to NOT be registered in BridgeBank
      let denomInBridgeBank = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank).to.be.equal("");

      // add token 2 as BridgeToken
      await state.bridgeBank.connect(owner).addExistingBridgeToken(state.token2.address);

      // userOne calls forceSetBridgeTokenDenom
      await expect(state.bridgeBank.connect(userOne)
        .forceSetBridgeTokenDenom(state.token2.address))
        .to.be.fulfilled;

      // expect token2's denom to be registered in BridgeBank
      denomInBridgeBank = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank).to.be.equal(state.constants.denom.three);
    });

    it("should fail to call forceSetBridgeTokenDenom for non-cosmosWhitelisted tokens", async function () {
      // expect token2's denom to NOT be registered in BridgeBank
      let denomInBridgeBank = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank).to.be.equal("");
      
      // userOne calls forceSetBridgeTokenDenom
      await expect(state.bridgeBank.connect(userOne)
        .forceSetBridgeTokenDenom(state.token2.address))
        .to.be.rejectedWith('Token is not in Cosmos whitelist');

      // expect token2's denom to NOT be registered in BridgeBank
      denomInBridgeBank = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank).to.be.equal("");
    });

    it("should allow anyone to batchForceSetBridgeTokenDenom", async function () {
      // expect token2's denom to NOT be registered in BridgeBank
      let denomInBridgeBank2 = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank2).to.be.equal("");

      // expect token3's denom to NOT be registered in BridgeBank
      let denomInBridgeBank3 = await state.bridgeBank.contractDenom(state.token3.address);
      expect(denomInBridgeBank3).to.be.equal("");

      // add tokens as BridgeTokens
      await state.bridgeBank.connect(owner)
        .batchAddExistingBridgeTokens([state.token2.address, state.token3.address]);

      // userOne calls batchForceSetBridgeTokenDenom
      await expect(state.bridgeBank.connect(userOne)
        .batchForceSetBridgeTokenDenom([state.token2.address, state.token3.address]))
        .to.be.fulfilled;

      // expect token2's denom to be registered in BridgeBank
      denomInBridgeBank2 = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank2).to.be.equal(state.constants.denom.three);

      // expect token3's denom to be registered in BridgeBank
      denomInBridgeBank3 = await state.bridgeBank.contractDenom(state.token3.address);
      expect(denomInBridgeBank3).to.be.equal(state.constants.denom.four);
    });

    it("should fail to call batchForceSetBridgeTokenDenom for non-cosmosWhitelisted tokens", async function () {
      // expect token2's denom to NOT be registered in BridgeBank
      let denomInBridgeBank2 = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank2).to.be.equal("");

      // expect token3's denom to NOT be registered in BridgeBank
      let denomInBridgeBank3 = await state.bridgeBank.contractDenom(state.token3.address);
      expect(denomInBridgeBank3).to.be.equal("");

      // add token 2 as BridgeToken, BUT NOT TOKEN 3
      await state.bridgeBank.connect(owner).addExistingBridgeToken(state.token2.address);
      
      // userOne calls forceSetBridgeTokenDenom
      await expect(state.bridgeBank.connect(userOne)
        .batchForceSetBridgeTokenDenom([state.token2.address, state.token3.address]))
        .to.be.rejectedWith('Token is not in Cosmos whitelist');

      // expect token2's denom to NOT be registered in BridgeBank
      denomInBridgeBank2 = await state.bridgeBank.contractDenom(state.token2.address);
      expect(denomInBridgeBank2).to.be.equal("");

      // expect token3's denom to NOT be registered in BridgeBank
      denomInBridgeBank3 = await state.bridgeBank.contractDenom(state.token3.address);
      expect(denomInBridgeBank3).to.be.equal("");
    });
  });
});
