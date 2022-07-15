import { setup, TestFixtureState } from "./helpers/testFixture";
import { upgrades } from "hardhat";
import { use, expect } from "chai";

import web3 from "web3";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { CosmosBridge, CosmosBridge__factory, MockCosmosBridgeUpgrade } from "../build";

const BigNumber = ethers.BigNumber;

describe("CosmosBridge Upgrade", function () {
  const consensusThreshold = 70;
  let userOne: SignerWithAddress;
  let userTwo: SignerWithAddress;
  let userThree: SignerWithAddress;
  let userFour: SignerWithAddress;
  let accounts: SignerWithAddress[];
  let signerAccounts: string[];
  let operator: SignerWithAddress;
  let owner: SignerWithAddress;
  let initialPowers: number[];
  let initialValidators: string[];
  let networkDescriptor: number;
  let state: TestFixtureState;
  let pauser: SignerWithAddress;
  let MockCosmosBridgeUpgrade: MockCosmosBridgeUpgrade;

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
    MockCosmosBridgeUpgrade = await ethers.getContractFactory("MockCosmosBridgeUpgrade");
  });

  describe("CosmosBridge smart contract deployment", function () {
    beforeEach(async function () {
      state = await setup(
        [userOne.address, userTwo.address, userThree.address, userFour.address],
        [30, 20, 21, 29],
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser,
        networkDescriptor,
      );

      state.cosmosBridge = (await upgrades.upgradeProxy(
        state.cosmosBridge.address,
        MockCosmosBridgeUpgrade as unknown as CosmosBridge__factory
      ) as CosmosBridge);
    });

    it("should be able to mint tokens for a user", async function () {
      const amount = 100000000000;
      expect(state.cosmosBridge).to.exist;

      await (state.cosmosBridge.connect(operator) as MockCosmosBridgeUpgrade).tokenFaucet();
      const operatorBalance = await (state.cosmosBridge as MockCosmosBridgeUpgrade).balanceOf(operator.address);
      expect(Number(operatorBalance)).to.equal(amount);
    });

    it("should be able to transfer tokens from the operator", async function () {
      const startingOperatorBalance = await (state.cosmosBridge as MockCosmosBridgeUpgrade).balanceOf(operator.address);
      expect(Number(startingOperatorBalance)).to.equal(0);

      const amount = 100000000000;
      expect(state.cosmosBridge).to.exist;

      await (state.cosmosBridge.connect(operator) as MockCosmosBridgeUpgrade).tokenFaucet();
      await (state.cosmosBridge.connect(operator) as MockCosmosBridgeUpgrade).transfer(userOne.address, amount);

      const operatorBalance = await (state.cosmosBridge as MockCosmosBridgeUpgrade).balanceOf(operator.address);
      const userOneBalance = await (state.cosmosBridge as MockCosmosBridgeUpgrade).balanceOf(userOne.address);

      expect(Number(operatorBalance)).to.equal(0);
      expect(Number(userOneBalance)).to.equal(amount);
    });

    it("should not be able to initialize cosmos bridge a second time", async function () {
      expect(state.cosmosBridge).to.exist;

      await expect(
        state.cosmosBridge.initialize(
          userFour.address,
          50,
          state.initialValidators,
          state.initialPowers,
          state.networkDescriptor
        )
      ).to.be.revertedWith("Initialized");
    });

    describe("Storage Remains Intact", function () {
      it("should not allow the operator to update the Bridge Bank once it has been set", async function () {
        await expect(
          state.cosmosBridge.connect(operator).setBridgeBank(state.bridgeBank.address)
        ).to.be.revertedWith("The Bridge Bank cannot be updated once it has been set");
      });
    });
  });
});
