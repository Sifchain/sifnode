const Web3Utils = require("web3-utils");
const web3 = require("web3");
const BigNumber = web3.BigNumber;

const { ethers, upgrades } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");
const { multiTokenSetup } = require("./helpers/testFixture");

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
  let CosmosBridge;
  let BridgeBank;
  let BridgeToken;
  // track the state of the deployed contracts
  let state;

  before(async function() {
    CosmosBridge = await ethers.getContractFactory("CosmosBridge");
    BridgeBank = await ethers.getContractFactory("BridgeBank");
    BridgeToken = await ethers.getContractFactory("BridgeToken");
    
    accounts = await ethers.getSigners();
    
    signerAccounts = accounts.map((e) => { return e.address });

    operator = accounts[0].address;
    userOne = accounts[1];
    userTwo = accounts[2];
    userFour = accounts[3];
    userThree = accounts[7].address;

    owner = accounts[5].address;
    pauser = accounts[6].address;

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);
  });

  beforeEach(async function () {
    state = await multiTokenSetup(
      initialValidators,
      initialPowers,
      operator,
      consensusThreshold,
      owner,
      userOne,
      userThree
    );
  });

  describe("BridgeBank", function () {
    it("should allow user to lock ERC20 tokens", async function () {
      await state.token1.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      // Attempt to lock tokens
      await state.bridgeBank.connect(userOne).lock(
        state.sender,
        state.token1.address,
        state.amount
      );

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await state.token1.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount);
    });

    it("should allow user to multi-lock ERC20 tokens", async function () {
      await state.token1.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      await state.token2.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      await state.token3.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      // Attempt to lock tokens
      await state.bridgeBank.connect(userOne).multiLock(
        [state.sender, state.sender, state.sender],
        [state.token1.address,state.token2.address,state.token3.address],
        [state.amount, state.amount, state.amount]
      );

      // Confirm that the user has been minted the correct token
      let afterUserBalance = Number(
        await state.token1.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount);

      afterUserBalance = Number(
        await state.token2.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount);

      afterUserBalance = Number(
        await state.token3.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount);
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
        Web3Utils.fromWei((+state.weiAmount).toString(), "ether")
      );
    });
  });
});
