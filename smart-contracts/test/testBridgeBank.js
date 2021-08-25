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
    });

    it("should not allow user to lock ERC20 tokens", async function () {
      const FakeToken = await ethers.getContractFactory("FakeERC20");
      fakeToken = await FakeToken.deploy();
      
      // Add the token into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(fakeToken.address, true)
        .should.be.fulfilled;

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
      await expect(state.bridgeBank.connect(operator)
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
        .to.be.revertedWith('!operator');
    });
  });
});
