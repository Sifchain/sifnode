import web3 from "web3";
import { ethers } from "hardhat";
import { use, expect } from "chai";
import { solidity } from "ethereum-waffle";
import { setup, getValidClaim, TestFixtureState } from "./helpers/testFixture";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { FailHardToken, FailHardToken__factory } from "../build";

const BigNumber = ethers.BigNumber;
// require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

describe("Fail Hard Token", function () {
  let userOne: SignerWithAddress;
  let userTwo: SignerWithAddress;
  let addresses: string[];
  let state: TestFixtureState;
  let failFactory: FailHardToken__factory;
  let failHardToken: FailHardToken;
  let accounts: SignerWithAddress[];

  before(async function () {
    accounts = await ethers.getSigners();
    addresses = accounts.map((e) => {
      return e.address;
    });

    failFactory = await ethers.getContractFactory("FailHardToken");

    userOne = accounts[6];
    userTwo = accounts[7];
  });

  beforeEach(async function () {
    state = await setup(
      addresses.slice(2, 6),
      [25, 25, 25, 25],
      accounts[0],
      75,
      accounts[1],
      userOne,
      userTwo,
      accounts[8],
      1,
    );
    state.amount = 1000;

    failHardToken = await failFactory.deploy(
      "Fail Hard Token",
      "FAIL",
      userOne.address,
      state.amount
    );

    await failHardToken.deployed();
  });

  describe("Burn, Unlock", function () {
    it("should successfully process an unlock claim when token reverts transfer(), but user does not receive any tokens back", async function () {
      // Get balances before locking
      const beforeBridgeBankBalance = Number(
        await failHardToken.balanceOf(state.bridgeBank.address)
      );
      expect(beforeBridgeBankBalance).to.equal(0);

      const beforeUserBalance = Number(await failHardToken.balanceOf(userOne.address));
      expect(beforeUserBalance).to.equal(state.amount);

      // Attempt to lock tokens (will work without a previous approval - only for FailHardToken)
      await state.bridgeBank
        .connect(userOne)
        .lock(state.sender, failHardToken.address, state.amount);

      // Confirm that the tokens have left the user's wallet
      const afterUserBalance = Number(await failHardToken.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(0);

      // Confirm that bridgeBank now owns the tokens:
      const afterBridgeBankBalance = Number(
        await failHardToken.balanceOf(state.bridgeBank.address)
      );
      expect(afterBridgeBankBalance).to.equal(state.amount);

      // Now, try to unlock those tokens...

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        userOne.address,
        failHardToken.address,
        state.amount,
        "Fail Hard Token",
        "FAIL",
        18,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.none,
        accounts.slice(2, 6),
      );

      await state.cosmosBridge
        .connect(accounts[2])
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      // The prophecy should be completed without reverting, but the user shouldn't receive anything
      const balance = Number(await failHardToken.balanceOf(userOne.address));
      expect(balance).to.be.equal(0);
    });
  });

  it("should revert on burn if user is blocklisted", async function () {
    await state.bridgeBank.connect(state.owner).addExistingBridgeToken(failHardToken.address);

    const beforeUserBalance = Number(await failHardToken.balanceOf(userOne.address));

    // Lock & burn tokens
    const tx = await expect(
      state.bridgeBank
        .connect(userOne)
        .multiLockBurn([state.sender], [failHardToken.address], [state.amount], [true])
    ).to.be.reverted;

    // assert that there was no change
    const afterUserBalance = Number(await failHardToken.balanceOf(userOne.address));
    expect(afterUserBalance).to.equal(beforeUserBalance);
  });
});
