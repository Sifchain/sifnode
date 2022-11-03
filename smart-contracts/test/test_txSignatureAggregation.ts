import { setup, getValidClaim, TestFixtureState, prefundAccount, preApproveAccount } from "./helpers/testFixture";

import { colorLog } from "./helpers/helpers";

import { expect } from "chai";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { ContractTransaction } from "ethers";

// Set `use` to `true` to compare a new implementation with the previous gas costs;
// Please set the previous gas costs accordingly in this object
const gasProfiling = {
  use: false,
  lock: 176141,
  mint: 182155,
  newBt: 1665776,
  multiLock: 346709,
  current: {
    lock: 0,
    mint:0,
    newBt: 0
  },
};

describe("Gas Cost Tests", function () {
  let userOne: SignerWithAddress;
  let userTwo: SignerWithAddress;
  let userThree: SignerWithAddress;
  let userFour: SignerWithAddress;
  let accounts: SignerWithAddress[];
  let operator: SignerWithAddress;
  let owner: SignerWithAddress;
  let pauser: SignerWithAddress;

  // Consensus threshold of 70%
  const consensusThreshold = 70;
  let initialPowers: number[];
  let initialValidators: string[];
  let networkDescriptor: number;
  let state: TestFixtureState;

  before(async function () {
    accounts = await ethers.getSigners();

    operator = accounts[0];
    userOne = accounts[1];
    userTwo = accounts[2];
    userThree = accounts[3];
    userFour = accounts[4];

    owner = accounts[5];
    pauser = accounts[6];

    initialPowers = [25, 25, 25, 25];
    initialValidators = [userOne.address, userTwo.address, userThree.address, userFour.address];

    networkDescriptor = 1;
  });

  beforeEach(async function () {
    // Deploy Valset contract
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
    
    // Send UserOne an initial balance
    const tokens = [state.token, state.token1, state.token2, state.token3, state.token_ibc, state.token_noDenom]
    await prefundAccount(userOne, state.amount, state.operator, tokens);
    await preApproveAccount(state.bridgeBank, userOne, state.amount, tokens);

    // Lock tokens on contract
    await expect(state.bridgeBank
      .connect(userOne)
      .lock(state.sender, state.token1.address, state.amount)).to.not.be.reverted;
  });

  describe("Gas Cost With 4 Validators", function () {
    it("should allow us to check the cost of submitting a prophecy claim lock", async function () {
      let balance = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(0);

      // Last nonce should now be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token1.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.one,
        accounts.slice(1, 5),
      );

      const tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);

      if (gasProfiling.use) {
        gasProfiling.current.lock = sum;
        logGasDiff("LOCK:", gasProfiling.lock, sum);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", sum);
      }

      // Bridge claim should be completed
      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      balance = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(state.amount);
    });

    it("should allow us to check the cost of submitting a prophecy claim mint", async function () {
      let balance = Number(await state.rowan.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(0);

      // Last nonce should now be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.rowan.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.rowan,
        accounts.slice(1, 5),
      );

      const tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);

      if (gasProfiling.use) {
        gasProfiling.current.mint = sum;
        logGasDiff("MINT:", gasProfiling.mint, sum);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", sum);
      }

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      // balance should have increased
      balance = Number(await state.rowan.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(state.amount);
    });

    it("should allow us to check the cost of creating a new BridgeToken", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token1.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce,
        state.constants.denom.one,
        accounts.slice(1, 5),
      );

      const expectedAddress = ethers.utils.getContractAddress({
        from: state.bridgeBank.address,
        nonce: 1,
      });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);

      if (gasProfiling.use) {
        gasProfiling.current.newBt = sum;
        logGasDiff("DoublePeg :: New BridgeToken:", gasProfiling.newBt, sum);
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", Number(receipt.gasUsed));
      }

      const newlyCreatedTokenAddress = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.constants.denom.one
      );
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

      // expect the token to have a denom
      const registeredDenom = await state.bridgeBank.contractDenom(newlyCreatedTokenAddress);
      expect(registeredDenom).to.be.equal(state.constants.denom.one);
    });

    it("should allow us to check the cost of submitting a batch prophecy claim lock", async function () {
      // Lock token2 on contract
      await expect(state.bridgeBank.connect(userOne).lock(state.sender, state.token2.address, state.amount))
        .not.to.be.reverted;

      // Lock token3 on contract
      await expect(state.bridgeBank.connect(userOne).lock(state.sender, state.token3.address, state.amount))
        .not.to.be.reverted;

      let balanceToken1 = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balanceToken1).to.equal(0);

      let balanceToken2 = Number(await state.token2.balanceOf(state.recipient.address));
      expect(balanceToken2).to.equal(0);

      let balanceToken3 = Number(await state.token3.balanceOf(state.recipient.address));
      expect(balanceToken3).to.equal(0);

      // Last nonce should now be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token1.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.one,
        accounts.slice(1, 5),
      );

      const {
        digest: digest2,
        claimData: claimData2,
        signatures: signatures2,
      } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token2.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce + 1,
        state.constants.denom.two,
        accounts.slice(1, 5),
      );

      const {
        digest: digest3,
        claimData: claimData3,
        signatures: signatures3,
      } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token3.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce + 2,
        state.constants.denom.three,
        accounts.slice(1, 5),
      );

      const tx = await state.cosmosBridge
        .connect(userOne)
        .batchSubmitProphecyClaimAggregatedSigs(
          [digest, digest2, digest3],
          [claimData, claimData2, claimData3],
          [signatures, signatures2, signatures3]
        );
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);

      if (gasProfiling.use) {
        const numberOfClaimsInBatch = 3;
        logGasDiff("BATCH, regarding previous implementation:", gasProfiling.multiLock, sum);
        logGasDiff(
          `${numberOfClaimsInBatch} Batched claims VS single claim:`,
          gasProfiling.current.lock * numberOfClaimsInBatch,
          sum,
          true
        );
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", sum);
      }

      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(3);

      balanceToken1 = Number(await state.token1.balanceOf(state.recipient.address));
      expect(balanceToken1).to.be.equal(state.amount);

      balanceToken2 = Number(await state.token2.balanceOf(state.recipient.address));
      expect(balanceToken2).to.be.equal(state.amount);

      balanceToken3 = Number(await state.token3.balanceOf(state.recipient.address));
      expect(balanceToken3).to.be.equal(state.amount);
    });

    it("should allow us to check the cost of batch creating new BridgeTokens", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token1.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce,
        state.constants.denom.one,
        accounts.slice(1, 5),
      );

      const {
        digest: digest2,
        claimData: claimData2,
        signatures: signatures2,
      } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token2.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce + 1,
        state.constants.denom.two,
        accounts.slice(1, 5),
      );

      const {
        digest: digest3,
        claimData: claimData3,
        signatures: signatures3,
      } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token3.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce + 2,
        state.constants.denom.three,
        accounts.slice(1, 5),
      );

      const expectedAddress1 = ethers.utils.getContractAddress({
        from: state.bridgeBank.address,
        nonce: 1,
      });
      const expectedAddress2 = ethers.utils.getContractAddress({
        from: state.bridgeBank.address,
        nonce: 2,
      });
      const expectedAddress3 = ethers.utils.getContractAddress({
        from: state.bridgeBank.address,
        nonce: 3,
      });

      const tx = await state.cosmosBridge
        .connect(userOne)
        .batchSubmitProphecyClaimAggregatedSigs(
          [digest, digest2, digest3],
          [claimData, claimData2, claimData3],
          [signatures, signatures2, signatures3]
        );
      const receipt = await tx.wait();
      const sum = Number(receipt.gasUsed);

      if (gasProfiling.use) {
        const numberOfClaimsInBatch = 3;
        logGasDiff(
          `DoublePeg :: ${numberOfClaimsInBatch} Batched claims VS single claim:`,
          gasProfiling.current.newBt * numberOfClaimsInBatch,
          sum
        );
      } else {
        console.log("~~~~~~~~~~~~\nTotal: ", Number(receipt.gasUsed));
      }

      const newlyCreatedTokenAddress1 = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.constants.denom.one
      );
      expect(newlyCreatedTokenAddress1).to.be.equal(expectedAddress1);

      const newlyCreatedTokenAddress2 = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.constants.denom.two
      );
      expect(newlyCreatedTokenAddress2).to.be.equal(expectedAddress2);

      const newlyCreatedTokenAddress3 = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.constants.denom.three
      );
      expect(newlyCreatedTokenAddress3).to.be.equal(expectedAddress3);
    });
  });
});

// Helper function to aid comparing implementations wrt gas costs
function logGasDiff(title: string, original: number, current: number, useShortTitle?: boolean) {
  const separator = useShortTitle ? "---" : "~~~~~~~~~~~~";
  colorLog("cyan", `${separator}\n${title}`);
  console.log("Original:", original);
  console.log("Current :", current);
  const pct = Math.abs((1 - current / original) * 100).toFixed(2);
  const diff = current - original;
  colorLog(getColorName(diff), `Diff    : ${diff} (${pct}%)`);
}

function getColorName(value: number) {
  if (value > 0) {
    return "red";
  } else if (value < 0) {
    return "green";
  } else {
    return "white";
  }
}

/**
 * 
 * 
Unlock Gas Cost With 4 Validators
tx0  173978
~~~~~~~~~~~~
Total:  173978

Mint Gas Cost With 4 Validators
tx0  179749
~~~~~~~~~~~~
Total:  179749

Create new BridgeToken Gas Cost With 4 Validators
tx0  1162769
~~~~~~~~~~~~
Total:  1162769
 * 
 */
