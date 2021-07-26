const { ethers } = require("hardhat");
const web3 = require("web3");
const { use, expect } = require("chai");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

const { singleSetup } = require('./helpers/testFixture');

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

describe("Test Valset", function () {
  let userOne;
  let userTwo;
  let userThree;
  let userFour;
  let accounts;
  let signerAccounts;
  let operator;
  let owner;
  const consensusThreshold = 80;
  let initialPowers;
  let initialValidators;
  let networkDescriptor;
  // track the state of the deployed contracts
  let state;

  describe("Valset contract deployment", function () {

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
      let initialValidators = [userOne.address, userTwo.address, userThree];
      let initialPowers = [5, 8, 12];
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

    it("should deploy the Valset and correctly set the current valset version", async function () {
      state.cosmosBridge.should.exist;

      const valsetValsetVersion = await state.cosmosBridge.currentValsetVersion();
      Number(valsetValsetVersion).should.be.bignumber.equal(1);
    });

    it("should correctly set initial validators and initial validator count", async function () {
      const userOneValidator = await state.cosmosBridge.isActiveValidator(
        userOne.address
      );
      const userTwoValidator = await state.cosmosBridge.isActiveValidator(
        userTwo.address
      );
      const userThreeValidator = await state.cosmosBridge.isActiveValidator(
        userThree
      );
      const valsetValidatorCount = await state.cosmosBridge.validatorCount();

      userOneValidator.should.be.equal(true);
      userTwoValidator.should.be.equal(true);
      userThreeValidator.should.be.equal(true);
      Number(valsetValidatorCount).should.be.bignumber.equal(
        state.initialValidators.length
      );
    });

    it("should correctly set initial validator powers ", async function () {
      const userOnePower = await state.cosmosBridge.getValidatorPower(userOne.address);
      const userTwoPower = await state.cosmosBridge.getValidatorPower(userTwo.address);
      const userThreePower = await state.cosmosBridge.getValidatorPower(
        userThree
      );

      Number(userOnePower).should.be.bignumber.equal(state.initialPowers[0]);
      Number(userTwoPower).should.be.bignumber.equal(state.initialPowers[1]);
      Number(userThreePower).should.be.bignumber.equal(state.initialPowers[2]);
    });

    it("should correctly set the initial total power", async function () {
      const valsetTotalPower = await state.cosmosBridge.totalPower();

      Number(valsetTotalPower).should.be.bignumber.equal(
        state.initialPowers[0] + state.initialPowers[1] + state.initialPowers[2]
      );
    });
  });

  describe("Dynamic validator set", function () {
    describe("Adding validators", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address];
        state.initialPowers = [5];

        state = await singleSetup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor
        );

        state.userTwoPower = 11;
        state.userThreePower = 44;
      });
      
      it("should correctly update the valset when the operator adds a new validator", async function () {
        // Confirm initial validator count
        const priorValsetValidatorCount = await state.cosmosBridge.validatorCount();
        Number(priorValsetValidatorCount).should.be.bignumber.equal(1);

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();

        Number(priorTotalPower).should.be.bignumber.equal(
          state.initialPowers[0]
        );

        // Operator adds a validator
        await state.cosmosBridge.connect(operator)
          .addValidator(userTwo.address, state.userTwoPower)
          .should.be.fulfilled;

        // Confirm that userTwo has been set as a validator
        const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(
          userTwo.address
        );
        isUserTwoValidator.should.be.equal(true);

        // Confirm that userTwo's power has been correctly set
        const userTwoSetPower = await state.cosmosBridge.getValidatorPower(
          userTwo.address
        );
        Number(userTwoSetPower).should.be.bignumber.equal(state.userTwoPower);

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(2);

        // Confirm updated total power
        const postTotalPower = await state.cosmosBridge.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(
          state.initialPowers[0] + state.userTwoPower
        );
      });

      it("should be able to add a new validator and get its power", async function () {
        // Get the event logs from the addition of a new validator
        await state.cosmosBridge.connect(operator).addValidator(
          userTwo.address,
          state.userTwoPower
        );

        const userTwoSetPower = await state.cosmosBridge.getValidatorPower(
          userTwo.address
        );
        Number(userTwoSetPower).should.be.bignumber.equal(state.userTwoPower);
      });

      it("should allow the operator to add multiple new validators", async function () {
        // Fail if not operator
        await expect(
            state.cosmosBridge.connect(userOne)
              .addValidator(userTwo.address, state.userTwoPower),
        ).to.be.revertedWith("Must be the operator.");

        await state.cosmosBridge.connect(operator)
          .addValidator(userTwo.address, state.userTwoPower)
          .should.be.fulfilled;
        await state.cosmosBridge.connect(operator)
          .addValidator(userThree, state.userThreePower)
          .should.be.fulfilled;
        await state.cosmosBridge.connect(operator)
          .addValidator(accounts[4].address, 77)
          .should.be.fulfilled;
        await state.cosmosBridge.connect(operator)
          .addValidator(accounts[5].address, 23)
          .should.be.fulfilled;

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(5);

        // Confirm updated total power
        const valsetTotalPower = await state.cosmosBridge.totalPower();
        Number(valsetTotalPower).should.be.bignumber.equal(
          state.initialPowers[0] + state.userTwoPower + state.userThreePower + 100 // (23 + 77)
        );
      });
    });

    describe("Updating validator's power", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address];
        state.initialPowers = [5];

        // Deploy CosmosBridge contract
        state = await singleSetup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor
        );

        state.userTwoPower = 11;
        state.userThreePower = 44;
      });

      it("should allow the operator to update a validator's power", async function () {
        const NEW_POWER = 515;

        // Confirm userOne's initial power
        const userOneInitialPower = await state.cosmosBridge.getValidatorPower(
          userOne.address
        );

        Number(userOneInitialPower).should.be.bignumber.equal(
          state.initialPowers[0]
        );

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();
        Number(priorTotalPower).should.be.bignumber.equal(
          state.initialPowers[0]
        );

        // Fail if not operator
        await expect(
          state.cosmosBridge.connect(userTwo)
            .updateValidatorPower(userOne.address, NEW_POWER),
        ).to.be.revertedWith("Must be the operator.");

        // Operator updates the validator's initial power
        await state.cosmosBridge.connect(operator)
          .updateValidatorPower(userOne.address, NEW_POWER)
          .should.be.fulfilled;

        // Confirm userOne's power has increased
        const userOnePostPower = await state.cosmosBridge.getValidatorPower(
          userOne.address
        );
        Number(userOnePostPower).should.be.bignumber.equal(NEW_POWER);

        // Confirm total power has been updated
        const postTotalPower = await state.cosmosBridge.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(NEW_POWER);
      });

      it("should update of a validator's power", async function () {
        const NEW_POWER = 111;

        await state.cosmosBridge.connect(operator)
          .updateValidatorPower(userOne.address, NEW_POWER);

        const userTwoPower = await state.cosmosBridge.getValidatorPower(userOne.address);
        Number(userTwoPower).should.be.bignumber.equal(NEW_POWER);
      });
    });

    describe("Removing validators", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address, userTwo.address];
        state.initialPowers = [33, 21];

        // Deploy CosmosBridge contract
        state = await singleSetup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor
        );
      });

      it("should correctly update the valset when the operator removes a validator", async function () {
        // Confirm initial validator count
        const priorValsetValidatorCount = await state.cosmosBridge.validatorCount();
        Number(priorValsetValidatorCount).should.be.bignumber.equal(
          state.initialValidators.length
        );

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();
        Number(priorTotalPower).should.be.bignumber.equal(
          state.initialPowers[0] + state.initialPowers[1]
        );

        // Fail if not operator
        await expect(
          state.cosmosBridge.connect(userOne)
            .removeValidator(userTwo.address),
        ).to.be.revertedWith("Must be the operator.");

        // Operator removes a validator
        await state.cosmosBridge.connect(operator)
          .removeValidator(userTwo.address)
          .should.be.fulfilled;

        // Confirm that userTwo is no longer an active validator
        const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(userTwo.address);
        isUserTwoValidator.should.be.equal(false);

        // Confirm that userTwo's power has been reset
        const userTwoPower = await state.cosmosBridge.getValidatorPower(userTwo.address);
        Number(userTwoPower).should.be.bignumber.equal(0);

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(1);

        // Confirm updated total power
        const postTotalPower = await state.cosmosBridge.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(state.initialPowers[0]);
      });

      it("should emit a LogValidatorRemoved event upon the removal of a validator", async function () {
        // Get the event logs from the update of a validator's power
        await state.cosmosBridge.connect(operator)
          .removeValidator(userTwo.address);

        const userTwoActive = await state.cosmosBridge.isActiveValidator(userTwo.address);
        expect(userTwoActive).to.be.equal(false);
      });
    });

    describe("Updating the entire valset", function () {
      beforeEach(async function () {
        state.initialValidators = [userOne.address, userTwo.address];
        state.initialPowers = [33, 21];

        state = await singleSetup(
          state.initialValidators,
          state.initialPowers,
          operator,
          consensusThreshold,
          owner,
          userOne,
          userThree,
          pauser,
          networkDescriptor
        );

        state.secondValidators = [userThree, accounts[4].address, accounts[5].address];
        state.secondPowers = [4, 19, 50];
      });

      it("should correctly update the valset", async function () {
        // Confirm current valset version number
        const priorValsetVersion = await state.cosmosBridge.currentValsetVersion();
        Number(priorValsetVersion).should.be.bignumber.equal(1);

        // Confirm initial validator count
        const priorValsetValidatorCount = await state.cosmosBridge.validatorCount();
        Number(priorValsetValidatorCount).should.be.bignumber.equal(
          state.initialValidators.length
        );

        // Confirm initial total power
        const priorTotalPower = await state.cosmosBridge.totalPower();
        Number(priorTotalPower).should.be.bignumber.equal(
          state.initialPowers[0] + state.initialPowers[1]
        );

        // Fail if not operator
        await expect(
          state.cosmosBridge.connect(userOne).updateValset(
            state.secondValidators,
            state.secondPowers,
          ),
        ).to.be.revertedWith("Must be the operator.");

        // Operator resets the valset
        await state.cosmosBridge.connect(operator).updateValset(
          state.secondValidators,
          state.secondPowers
        ).should.be.fulfilled;

        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidator = await state.cosmosBridge.isActiveValidator(
          userOne.address
        );
        isUserOneValidator.should.be.equal(false);

        const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(
          userTwo.address
        );
        isUserTwoValidator.should.be.equal(false);

        // Confirm that all three secondary validators are now active validators
        const isUserThreeValidator = await state.cosmosBridge.isActiveValidator(
          userThree
        );
        isUserThreeValidator.should.be.equal(true);
        const isUserFourValidator = await state.cosmosBridge.isActiveValidator(
          accounts[4].address
        );
        isUserFourValidator.should.be.equal(true);
        const isUserFiveValidator = await state.cosmosBridge.isActiveValidator(
          accounts[5].address
        );
        isUserFiveValidator.should.be.equal(true);

        // Confirm updated valset version number
        const postValsetVersion = await state.cosmosBridge.currentValsetVersion();
        Number(postValsetVersion).should.be.bignumber.equal(2);

        // Confirm updated validator count
        const postValsetValidatorCount = await state.cosmosBridge.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(
          state.secondValidators.length
        );

        // Confirm updated total power
        const postTotalPower = await state.cosmosBridge.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(
          state.secondPowers[0] + state.secondPowers[1] + state.secondPowers[2]
        );
      });

      it("should allow active validators to remain active if they are included in the new valset", async function () {
        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidatorFirstValsetVersion = await state.cosmosBridge.isActiveValidator(
          userOne.address
        );
        isUserOneValidatorFirstValsetVersion.should.be.equal(true);

        // Operator resets the valset
        await state.cosmosBridge.connect(operator).updateValset(
          [state.initialValidators[0]],
          [state.initialPowers[0]],
        ).should.be.fulfilled;

        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidatorSecondValsetVersion = await state.cosmosBridge.isActiveValidator(
          userOne.address
        );
        isUserOneValidatorSecondValsetVersion.should.be.equal(true);
      });

      it("should emit LogValsetReset and LogValsetUpdated events upon the update of the valset", async function () {
        // Get the event logs from the valset update
        await state.cosmosBridge.connect(operator).updateValset(
          state.secondValidators,
          state.secondPowers,
        ).should.be.fulfilled;

        for (let i = 0; i < state.secondValidators.length; i++) {
          const isWhitelisted = await state.cosmosBridge
            .isActiveValidator(state.secondValidators[i]);

          const validatorPower = await state.cosmosBridge
            .getValidatorPower(state.secondValidators[i]);

          expect(isWhitelisted).to.be.equal(true);
          expect(Number(validatorPower)).to.be.equal(state.secondPowers[i]);
        }
      });
    });
  });

  describe("Gas recovery", function () {
    beforeEach(async function () {
      state.initialValidators = [userOne.address, userTwo.address];
      state.initialPowers = [50, 60];

      state = await singleSetup(
        state.initialValidators,
        state.initialPowers,
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser,
        networkDescriptor
      );

      state.secondValidators = [userThree];
      state.secondPowers = [5];
    });

    it("should not allow the gas recovery of storage in use by active validators", async function () {
      // Operator attempts to recover gas from userOne's storage slot
      await state.cosmosBridge.connect(operator)
        .recoverGas(1, userOne.address)
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow the gas recovery of inactive validator storage", async function () {
      // Confirm that both initial validators are active validators
      const isUserOneValidatorPrior = await state.cosmosBridge.isActiveValidator(
        userOne.address
      );
      isUserOneValidatorPrior.should.be.equal(true);
      const isUserTwoValidatorPrior = await state.cosmosBridge.isActiveValidator(
        userTwo.address
      );
      isUserTwoValidatorPrior.should.be.equal(true);

      // Operator updates the valset, making userOne and userTwo inactive validators
      await state.cosmosBridge.connect(operator)
        .updateValset(state.secondValidators, state.secondPowers)
        .should.be.fulfilled;

      // Confirm that both initial validators are no longer an active validators
      const isUserOneValidatorPost = await state.cosmosBridge.isActiveValidator(
        userOne.address
      );
      isUserOneValidatorPost.should.be.equal(false);
      const isUserTwoValidatorPost = await state.cosmosBridge.isActiveValidator(
        userTwo.address
      );
      isUserTwoValidatorPost.should.be.equal(false);

      // Fail if not operator
      await expect(
        state.cosmosBridge.connect(userTwo)
        .recoverGas(1, userOne.address),
      ).to.be.revertedWith("Must be the operator.");

      // Operator recovers gas from inactive validator userOne
      await state.cosmosBridge.connect(operator)
        .recoverGas(1, userOne.address)
        .should.be.fulfilled;

      // Operator recovers gas from inactive validator userOne
      await state.cosmosBridge.connect(operator)
        .recoverGas(1, userTwo.address)
        .should.be.fulfilled;
    });
  });
});