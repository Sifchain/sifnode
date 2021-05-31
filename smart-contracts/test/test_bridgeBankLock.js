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
  // track the state of the deployed contracts
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


  describe("Multi Lock ERC20 Tokens", function () {
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

    it("should allow user to multi-lock ERC20 tokens with multiLockBurn method", async function () {
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
      await state.bridgeBank.connect(userOne).multiLockBurn(
        [state.sender, state.sender, state.sender],
        [state.token1.address,state.token2.address,state.token3.address],
        [state.amount, state.amount, state.amount],
        [false, false, false]
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

    it("should allow user to multi-lock and burn ERC20 tokens and rowan with multiLockBurn method", async function () {
      await state.token1.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      await state.token2.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      // approve bridgebank to spend rowan
      await state.rowan.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      // Attempt to lock tokens
      await state.bridgeBank.connect(userOne).multiLockBurn(
        [state.sender, state.sender, state.sender],
        [state.token1.address, state.token2.address, state.rowan.address],
        [state.amount, state.amount, state.amount],
        [false, false, true]
      );

      // Confirm that the user has the proper balance after the multiLockBurn
      let afterUserBalance = Number(
        await state.token1.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount);

      afterUserBalance = Number(
        await state.token2.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount);

      afterUserBalance = Number(
        await state.rowan.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount);
    });

    it("should not allow user to multi-lock ERC20 tokens if one token is not fully approved", async function () {
      await state.token2.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      await state.token3.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLock(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount]
        )
      ).to.be.revertedWith("transfer amount exceeds allowance");

      // Confirm that user token balances have stayed the same
      let afterUserBalance = Number(
        await state.token1.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount * 2);

      afterUserBalance = Number(
        await state.token2.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount * 2);

      afterUserBalance = Number(
        await state.token3.balanceOf(userOne.address)
      );
      afterUserBalance.should.be.bignumber.equal(state.amount * 2);
    });

    it("should not allow user to multi-lock when parameters are malformed, not enough token amounts", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLock(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount]
        )
      ).to.be.revertedWith("M_P");
    });

    it("should not allow user to multi-lock when parameters are malformed, not enough token addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLock(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address],
          [state.amount, state.amount, state.amount]
        )
      ).to.be.revertedWith("M_P");
    });

    it("should not allow user to multi-lock when parameters are malformed, not enough sif addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLock(
          [state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount]
        )
      ).to.be.revertedWith("M_P");
    });

    it("should not allow user to multi-lock when parameters are malformed, invalid sif addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLock(
          [state.sender + "ee", state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount]
        )
      ).to.be.revertedWith("INV_ADR");
    });
  });

  describe("Multi Lock Burn ERC20 Tokens", function () {
    it("should revert when parameters are malformed, not enough token amounts", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount],
          [false, false, false],

        )
      ).to.be.revertedWith("M_P");
    });

    it("should revert when multi-lock parameters are malformed, not enough token addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address],
          [state.amount, state.amount, state.amount],
          [false, false, false]
        )
      ).to.be.revertedWith("M_P");
    });

    it("should revert when multi-lock parameters are malformed, not enough sif addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount],
          [false, false, false]
        )
      ).to.be.revertedWith("M_P");
    });

    it("should revert when multi-lock parameters are malformed, invalid sif addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender + "ee", state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount],
          [false, false, false]
        )
      ).to.be.revertedWith("INV_ADR");
    });
  });
});
