import web3 from "web3";
import { ethers } from "hardhat";
import { use, expect } from "chai";
import { solidity } from "ethereum-waffle";
import { setup, TestFixtureState } from "./helpers/testFixture";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BridgeToken__factory, ManyDecimalsToken, ManyDecimalsToken__factory } from "../build";

const BigNumber = ethers.BigNumber;

use(solidity);

describe("Many Decimals Token", function () {
  let userOne: SignerWithAddress;
  let userTwo: SignerWithAddress;
  let addresses: string[];
  let state: TestFixtureState;
  let tokenFactory: ManyDecimalsToken__factory;
  let manyDecimalsToken: ManyDecimalsToken;
  let accounts: SignerWithAddress[];

  before(async function () {
    accounts = await ethers.getSigners();
    addresses = accounts.map((e) => {
      return e.address;
    });

    tokenFactory = await ethers.getContractFactory("ManyDecimalsToken");

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

    manyDecimalsToken = await tokenFactory.deploy(
      "Many Decimals Token",
      "DEC",
      userOne.address,
      state.amount
    );

    await manyDecimalsToken.deployed();
  });

  describe("Lock", function () {
    it("should allow user to lock a token that has 100 decimals, but not twice!", async function () {
      // approve bridgebank to spend rowan
      await manyDecimalsToken.connect(userOne).approve(state.bridgeBank.address, state.amount);

      // Lock & burn tokens
      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn([state.sender], [manyDecimalsToken.address], [state.amount / 2], [false])
      )
        .to.emit(state.bridgeBank, "LogLock")
        .withArgs(
          userOne.address,
          state.sender,
          manyDecimalsToken.address,
          state.amount / 2,
          1,
          255,
          "DEC",
          "Many Decimals Token",
          1
        );

      const afterUserBalance = Number(await manyDecimalsToken.balanceOf(userOne.address));
      expect(afterUserBalance).to.equal(state.amount / 2);

      await expect(
        state.bridgeBank
          .connect(userOne)
          .multiLockBurn([state.sender], [manyDecimalsToken.address], [state.amount / 2], [false])
      )
        .to.emit(state.bridgeBank, "LogLock")
        .withArgs(
          userOne.address,
          state.sender,
          manyDecimalsToken.address,
          state.amount / 2,
          2,
          255,
          "DEC",
          "Many Decimals Token",
          1
        );
    });
  });
});
