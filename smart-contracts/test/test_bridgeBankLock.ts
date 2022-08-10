const Web3Utils = require("web3-utils");

import { ethers, network } from "hardhat";
import { use, expect } from "chai";
import { solidity } from "ethereum-waffle";
import { setup, deployCommissionToken, TestFixtureState } from "./helpers/testFixture";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { ContractTransaction } from "ethers";
import { BridgeBank } from "../build";

use(solidity);

const getBalance = async function (address: string) {
  return await network.provider.send("eth_getBalance", [address]);
};

const BigNumber = ethers.BigNumber;

describe("Test Bridge Bank", function () {
  let userOne: SignerWithAddress;
  let userTwo: SignerWithAddress;
  let userThree: SignerWithAddress;
  let userFour: SignerWithAddress;
  let accounts: SignerWithAddress[];
  let signerAccounts: string[];
  let operator: SignerWithAddress;
  let owner: SignerWithAddress;
  let pauser: SignerWithAddress;
  const consensusThreshold = 75;
  let initialPowers: number[];
  let initialValidators: string[];
  let networkDescriptor: number;
  // track the state of the deployed contracts
  let state: TestFixtureState;

  before(async function () {
    accounts = await ethers.getSigners();

    signerAccounts = accounts.map((e) => {
      return e.address;
    });

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
    state = await setup(
      initialValidators,
      initialPowers,
      operator,
      consensusThreshold,
      owner,
      userOne,
      userThree,
      pauser,
      networkDescriptor,
    );

    const tokens = [state.token, state.token1, state.token2, state.token3, state.token_ibc, state.token_noDenom];
    const users = [userOne, userTwo, userThree, userFour];
    const mintPromises: Promise<ContractTransaction>[] = [];
    const approvePromises: Promise<ContractTransaction>[] = [];

    for (const user of users) {
      for (const token of tokens) {
        mintPromises.push(token.connect(operator).mint(user.address, state.amount * 2));
        approvePromises.push(token.connect(user).approve(state.bridgeBank.address, state.amount * 2));
      }
    }
    await Promise.all(mintPromises);
    await Promise.all(approvePromises);
  });

  describe("BridgeBank", function () {
    it("should deploy the BridgeBank, correctly setting the operator", async function () {
      expect(state.bridgeBank).to.exist;

      const bridgeBankOperator = await state.bridgeBank.operator();
      expect(bridgeBankOperator).to.equal(operator.address);
    });

    it("should allow user to lock ERC20 tokens", async function () {
      // Get balances before locking
      const beforeBridgeBankBalance = Number(
        await state.token1.balanceOf(state.bridgeBank.address)
      );
      expect(beforeBridgeBankBalance).to.equal(0);

      const beforeUserBalance = Number(await state.token1.balanceOf(userOne.address));
       expect(beforeUserBalance).to.equal(state.amount * 2);

      // Attempt to lock tokens
      await state.bridgeBank
        .connect(userOne)
        .lock(state.sender, state.token1.address, state.amount);

      // Confirm that the tokens have left the user's wallet
      const afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      // Confirm that bridgeBank now owns the tokens:
      const afterBridgeBankBalance = Number(await state.token1.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(state.amount);
    });

    it("should allow user to lock Commission Charging ERC20 tokens", async function () { 
      const devAccount = userOne.address;
      const user = userTwo;
      const devFee = 500; // 5% dev fee
      const initialBalance = 100_000
      const token = await deployCommissionToken(devAccount, devFee, user.address, initialBalance);
      const transferAmount = 10_000
      const remaining = 90_000;
      const afterTransfer = 9_500;
      // Get balances before locking
      const beforeBridgeBankBalance = ethers.BigNumber.from(
        await token.balanceOf(state.bridgeBank.address)
      );
      expect(beforeBridgeBankBalance).to.equal(0);

      const beforeUserBalance = ethers.BigNumber.from(await token.balanceOf(user.address));
      expect(beforeUserBalance).to.equal(initialBalance);
      // Approve transfer of tokens first
      await token.connect(user).approve(state.bridgeBank.address, transferAmount);
      // Attempt to lock tokens
      await state.bridgeBank
        .connect(user)
        .lock(state.sender, token.address, transferAmount);

      // Confirm that the tokens have left the user's wallet
      const afterUserBalance = ethers.BigNumber.from(await token.balanceOf(user.address));
      expect(afterUserBalance).to.equal(remaining);

      // Confirm that bridgeBank now owns the tokens:
      const afterBridgeBankBalance = ethers.BigNumber.from(await token.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(afterTransfer);
    });

    it("should allow users to lock Ethereum in the bridge bank", async function () {
      const tx = await state.bridgeBank
        .connect(userOne)
        .lock(state.sender, state.constants.zeroAddress, state.weiAmount, {
          value: state.weiAmount,
        });
      await tx.wait();

      const contractBalanceWei = await getBalance(state.bridgeBank.address);
      const contractBalance = Web3Utils.fromWei(contractBalanceWei, "ether");

      expect(contractBalance).to.equal(
        Web3Utils.fromWei((+state.weiAmount).toString(), "ether")
      );
    });

    it("should NOT allow a blocklisted user to lock ERC20 tokens", async function () {
      // Add userOne to the blocklist:
      await expect(state.blocklist.connect(operator).addToBlocklist(userOne.address)).to.not.be.reverted;

      // Get balances before locking
      const beforeBridgeBankBalance = Number(
        await state.token1.balanceOf(state.bridgeBank.address)
      );
      expect(beforeBridgeBankBalance).to.equal(0);

      const beforeUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(beforeUserBalance).to.equal(state.amount * 2);

      // Attempt to lock tokens and fail
      await expect(
        state.bridgeBank.connect(userOne).lock(state.sender, state.token1.address, state.amount)
      ).to.be.revertedWith("Address is blocklisted");

      // Confirm that the tokens have NOT left the user's wallet
      const afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(beforeUserBalance);

      // Confirm that bridgeBank did not receive the tokens:
      const afterBridgeBankBalance = Number(await state.token1.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(beforeBridgeBankBalance);
    });

    it("should NOT allow a blocklisted user to lock Ethereum in the bridge bank", async function () {
      // Add userOne to the blocklist:
      await expect(state.blocklist.connect(operator).addToBlocklist(userOne.address)).to.not.be.reverted;

      await expect(
        state.bridgeBank
          .connect(userOne)
          .lock(state.sender, state.constants.zeroAddress, state.weiAmount, {
            value: state.weiAmount,
          })
      ).to.be.revertedWith("Address is blocklisted");

      const contractBalanceWei = await getBalance(state.bridgeBank.address);
      const expectedValue = BigNumber.from(0);

      expect(contractBalanceWei).to.equal(expectedValue);
    });
  });

  describe("Multi Lock ERC20 Tokens", function () {
    it("should allow user to multi-lock ERC20 tokens", async function () {
      const previousNonce = await state.bridgeBank.connect(userOne).lockBurnNonce();

      // Attempt to lock tokens
      await state.bridgeBank
        .connect(userOne)
        .multiLockBurn(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount],
          [false, false, false]
        );

      expect(await state.bridgeBank.connect(userOne).lockBurnNonce()).to.equal(previousNonce.add(3));

      // Confirm that the user has been minted the correct token
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.token3.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);
    });

     it("should allow user to multi-lock ERC20 tokens including commission tokens", async function () {
      const devAccount = userTwo.address;
      const user = userOne;
      const devFee = 500; // 5% dev fee
      const initialBalance = 100_000
      const token = await deployCommissionToken(devAccount, devFee, user.address, initialBalance);
      const transferAmount = 10_000
      const remaining = 90_000;
      const afterTransfer = 9_500;

      const beforeBridgeBankBalance = Number(
        await token.balanceOf(state.bridgeBank.address)
      );
      expect(beforeBridgeBankBalance).to.equal(0);

      // Approve bridgebank as a spender
      await token.connect(user).approve(state.bridgeBank.address, transferAmount);

      // Attempt to lock tokens
      await state.bridgeBank
        .connect(userOne)
        .multiLockBurn(
          [state.sender, state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address, token.address],
          [state.amount, state.amount, state.amount, transferAmount],
          [false, false, false, false]
        );

      // Confirm that the user has been minted the correct token
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.token3.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      // Confirm that the commission tokens have left the user's wallet
      afterUserBalance = Number(await token.balanceOf(user.address));
      expect(afterUserBalance).to.equal(remaining);

      // Confirm that bridgeBank now owns the commission tokens:
      const afterBridgeBankBalance = Number(await token.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(afterTransfer);
    });

    it("should NOT allow a blocklisted user to multi-lock ERC20 tokens", async function () {
      // Add userOne to the blocklist:
      await expect(state.blocklist.connect(operator).addToBlocklist(userOne.address)).to.not.be.reverted;

      // Attempt to lock tokens and fail
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("Address is blocklisted");

      // Confirm that the tokens have not left the user's wallet
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.token3.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      let afterBridgeBankBalance = Number(await state.token1.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);

      afterBridgeBankBalance = Number(await state.token2.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);

      afterBridgeBankBalance = Number(await state.token3.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);
    });

    it("should allow user to multi-lock ERC20 tokens with multiLockBurn method", async function () {
      // Attempt to lock tokens
      await state.bridgeBank
        .connect(userOne)
        .multiLockBurn(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.token3.address],
          [state.amount, state.amount, state.amount],
          [false, false, false]
        );

      // Confirm that the user has been minted the correct token
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.token3.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);
    });

    it("should NOT allow a blocklisted user to multi-lock ERC20 tokens with multiLockBurn method", async function () {
      // Add userOne to the blocklist:
      await expect(state.blocklist.connect(operator).addToBlocklist(userOne.address)).to.not.be.reverted;

      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("Address is blocklisted");

      // Confirm that the user has been minted the correct token
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.token3.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      let afterBridgeBankBalance = Number(await state.token1.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);

      afterBridgeBankBalance = Number(await state.token2.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);

      afterBridgeBankBalance = Number(await state.token3.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);
    });

    it("should NOT allow user to multi-burn ERC20 tokens that are not cosmos native assets", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [true, false, false]
          )
      ).to.be.revertedWith("Token is not in Cosmos whitelist");
    });

    it("should allow user to multi-lock and burn ERC20 tokens and rowan with multiLockBurn method", async function () {
      // approve bridgebank to spend rowan
      await state.rowan.connect(userOne).approve(state.bridgeBank.address, state.amount);

      // Lock & burn tokens
      const tx = await state.bridgeBank
        .connect(userOne)
        .multiLockBurn(
          [state.sender, state.sender, state.sender],
          [state.token1.address, state.token2.address, state.rowan.address],
          [state.amount, state.amount, state.amount],
          [false, false, true]
        );

      await tx.wait();

      // Confirm that the user has the proper balance after the multiLockBurn
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);

      afterUserBalance = Number(await state.rowan.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount);
    });

    it("should NOT allow a blocklisted user to multi-lock and burn ERC20 tokens and rowan with multiLockBurn method", async function () {
      // Add userOne to the blocklist:
      await expect(state.blocklist.connect(operator).addToBlocklist(userOne.address)).to.not.be.reverted;

      // approve bridgebank to spend rowan
      await state.rowan.connect(userOne).approve(state.bridgeBank.address, state.amount);

      // Lock & burn tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.rowan.address],
            [state.amount, state.amount, state.amount],
            [false, false, true]
          )
      ).to.be.revertedWith("Address is blocklisted");

      // Confirm that the user has the proper balance after the multiLockBurn
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.rowan.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      let afterBridgeBankBalance = Number(await state.token1.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);

      afterBridgeBankBalance = Number(await state.token2.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);

      afterBridgeBankBalance = Number(await state.rowan.balanceOf(state.bridgeBank.address));
      expect(afterBridgeBankBalance).to.equal(0);
    });

    it("should NOT allow user to multi-lock ERC20 tokens if one token is not fully approved", async function () {
      const tx = await state.token1.connect(userOne).approve(state.bridgeBank.address, 0);
      const receipt = await tx.wait();

      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("transfer amount exceeds allowance");

      // Confirm that user token balances have stayed the same
      let afterUserBalance = Number(await state.token1.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.token2.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);

      afterUserBalance = Number(await state.token3.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount * 2);
    });

    it("should NOT allow user to multi-lock when parameters are malformed, not enough token amounts", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("M_P");
    });

    it("should NOT allow user to multi-lock when parameters are malformed, not enough token addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("M_P");
    });

    it("should NOT allow user to multi-lock when parameters are malformed, not enough sif addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("M_P");
    });

    it("should NOT allow user to multi-lock when parameters are malformed, invalid sif addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender + "ee", state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("INV_SIF_ADDR");
    });

    it("should NOT allow user to multi-lock when bridgebank is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("Pausable: paused");
    });

    it("should NOT allow user to multi-lock ERC20 tokens and Eth in the same call", async function () {
      // Attempt to lock tokens and Ether in the same call
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender, state.sender],
            [
              state.token1.address,
              state.token2.address,
              state.token3.address,
              state.constants.zeroAddress,
            ],
            [state.amount, state.amount, state.amount, state.amount],
            [false, false, false, false],
            { value: 100 } as any // Typescript and typechain are smart enough to know this is a non-payable function so we must override that for this check
          ),
      ).to.be.reverted; // Non-Payable function may prevent ethersjs from getting to the revert
    });

    it("should NOT allow user to multi-burn tokens and Eth in the same call", async function () {
      // Add the tokens into whitelist
      // Also, add Ether into whitelist, which shouldn't be done but
      // we'll indulge in this scenario to bypass the whitelist requirements
      await state.bridgeBank
        .connect(owner)
        .batchAddExistingBridgeTokens([
          state.token1.address,
          state.token2.address,
          state.token3.address,
          state.constants.zeroAddress,
        ]);

      // Attempt to burn tokens and Ether in the same call
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender, state.sender],
            [
              state.token1.address,
              state.token2.address,
              state.token3.address,
              state.constants.zeroAddress,
            ],
            [state.amount, state.amount, state.amount, state.amount],
            [true, true, true, true]
          )
      ).to.be.reverted;
    });
  });

  describe("Multi Lock Burn ERC20 Tokens", function () {
    it("should revert when parameters are malformed, not enough token amounts", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("M_P");
    });

    it("should revert when multi-lock parameters are malformed, not enough token addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
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
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
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
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender],
            [state.token1.address, state.token2.address],
            [state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("M_P");
    });

    it("should revert when multi-lock parameters are malformed, invalid sif addresses", async function () {
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender + "ee", state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("INV_SIF_ADDR");
    });

    it("should NOT allow user to multi-lock/burn when bridgebank is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("Pausable: paused");
    });
  });

  describe("Whitelist", function () {
    it("should NOT allow user to lock ERC20 tokens that are in Cosmos whitelist", async function () {
      // add token as BridgeToken
      await state.bridgeBank.connect(owner).addExistingBridgeToken(state.token1.address);

      // Attempt to lock tokens
      await expect(
        state.bridgeBank.connect(userOne).lock(state.sender, state.token1.address, state.amount)
      ).to.be.revertedWith("Only token not in cosmos whitelist can be locked");
    });

    it("should NOT allow user to multi-lock ERC20 tokens if at least one of them is in cosmos whitelist", async function () {
      // add token1 as BridgeToken
      await state.bridgeBank.connect(owner).addExistingBridgeToken(state.token1.address);

      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("Only token not in cosmos whitelist can be locked");
    });

    it("should NOT allow user to multi-lock ERC20 tokens with multiLockBurn method if one of them is cosmos whitelist", async function () {
      // add token1 as BridgeToken
      await state.bridgeBank.connect(owner).addExistingBridgeToken(state.token1.address);

      // Attempt to lock tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.token3.address],
            [state.amount, state.amount, state.amount],
            [false, false, false]
          )
      ).to.be.revertedWith("Only token not in cosmos whitelist can be locked");
    });

    it("should NOT allow user to multi-lock and burn ERC20 tokens and rowan with multiLockBurn method if at least one of them is in cosmos whitelist ", async function () {
      // add token1 as BridgeToken
      await state.bridgeBank.connect(owner).addExistingBridgeToken(state.token1.address);

      // approve bridgebank to spend rowan
      await state.rowan.connect(userOne).approve(state.bridgeBank.address, state.amount);

      // Lock & burn tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn(
            [state.sender, state.sender, state.sender],
            [state.token1.address, state.token2.address, state.rowan.address],
            [state.amount, state.amount, state.amount],
            [false, false, true]
          )
      ).to.be.revertedWith("Only token not in cosmos whitelist can be locked");
    });
  });
});
