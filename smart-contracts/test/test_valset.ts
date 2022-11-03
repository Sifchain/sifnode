import { ethers } from "hardhat";
import { use, expect } from "chai";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { TestFixtureState, setup } from "./helpers/testFixture";

const EVMRevert = "revert";
const BigNumber = ethers.BigNumber;

interface TestValsetState extends TestFixtureState {
  userTwoPower: number;
  userThreePower: number;
  secondValidators: string[];
  secondPowers: number[];
}

describe("Test Valset", function () {
  let userOne: SignerWithAddress;
  let userTwo: SignerWithAddress;
  let userThree: SignerWithAddress;
  let userFour: SignerWithAddress;
  let accounts: SignerWithAddress[];
  let signerAccounts: string[];
  let operator: SignerWithAddress;
  let owner: SignerWithAddress;
  let pauser: SignerWithAddress;
  const consensusThreshold = 80;
  let initialPowers: number[];
  let initialValidators: string[];
  let networkDescriptor: number;
  // track the state of the deployed contracts
  let state: TestValsetState;

  describe("Valset contract deployment", function () {
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
      state = (await setup(
        [userOne.address, userTwo.address, userThree.address],
        [5, 8, 12],
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser,
        networkDescriptor,
      ) as TestValsetState);
    });

    it("should deploy the Valset and correctly set the current valset version", async function () {
      expect(state.cosmosBridge).to.exist;

      const valsetValsetVersion = await state.cosmosBridge.currentValsetVersion();
      expect(Number(valsetValsetVersion)).to.equal(1);
    });

    it("should correctly set initial validators and initial validator count", async function () {
      const userOneValidator = await state.cosmosBridge.isActiveValidator(userOne.address);
      const userTwoValidator = await state.cosmosBridge.isActiveValidator(userTwo.address);
      const userThreeValidator = await state.cosmosBridge.isActiveValidator(userThree.address);
      const valsetValidatorCount = await state.cosmosBridge.validatorCount();

      expect(userOneValidator).to.equal(true);
      expect(userTwoValidator).to.equal(true);
      expect(userThreeValidator).to.equal(true);
      expect(Number(valsetValidatorCount)).to.equal(state.initialValidators.length);
    });

    it("should correctly set initial validator powers ", async function () {
      const userOnePower = await state.cosmosBridge.getValidatorPower(userOne.address);
      const userTwoPower = await state.cosmosBridge.getValidatorPower(userTwo.address);
      const userThreePower = await state.cosmosBridge.getValidatorPower(userThree.address);

      expect(Number(userOnePower)).to.equal(state.initialPowers[0]);
      expect(Number(userTwoPower)).to.equal(state.initialPowers[1]);
      expect(Number(userThreePower)).to.equal(state.initialPowers[2]);
    });

    it("should correctly set the initial total power", async function () {
      const valsetTotalPower = await state.cosmosBridge.totalPower();

      expect(Number(valsetTotalPower)).to.equal(
        state.initialPowers[0] + state.initialPowers[1] + state.initialPowers[2]
      );
    });
  });

  describe("Dynamic validator set", function () {
    describe("Adding validators", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address];
        state.initialPowers = [5];

        state = (await setup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor,
        ) as TestValsetState);

        state.userTwoPower = 11;
        state.userThreePower = 44;
      });

      it("should correctly update the valset when the operator adds a new validator", async function () {
        // Confirm initial validator count
        const priorValsetValidatorCount = await state.cosmosBridge.validatorCount();
        expect(Number(priorValsetValidatorCount)).to.equal(1);

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();

        expect(Number(priorTotalPower)).to.equal(state.initialPowers[0]);

        // Operator adds a validator
        await expect(state.cosmosBridge.connect(operator).addValidator(userTwo.address, state.userTwoPower))
          .not.to.be.reverted;

        // Confirm that userTwo has been set as a validator
        const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(userTwo.address);
        expect(isUserTwoValidator).to.equal(true);

        // Confirm that userTwo's power has been correctly set
        const userTwoSetPower = await state.cosmosBridge.getValidatorPower(userTwo.address);
        expect(Number(userTwoSetPower)).to.equal(state.userTwoPower);

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        expect(Number(postValsetValidatorCount)).to.equal(2);

        // Confirm updated total power
        const postTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(postTotalPower)).to.equal(
          state.initialPowers[0] + state.userTwoPower
        );
      });

      it("should be able to add a new validator and get its power", async function () {
        // Get the event logs from the addition of a new validator
        await state.cosmosBridge
          .connect(operator)
          .addValidator(userTwo.address, state.userTwoPower);

        const userTwoSetPower = await state.cosmosBridge.getValidatorPower(userTwo.address);
        expect(Number(userTwoSetPower)).to.equal(state.userTwoPower);
      });

      it("should allow the operator to add multiple new validators", async function () {
        // Fail if not operator
        await expect(
          state.cosmosBridge.connect(userOne).addValidator(userTwo.address, state.userTwoPower)
        ).to.be.revertedWith("Must be the operator.");

        await expect(state.cosmosBridge.connect(operator).addValidator(userTwo.address, state.userTwoPower))
          .not.to.be.reverted;
        await expect(state.cosmosBridge
          .connect(operator)
          .addValidator(userThree.address, state.userThreePower)).to.not.be.reverted;
        await expect(state.cosmosBridge.connect(operator).addValidator(accounts[4].address, 77)).to.not.be.reverted;
        await expect(state.cosmosBridge.connect(operator).addValidator(accounts[5].address, 23)).to.not.be.reverted;

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        expect(Number(postValsetValidatorCount)).to.equal(5);

        // Confirm updated total power
        const valsetTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(valsetTotalPower)).to.equal(
          state.initialPowers[0] + state.userTwoPower + state.userThreePower + 100 // (23 + 77)
        );
      });

      it("should not let you add the same validator twice", async function () {
        await expect(state.cosmosBridge.
            connect(operator)
            .addValidator(userOne.address, state.userThreePower))
            .to.be.revertedWith("Already a validator");
      })
    });

    describe("Updating validator's power", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address];
        state.initialPowers = [5];

        // Deploy CosmosBridge contract
        state = (await setup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor,
        ) as TestValsetState);

        state.userTwoPower = 11;
        state.userThreePower = 44;
      });

      it("should allow the operator to update a validator's power", async function () {
        const NEW_POWER = 515;

        // Confirm userOne's initial power
        const userOneInitialPower = await state.cosmosBridge.getValidatorPower(userOne.address);

        expect(Number(userOneInitialPower)).to.equal(state.initialPowers[0]);

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(priorTotalPower)).to.equal(state.initialPowers[0]);

        // Fail if not operator
        await expect(
          state.cosmosBridge.connect(userTwo).updateValidatorPower(userOne.address, NEW_POWER)
        ).to.be.revertedWith("Must be the operator.");

        // Operator updates the validator's initial power
        await expect(state.cosmosBridge.connect(operator).updateValidatorPower(userOne.address, NEW_POWER))
          .to.not.be.reverted;

        // Confirm userOne's power has increased
        const userOnePostPower = await state.cosmosBridge.getValidatorPower(userOne.address);
        expect(Number(userOnePostPower)).to.equal(NEW_POWER);

        // Confirm total power has been updated
        const postTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(postTotalPower)).to.equal(NEW_POWER);
      });

      it("should update of a validator's power", async function () {
        const NEW_POWER = 111;

        await state.cosmosBridge.connect(operator).updateValidatorPower(userOne.address, NEW_POWER);

        const userTwoPower = await state.cosmosBridge.getValidatorPower(userOne.address);
        expect(Number(userTwoPower)).to.equal(NEW_POWER);
      });
    });

    describe("Removing validators", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address, userTwo.address];
        state.initialPowers = [33, 21];

        // Deploy CosmosBridge contract
        state = (await setup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor,
        ) as TestValsetState);
      });

      it("should correctly update the valset when the operator removes a validator", async function () {
        // Confirm initial validator count
        const priorValsetValidatorCount = await state.cosmosBridge.validatorCount();
        expect(Number(priorValsetValidatorCount)).to.equal(state.initialValidators.length);

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(priorTotalPower)).to.equal(
          state.initialPowers[0] + state.initialPowers[1]
        );

        // Fail if not operator
        await expect(
          state.cosmosBridge.connect(userOne).removeValidator(userTwo.address)
        ).to.be.revertedWith("Must be the operator.");

        // Operator removes a validator
        await expect(state.cosmosBridge.connect(operator).removeValidator(userTwo.address)).to.not.be.reverted;

        // Confirm that userTwo is no longer an active validator
        const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(userTwo.address);
        expect(isUserTwoValidator).to.equal(false);

        // Confirm that userTwo's power has been reset
        const userTwoPower = await state.cosmosBridge.getValidatorPower(userTwo.address);
        expect(Number(userTwoPower)).to.equal(0);

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        expect(Number(postValsetValidatorCount)).to.equal(1);

        // Confirm updated total power
        const postTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(postTotalPower)).to.equal(state.initialPowers[0]);
      });

      it("should emit a LogValidatorRemoved event upon the removal of a validator", async function () {
        // Get the event logs from the update of a validator's power
        await state.cosmosBridge.connect(operator).removeValidator(userTwo.address);

        const userTwoActive = await state.cosmosBridge.isActiveValidator(userTwo.address);
        expect(userTwoActive).to.be.equal(false);
      });
    });

    describe("Updating the entire valset", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address, userTwo.address];
        state.initialPowers = [33, 21];

        state = (await setup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor,
        ) as TestValsetState);

        state.secondValidators = [userThree.address, accounts[4].address, accounts[5].address];
        state.secondPowers = [4, 19, 50];
      });

      it("should correctly update the valset", async function () {
        // Confirm current valset version number
        const priorValsetVersion = await state.cosmosBridge.currentValsetVersion();
        expect(Number(priorValsetVersion)).to.equal(1);

        // Confirm initial validator count
        const priorValsetValidatorCount = await state.cosmosBridge.validatorCount();
        expect(Number(priorValsetValidatorCount)).to.equal(state.initialValidators.length);

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(priorTotalPower)).to.equal(
          state.initialPowers[0] + state.initialPowers[1]
        );

        // Fail if not operator
        await expect(
          state.cosmosBridge
            .connect(userOne)
            .updateValset(state.secondValidators, state.secondPowers)
        ).to.be.revertedWith("Must be the operator.");

        // Operator resets the valset
        await expect(state.cosmosBridge
          .connect(operator)
          .updateValset(state.secondValidators, state.secondPowers)).to.not.be.reverted;

        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidator = await state.cosmosBridge.isActiveValidator(userOne.address);
        expect(isUserOneValidator).to.equal(false);

        const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(userTwo.address);
        expect(isUserTwoValidator).to.equal(false);

        // Confirm that all three secondary validators are now active validators
        const isUserThreeValidator = await state.cosmosBridge.isActiveValidator(userThree.address);
        expect(isUserThreeValidator).to.equal(true);
        const isUserFourValidator = await state.cosmosBridge.isActiveValidator(accounts[4].address);
        expect(isUserFourValidator).to.equal(true);
        const isUserFiveValidator = await state.cosmosBridge.isActiveValidator(accounts[5].address);
        expect(isUserFiveValidator).to.equal(true);

        // Confirm updated valset version number
        const postValsetVersion = await state.cosmosBridge.currentValsetVersion();
        expect(Number(postValsetVersion)).to.equal(2);

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        expect(Number(postValsetValidatorCount)).to.equal(state.secondValidators.length);

        // Confirm updated total power
        const postTotalPower = await state.cosmosBridge.totalPower();
        expect(Number(postTotalPower)).to.equal(
          state.secondPowers[0] + state.secondPowers[1] + state.secondPowers[2]
        );
      });

      it("should allow active validators to remain active if they are included in the new valset", async function () {
        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidatorFirstValsetVersion = await state.cosmosBridge.isActiveValidator(
          userOne.address
        );
        expect(isUserOneValidatorFirstValsetVersion).to.equal(true);

        // Operator resets the valset
        await expect(state.cosmosBridge
          .connect(operator)
          .updateValset([state.initialValidators[0]], [state.initialPowers[0]])).to.not.be.reverted;

        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidatorSecondValsetVersion = await state.cosmosBridge.isActiveValidator(
          userOne.address
        );
        expect(isUserOneValidatorSecondValsetVersion).to.equal(true);
      });

      it("should emit LogValsetReset and LogValsetUpdated events upon the update of the valset", async function () {
        // Get the event logs from the valset update
        await expect(state.cosmosBridge
          .connect(operator)
          .updateValset(state.secondValidators, state.secondPowers)).to.not.be.reverted;

        for (let i = 0; i < state.secondValidators.length; i++) {
          const isWhitelisted = await state.cosmosBridge.isActiveValidator(
            state.secondValidators[i]
          );

          const validatorPower = await state.cosmosBridge.getValidatorPower(
            state.secondValidators[i]
          );

          expect(isWhitelisted).to.equal(true);
          expect(Number(validatorPower)).to.equal(state.secondPowers[i]);
        }
      });
    });
  });

  describe("Gas recovery", function () {
    beforeEach(async function () {
      state.initialValidators = [userOne.address, userTwo.address];
      state.initialPowers = [50, 60];

      state = (await setup(
        state.initialValidators,
        state.initialPowers,
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser,
        networkDescriptor,
      ) as TestValsetState);

      state.secondValidators = [userThree.address];
      state.secondPowers = [5];
    });

    it("should not allow the gas recovery of storage in use by active validators", async function () {
      // Operator attempts to recover gas from userOne's storage slot
      await expect(state.cosmosBridge
        .connect(operator)
        .recoverGas(1, userOne.address))
        .to.be.revertedWith(EVMRevert);
    });

    it("should allow the gas recovery of inactive validator storage", async function () {
      // Confirm that both initial validators are active validators
      const isUserOneValidatorPrior = await state.cosmosBridge.isActiveValidator(userOne.address);
      expect(isUserOneValidatorPrior).to.equal(true);
      const isUserTwoValidatorPrior = await state.cosmosBridge.isActiveValidator(userTwo.address);
      expect(isUserTwoValidatorPrior).to.equal(true);

      // Operator updates the valset, making userOne and userTwo inactive validators
      await expect(state.cosmosBridge
        .connect(operator)
        .updateValset(state.secondValidators, state.secondPowers)).not.be.reverted;

      // Confirm that both initial validators are no longer an active validators
      const isUserOneValidatorPost = await state.cosmosBridge.isActiveValidator(userOne.address);
      expect(isUserOneValidatorPost).to.equal(false);
      const isUserTwoValidatorPost = await state.cosmosBridge.isActiveValidator(userTwo.address);
      expect(isUserTwoValidatorPost).to.equal(false);

      // Fail if not operator
      await expect(
        state.cosmosBridge.connect(userTwo).recoverGas(1, userOne.address)
      ).to.be.revertedWith("Must be the operator.");

      // Operator recovers gas from inactive validator userOne
      await expect(state.cosmosBridge.connect(operator).recoverGas(1, userOne.address)).to.not.be.reverted;

      // Operator recovers gas from inactive validator userOne
      await expect(state.cosmosBridge.connect(operator).recoverGas(1, userTwo.address)).to.not.be.reverted;
    });
  });
});
