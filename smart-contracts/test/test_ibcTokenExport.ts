import Web3Utils from "web3-utils";
import web3 from "web3";

import { ethers } from "hardhat";
import { use, expect } from "chai";
import { solidity } from "ethereum-waffle";
import { setup, getValidClaim, TestFixtureState } from "./helpers/testFixture";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";

const BigNumber = ethers.BigNumber;
const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

use(solidity);

describe("Test Cosmos Bridge", function () {
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
    userThree = accounts[9];

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
      true,
    );
  });


  it("should deploy a new token upon the successful processing of a ibc token burn prophecy claim just for first time", async function () {
    state.nonce = 1;

    const { digest, claimData, signatures } = await getValidClaim(
      state.sender,
      state.senderSequence,
      state.recipient.address,
      ZERO_ADDRESS,
      state.amount,
      state.name,
      state.symbol,
      state.decimals,
      state.networkDescriptor,
      true,
      state.nonce,
      state.constants.denom.ibc,
      [userOne, userTwo, userFour],
    );

    const expectedAddress = ethers.utils.getContractAddress({
      from: state.bridgeBank.address,
      nonce: 1,
    });
    await expect(
      state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
    )
      .to.emit(state.cosmosBridge, "LogNewBridgeTokenCreated")
      .withArgs(
        state.decimals,
        state.networkDescriptor,
        state.name,
        state.symbol,
        ZERO_ADDRESS,
        expectedAddress,
        state.constants.denom.ibc
      );

    const newlyCreatedTokenAddress = await state.cosmosBridge.cosmosDenomToDestinationAddress(
      state.constants.denom.ibc
    );
    expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

    // Everything again, but this time submitProphecyClaimAggregatedSigs should NOT emit the event
    const {
      digest: digest2,
      claimData: claimData2,
      signatures: signatures2,
    } = await getValidClaim(
      state.sender,
      state.senderSequence + 1,
      state.recipient.address,
      ZERO_ADDRESS,
      state.amount,
      state.name,
      state.symbol,
      state.decimals,
      state.networkDescriptor,
      true,
      state.nonce + 1,
      state.constants.denom.ibc,
      [userOne, userTwo, userFour],
    );

    await expect(
      state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest2, claimData2, signatures2)
    ).to.not.emit(state.cosmosBridge, "LogNewBridgeTokenCreated");

    // But should have minted:
    const deployedTokenFactory = await ethers.getContractFactory("BridgeToken");
    const deployedToken = await deployedTokenFactory.attach(newlyCreatedTokenAddress);
    const endingBalance = await deployedToken.balanceOf(state.recipient.address);
    expect(endingBalance).to.be.equal(state.amount * 2);
  });
});
