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
  let networkDescriptor;
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
    pauser = accounts[6];

    initialPowers = [25, 25, 25, 25];
    initialValidators = signerAccounts.slice(0, 4);

    networkDescriptor = 1;
  });

  beforeEach(async function () {
    state = await multiTokenSetup(
      initialValidators,
      initialPowers,
      operator,
      consensusThreshold,
      owner,
      userOne,
      userThree,
      pauser.address,
      networkDescriptor
    );
  });

  describe("BridgeBank", function () {
    it("should deploy the BridgeBank, correctly setting the operator", async function () {
      state.bridgeBank.should.exist;

      const bridgeBankOperator = await state.bridgeBank.operator();
      bridgeBankOperator.should.be.equal(operator.address);
    });

    it("should allow user to lock ERC20 tokens", async function () {
      // Add the token into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

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
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token2.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token3.address, true)
        .should.be.fulfilled;

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
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token2.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token3.address, true)
        .should.be.fulfilled;

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

    it("should not allow user to multi-burn ERC20 tokens that are not cosmos native assets", async function () {
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
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender, state.sender, state.sender],
          [state.token1.address,state.token2.address,state.token3.address],
          [state.amount, state.amount, state.amount],
          [true, false, false]
        ),
      ).to.be.revertedWith("Only token in whitelist can be burned");
    });

    it("should allow user to multi-lock and burn ERC20 tokens and rowan with multiLockBurn method", async function () {
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token2.address, true)
        .should.be.fulfilled;

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

      // Lock & burn tokens
      const tx = await state.bridgeBank.connect(userOne).multiLockBurn(
        [state.sender, state.sender, state.sender],
        [state.token1.address, state.token2.address, state.rowan.address],
        [state.amount, state.amount, state.amount],
        [false, false, true]
      );

      await tx.wait();

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
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token2.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token3.address, true)
        .should.be.fulfilled;

      const tx = await state.token1.connect(userOne).approve(
        state.bridgeBank.address,
        0
      );
      const receipt = await tx.wait();

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
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token2.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token3.address, true)
        .should.be.fulfilled;

      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLock(
          [state.sender + "ee", state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount]
        )
      ).to.be.revertedWith("INV_SIF_ADDR");
    });

    it("should not allow user to multi-lock when bridgebank is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLock(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount]
        )
      ).to.be.revertedWith("Pausable: paused");
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

    it("should revert when multi-lock parameters are malformed, not enough booleans", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender, state.sender],
          [state.token1.address, state.token2.address],
          [state.amount, state.amount],
          [false, false, false]
        )
      ).to.be.revertedWith("M_P");
    });

    it("should revert when multi-lock parameters are malformed, invalid sif addresses", async function () {
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token2.address, true)
        .should.be.fulfilled;

      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token3.address, true)
        .should.be.fulfilled;

      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender + "ee", state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount],
          [false, false, false]
        )
      ).to.be.revertedWith("INV_SIF_ADDR");
    });

    it("should not allow user to multi-lock/burn when bridgebank is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).multiLockBurn(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount],
          [false, false, false]
        )
      ).to.be.revertedWith("Pausable: paused");
    });
  });

  describe("Whitelist", function () {
    it("should allow the operator to add a token to the whitelist", async function () {
      // Add the tokens into white list
      await expect(state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true))
        .to.emit(state.bridgeBank, 'LogWhiteListUpdate')
        .withArgs(state.token1.address, true);

      expect(await state.bridgeBank.getTokenInEthWhiteList(state.token1.address)).to.be.equal(true);
    });

    it("should not allow user to add a token to the whitelist", async function () {
      // Add the tokens into white list
      await expect(state.bridgeBank.connect(userOne)
        .updateEthWhiteList(state.token1.address, true))
        .to.be.revertedWith("!operator");

      expect(await state.bridgeBank.getTokenInEthWhiteList(state.token1.address)).to.be.equal(false);
    });

    it("should allow the operator to remove a token from the whitelist", async function () {
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      // Remove the token from whitelist
      await expect(state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, false))
        .to.emit(state.bridgeBank, 'LogWhiteListUpdate')
        .withArgs(state.token1.address, false);

      expect(await state.bridgeBank.getTokenInEthWhiteList(state.token1.address)).to.be.equal(false);
    });

    it("should not allow user to remove a token from the whitelist", async function () {
      // Add the tokens into white list
      await state.bridgeBank.connect(operator)
        .updateEthWhiteList(state.token1.address, true)
        .should.be.fulfilled;

      await expect(state.bridgeBank.connect(userOne)
        .updateEthWhiteList(state.token1.address, false))
        .to.be.revertedWith("!operator");

      expect(await state.bridgeBank.getTokenInEthWhiteList(state.token1.address)).to.be.equal(true);
    });

    it("should not allow user to lock ERC20 tokens that are not in whitelist", async function () {
      await state.token1.connect(userOne).approve(
        state.bridgeBank.address,
        state.amount
      );

      // Attempt to lock tokens
      await expect(state.bridgeBank.connect(userOne).lock(
        state.sender,
        state.token1.address,
        state.amount
      )).to.be.revertedWith("Only token in whitelist can be transferred to cosmos");
    });

    it("should not allow user to multi-lock ERC20 tokens that are not in whitelist", async function () {
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
      await expect(state.bridgeBank.connect(userOne).multiLock(
        [state.sender, state.sender, state.sender],
        [state.token1.address,state.token2.address,state.token3.address],
        [state.amount, state.amount, state.amount]
      )).to.be.revertedWith("Only token in whitelist can be transferred to cosmos");
    });

    it("should not allow user to multi-lock ERC20 tokens that are not in whitelist with multiLockBurn method", async function () {
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
      await expect(state.bridgeBank.connect(userOne).multiLockBurn(
        [state.sender, state.sender, state.sender],
        [state.token1.address,state.token2.address,state.token3.address],
        [state.amount, state.amount, state.amount],
        [false, false, false]
      )).to.be.revertedWith("Only token in whitelist can be transferred to cosmos");
    });

    it("should not allow user to multi-lock and burn ERC20 tokens not in whitelist and rowan with multiLockBurn method", async function () {
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

      // Lock & burn tokens
      await expect(state.bridgeBank.connect(userOne).multiLockBurn(
        [state.sender, state.sender, state.sender],
        [state.token1.address, state.token2.address, state.rowan.address],
        [state.amount, state.amount, state.amount],
        [false, false, true]
      )).to.be.revertedWith("Only token in whitelist can be transferred to cosmos");
    });

    it("should not allow the operator to add a token to the whitelist if it's in cosmosWhitelist", async function () {
      // assert that the cosmos bridge token has not been created
      let bridgeToken = await state.cosmosBridge.sourceAddressToDestinationAddress(
        state.token1.address
      );
      expect(bridgeToken).to.be.equal(state.ethereumToken);

      await state.cosmosBridge.connect(userOne).createNewBridgeToken(
        state.symbol,
        state.name,
        state.token1.address,
        18,
        1,
      );

      // now assert that the bridge token has been created
      bridgeToken = await state.cosmosBridge.sourceAddressToDestinationAddress(
        state.token1.address
      );
      expect(bridgeToken).to.not.be.equal(state.ethereumToken);

      // assert the bridgeToken is in cosmosWhitelist
      const isInCosmosWhitelist = await state.bridgeBank.getCosmosTokenInWhiteList(
        bridgeToken
      );
      expect(isInCosmosWhitelist).to.be.equal(true);

      // Try adding the token into white list
      await expect(state.bridgeBank.connect(operator)
        .updateEthWhiteList(bridgeToken, true))
        .to.be.revertedWith('whitelisted');

      expect(await state.bridgeBank.getTokenInEthWhiteList(bridgeToken)).to.be.equal(false);
    });
  });
});
