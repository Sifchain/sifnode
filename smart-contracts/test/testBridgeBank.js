const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { singleSetup } = require("./helpers/testFixture");

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
    userThree = accounts[7].address;

    owner = accounts[5];
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);

    networkDescriptor = 1;
  });

  beforeEach(async function () {
    state = await singleSetup(
        initialValidators,
        initialPowers,
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser,
        networkDescriptor
    );
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
        state.ethereumToken,
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
          state.ethereumToken,
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
      const bridgeToken = await BridgeToken.deploy("rowan", "rowan", 18);

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
});
