import Web3Utils from "web3-utils";
import web3 from "web3";
import { ethers, network } from "hardhat";
import { use, expect } from "chai";
import { solidity } from "ethereum-waffle";
import { setup, getValidClaim, TestFixtureState, SignedData } from "./helpers/testFixture";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";

use(solidity);
const BigNumber = ethers.BigNumber;

const getBalance = async function (address: string, fromWei = false): Promise<string> {
  let balance = await network.provider.send("eth_getBalance", [address]);

  if (fromWei) {
    balance = Web3Utils.fromWei(balance);
  }

  return balance;
};

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
      false,
    );
  });

  describe("CosmosBridge:Oracle", function () {
    it("should be able to getProphecyStatus with 1/3 of power", async function () {
      const status = await state.cosmosBridge.connect(userOne).getProphecyStatus(25);
      expect(status).to.equal(false);
    });

    it("should be able to getProphecyStatus with 2/3 of power", async function () {
      const status = await state.cosmosBridge.getProphecyStatus(50);
      expect(status).to.equal(false);
    });

    it("should be able to getProphecyStatus with 3/3 of power", async function () {
      const status = await state.cosmosBridge.getProphecyStatus(75);
      expect(status).to.equal(true);
    });

    it("should be able to getProphecyStatus with 100% of power", async function () {
      const status = await state.cosmosBridge.getProphecyStatus(100);
      expect(status).to.equal(true);
    });

    it("should allow operator to update consensus threshold", async function () {
      const newThreshold = 35;

      let status = await state.cosmosBridge.getProphecyStatus(50);
      expect(status).to.equal(false, "Expected not to pass, 50 is below default of 75");

      await expect(state.cosmosBridge.connect(operator).updateConsensusThreshold(newThreshold)).to
        .not.be.reverted;

      status = await state.cosmosBridge.getProphecyStatus(50);
      expect(status).to.equal(
        true,
        "signedPower of 50 should now pass because it is above newThreshold"
      );
    });

    it("should not allow non-operator to update consensus threshold", async function () {
      await expect(
        state.cosmosBridge.connect(userOne).updateConsensusThreshold(50)
      ).to.be.revertedWith("Must be the operator.");
    });

    it("should allows updating consensus threshold to 100", async function () {
      const newThreshold = 100;
      let status = await state.cosmosBridge.connect(operator).getProphecyStatus(newThreshold);
      expect(status).to.equal(true);
    });

    it("should not allow updating consensus threshold to lte 0", async function () {
      const newThreshold = 0;
      await expect(
        state.cosmosBridge.connect(operator).updateConsensusThreshold(newThreshold)
      ).to.be.revertedWith("Consensus threshold must be positive.");
    });

    it("should not allow updating consensus threshold to value greater than 100", async function () {
      const newThreshold = 101;
      await expect(
        state.cosmosBridge.connect(operator).updateConsensusThreshold(newThreshold)
      ).to.be.revertedWith("Invalid consensus threshold.");
    });

  });

  describe("CosmosBridge", function () {
    it("Can update the valset", async function () {
      // Operator resets the valset
      await expect(state.cosmosBridge
        .connect(operator)
        .updateValset([userOne.address, userTwo.address], [50, 50])).not.to.be.reverted;

      // Confirm that both initial validators are now active validators
      const isUserOneValidator = await state.cosmosBridge.isActiveValidator(userOne.address);
      expect(isUserOneValidator).to.equal(true);
      const isUserTwoValidator = await state.cosmosBridge.isActiveValidator(userTwo.address);
      expect(isUserTwoValidator).to.equal(true);

      // Confirm that all both secondary validators are not active validators
      const isUserThreeValidator = await state.cosmosBridge.isActiveValidator(userThree.address);
      expect(isUserThreeValidator).to.equal(false);
      const isUserFourValidator = await state.cosmosBridge.isActiveValidator(userFour.address);
      expect(isUserFourValidator).to.equal(false);
    });

    it("Can change the operator", async function () {
      // Confirm that the operator has changed
      const originalOperator = await state.cosmosBridge.operator();
      expect(originalOperator).to.be.equal(operator.address);

      // Operator resets the valset
      await expect(state.cosmosBridge.connect(operator).changeOperator(userOne.address)).not.to.be.reverted;

      // Confirm that the operator has changed
      const newOperator = await state.cosmosBridge.operator();
      expect(newOperator).to.be.equal(userOne.address);
    });

    it("should NOT allow to change the operator to the zero address", async function () {
      // Confirm that the operator has changed
      const originalOperator = await state.cosmosBridge.operator();
      expect(originalOperator).to.be.equal(operator.address);

      // Operator resets the valset
      await expect(
        state.cosmosBridge.connect(operator).changeOperator(state.constants.zeroAddress)
      ).to.be.revertedWith("invalid address");

      // Confirm that the operator has NOT changed
      const newOperator = await state.cosmosBridge.operator();
      expect(newOperator).to.be.equal(operator.address);
    });

    it("Can update the validator set", async function () {
      // Also make sure everything runs fourth time after switching validators a second time.
      // Operator resets the valset
      await expect(state.cosmosBridge
        .connect(operator)
        .updateValset([userThree.address, userFour.address], [50, 50])).not.to.be.reverted;

      // Confirm that both initial validators are no longer an active validators
      const isUserOneValidator2 = await state.cosmosBridge.isActiveValidator(userOne.address);
      expect(isUserOneValidator2).to.equal(false);
      const isUserTwoValidator2 = await state.cosmosBridge.isActiveValidator(userTwo.address);
      expect(isUserTwoValidator2).to.equal(false);

      // Confirm that both secondary validators are now active validators
      const isUserThreeValidator2 = await state.cosmosBridge.isActiveValidator(userThree.address);
      expect(isUserThreeValidator2).to.equal(true);
      const isUserFourValidator2 = await state.cosmosBridge.isActiveValidator(userFour.address);
      expect(isUserFourValidator2).to.equal(true);
    });

    it("should return true if a sifchain address prefix is correct", async function () {
      expect(await state.bridgeBank.VSA(state.sender)).to.equal(true);
    });

    it("should return false if a sifchain address length is incorrect", async function () {
      const incorrectSifAddress = web3.utils.utf8ToHex(
        "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpaceee"
      );
      expect(await state.bridgeBank.VSA(incorrectSifAddress)).to.equal(false);
    });

    it("should return false if a sifchain address has an incorrect `sif` prefix", async function () {
      const incorrectSifAddress = web3.utils.utf8ToHex(
        "eif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
      );
      expect(await state.bridgeBank.VSA(incorrectSifAddress)).to.equal(false);
    });

    it("should deploy cosmos bridge and bridge bank", async function () {
      expect((await state.cosmosBridge.consensusThreshold()).toString()).to.equal(
        consensusThreshold.toString()
      );

      // iterate over all validators and ensure they have the proper
      // powers and that they have been succcessfully whitelisted
      for (let i = 0; i < initialValidators.length; i++) {
        const address = initialValidators[i];

        expect(await state.cosmosBridge.isActiveValidator(address)).to.be.true;

        expect((await state.cosmosBridge.getValidatorPower(address)).toString()).to.equal("25");
      }

      expect(await state.bridgeBank.cosmosBridge()).to.be.equal(state.cosmosBridge.address);
      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
    });

    it("should deploy cosmos bridge and bridge bank, correctly setting the networkDescriptor", async function () {
      expect(await state.cosmosBridge.networkDescriptor()).to.equal(state.networkDescriptor);
      expect(await state.bridgeBank.networkDescriptor()).to.equal(state.networkDescriptor);
    });

    it("should unlock tokens upon the successful processing of a burn prophecy claim", async function () {
      // Bridgebank should have a balance of tokens that would be unlocked when processing the claim
      await state.token.connect(operator).mint(state.bridgeBank.address, state.amount);
      
      const beforeUserBalance = Number(await state.token.balanceOf(state.recipient.address));
      expect(beforeUserBalance).to.equal(Number(0));

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.one,
        [userOne, userTwo, userFour],
      );

      await expect(state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)).not.to.be.reverted;

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.equal(1);

      let balance = Number(await state.token.balanceOf(state.recipient.address));
      expect(balance).to.equal(state.amount);
    });

    it("should NOT unlock tokens upon the successful processing of a burn prophecy claim if the recipient is blocklisted", async function () {
      // Add recipient to the blocklist
      await expect(state.blocklist.addToBlocklist(state.recipient.address)).to.be.not.be.reverted;

      const beforeUserBalance = Number(await state.token.balanceOf(state.recipient.address));
      expect(beforeUserBalance).to.equal(Number(0));

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.one,
        [userOne, userTwo, userFour],
      );

      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      // Last nonce should now be 1
      lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);

      // Balance should still be 0; tokens should not be unlocked
      const balance = Number(await state.token.balanceOf(state.recipient.address));
      expect(balance).to.be.equal(0);
    });

    it("should unlock eth upon the successful processing of a burn prophecy claim", async function () {
      // Send balance that can be unlocked first
      const seedEtherBalance = ethers.utils.parseEther("1");
      await state.bridgeBank.connect(operator)
        .lock(
          state.sender, 
          state.constants.zeroAddress, 
          seedEtherBalance, 
          {value: seedEtherBalance}
        )

      // assert recipient balance before receiving proceeds from newProphecyClaim is correct
      const startingBalance = await getBalance(state.recipient.address, true);
      expect(startingBalance).to.be.equal("10000");
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.constants.zeroAddress,
        state.amount,
        "Ether",
        "ETH",
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.ether,
        [userOne, userTwo, userFour],
      );

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      const endingBalance = await getBalance(state.recipient.address, true);
      const expectedEndingBalance = "10000.0000000000000001"; // added 100 weis
      expect(endingBalance).to.equal(expectedEndingBalance);
    });

    it("should NOT unlock eth upon the successful processing of a burn prophecy claim if the recipient is blocklisted", async function () {
      // Add recipient to the blocklist
      await expect(state.blocklist.addToBlocklist(state.recipient.address)).to.not.be.reverted;

      // assert recipient balance before receiving proceeds from newProphecyClaim is correct
      const startingBalance = await getBalance(state.recipient.address, true);
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.constants.zeroAddress,
        state.amount,
        "Ether",
        "ETH",
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.ether,
        [userOne, userTwo, userFour],
      );

      // get the prophecyID:
      const prophecyID = await state.cosmosBridge.getProphecyID(
        state.sender, // cosmosSender
        state.senderSequence, // cosmosSenderSequence
        state.recipient.address, // ethereumReceiver
        state.constants.zeroAddress, // tokenAddress
        state.amount, // amount
        "Ether", // tokenName
        "ETH", // tokenSymbol
        state.decimals, // tokenDecimals
        state.networkDescriptor, // networkDescriptor
        false, // bridge token
        state.nonce, // nonce
        state.constants.denom.ether, // cosmos denom
      );

      // Doesn't revert, but user doesn't get the funds either
      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      )
        .to.emit(state.cosmosBridge, "LogProphecyCompleted")
        .withArgs(prophecyID, false); // the second argument here is 'success'

      // Make sure the balance didn't change
      const endingBalance = await getBalance(state.recipient.address, true); 
      expect(endingBalance).to.be.equal(startingBalance);
    });

    it("should deploy a new token upon the successful processing of a double-pegged burn prophecy claim", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce,
        state.constants.denom.one,
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
          state.token.address,
          expectedAddress,
          state.constants.denom.one
        );

      // check if the new token has been correctly deployed
      const newlyCreatedTokenAddress = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.constants.denom.one
      );

      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

      // check if the user received minted tokens
      const deployedToken = await state.factories.BridgeToken.attach(newlyCreatedTokenAddress);
      const mintedTokensToRecipient = await deployedToken.balanceOf(state.recipient.address);
      const totalSupply = await deployedToken.totalSupply();
      expect(totalSupply).to.be.equal(state.amount);
      expect(mintedTokensToRecipient).to.be.equal(state.amount);
    });

    it("should deploy a new token upon the successful processing of a double-pegged burn prophecy claim if user is blocklisted, but should not mint", async function () {
      // Add recipient to the blocklist
      await expect(state.blocklist.addToBlocklist(state.recipient.address)).to.be.not.be.reverted;

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce,
        state.constants.denom.one,
        [userOne, userTwo, userFour],
      );

      // get the prophecyID:
      const prophecyID = await state.cosmosBridge.getProphecyID(
        state.sender, // cosmosSender
        state.senderSequence, // cosmosSenderSequence
        state.recipient.address, // ethereumReceiver
        state.token.address, // tokenAddress
        state.amount, // amount
        state.name, // tokenName
        state.symbol, // tokenSymbol
        state.decimals, // tokenDecimals
        state.networkDescriptor, // networkDescriptor
        true, // bridge token
        state.nonce, // nonce
        state.constants.denom.one, // cosmos denom
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
        .to.emit(state.cosmosBridge, "LogProphecyCompleted")
        .withArgs(prophecyID, false); // the second argument here is 'success'

      const newlyCreatedTokenAddress = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.constants.denom.one
      );
      expect(newlyCreatedTokenAddress).to.be.equal(expectedAddress);

      // check if tokens were minted (should not have)
      const deployedToken = await state.factories.BridgeToken.attach(newlyCreatedTokenAddress);
      const mintedTokensToRecipient = await deployedToken.balanceOf(state.recipient.address);
      const totalSupply = await deployedToken.totalSupply();
      expect(totalSupply).to.be.equal(0);
      expect(mintedTokensToRecipient).to.be.equal(0);
    });

    it("should NOT deploy a new token upon the successful processing of a normal burn prophecy claim", async function () {
      // Bridgebank should have an initial balance of token
      await state.token.connect(operator).mint(state.bridgeBank.address, state.amount);

      state.nonce = 1;

      const beforeUserBalance = Number(await state.token.balanceOf(state.recipient.address));
      expect(beforeUserBalance).to.equal(Number(0));

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.one,
        [userOne, userTwo, userFour],
      );

      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      const newlyCreatedTokenAddress = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.token.address
      );
      expect(newlyCreatedTokenAddress).to.be.equal(state.constants.zeroAddress);

      // assert that the recipient's balance of the token went up by the amount we specified in the claim
      const balance = Number(await state.token.balanceOf(state.recipient.address));
      expect(balance).to.equal(state.amount);
    });

    it("should NOT deploy a new token upon the successful processing of a double-pegged burn prophecy claim for an already managed token", async function () {
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce,
        state.constants.denom.one,
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
          state.token.address,
          expectedAddress,
          state.constants.denom.one
        );

      const newlyCreatedTokenAddress = await state.cosmosBridge.cosmosDenomToDestinationAddress(
        state.constants.denom.one
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
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        true,
        state.nonce + 1,
        state.constants.denom.one,
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

    it("should get an signer out of order exception", async function () {
      const beforeUserBalance = Number(await state.token.balanceOf(state.recipient.address));
      expect(beforeUserBalance).to.equal(Number(0));

      // Last nonce should be 0
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(0);

      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.token.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.one,
        [userOne, userTwo, userFour],
      );

      const outOfOrderSignatures: SignedData[] = []
      outOfOrderSignatures.push(signatures[2])
      outOfOrderSignatures.push(signatures[1])
      outOfOrderSignatures.push(signatures[0])

      await expect(
        state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, outOfOrderSignatures)
      ).to.be.revertedWith("custom error 'OutOfOrderSigner(0)'"); 
    });
  });
});
