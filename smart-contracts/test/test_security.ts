import { expect } from "chai";
import {
  signHash,
  setup,
  deployTrollToken,
  getDigestNewProphecyClaim,
  getValidClaim,
  TestFixtureState,
} from "./helpers/testFixture";
import { ethers } from "hardhat";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BridgeToken, BridgeToken__factory, CosmosBridge, CosmosBridge__factory, TrollToken, TrollToken__factory } from "../build";
import { Signer } from "ethers";

interface TestSecurityState extends TestFixtureState {
  troll: TrollToken;
}

// import web3 from "web3"
// const sifRecipient = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace");
const sifRecipient = ethers.utils.hexlify(ethers.utils.toUtf8Bytes("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"));

describe("Security Test", function () {
  let userOne: SignerWithAddress;
  let userTwo: SignerWithAddress;
  let userThree: SignerWithAddress;
  let userFour: SignerWithAddress;
  let accounts: SignerWithAddress[];
  let signerAccounts: string[];
  let operator: SignerWithAddress;
  let owner: SignerWithAddress;
  let pauser: SignerWithAddress;
  const consensusThreshold = 70;
  let initialPowers: number[];
  let initialValidators: string[];
  let networkDescriptor: number;
  // track the state of the deployed contracts
  let state: TestSecurityState;
  let CosmosBridgeFactory: CosmosBridge__factory;
  let BridgeTokenFactory: BridgeToken__factory;
  let TrollToken: TrollToken;

  before(async function () {
    CosmosBridgeFactory = await ethers.getContractFactory("CosmosBridge");
    BridgeTokenFactory = await ethers.getContractFactory("BridgeToken");
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
    initialValidators = [
      accounts[0].address,
      accounts[1].address,
      accounts[2].address,
      accounts[3].address,
    ];

    networkDescriptor = 1;
  });

  describe("BridgeBank Security", function () {
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
      ) as TestSecurityState;
    });

    it("should allow operator to call reinitialize after initialization, setting the correct values", async function () {
      // Change all values to test if the state actually changed
      await expect(
        state.bridgeBank.connect(operator).reinitialize(
          accounts[3].address, // was operator.address
          accounts[4].address, // was state.cosmosBridge.address
          accounts[5].address, // was owner.address
          accounts[6].address, // was pauser.address
          state.networkDescriptor + 1, // was state.networkDescriptor,
          state.rowan.address
        )
      ).to.not.be.reverted;

      expect(await state.bridgeBank.operator()).to.equal(accounts[3].address);
      expect(await state.bridgeBank.cosmosBridge()).to.equal(accounts[4].address);
      expect(await state.bridgeBank.owner()).to.equal(accounts[5].address);
      expect(await state.bridgeBank.pausers(accounts[6].address)).to.equal(true);
      expect(await state.bridgeBank.networkDescriptor()).to.equal(state.networkDescriptor + 1);

      // Expect to keep the previous pauser too
      expect(await state.bridgeBank.pausers(pauser.address)).to.equal(true);
    });

    it("should not allow operator to call reinitialize a second time", async function () {
      await expect(
        state.bridgeBank
          .connect(operator)
          .reinitialize(
            operator.address,
            state.cosmosBridge.address,
            owner.address,
            pauser.address,
            state.networkDescriptor,
            state.rowan.address
          )
      ).to.not.be.reverted;

      await expect(
        state.bridgeBank
          .connect(operator)
          .reinitialize(
            operator.address,
            state.cosmosBridge.address,
            owner.address,
            pauser.address,
            state.networkDescriptor,
            state.rowan.address
          )
      ).to.be.revertedWith("Already reinitialized");
    });

    it("should not allow user to call reinitialize", async function () {
      await expect(
        state.bridgeBank
          .connect(userOne)
          .reinitialize(
            operator.address,
            state.cosmosBridge.address,
            owner.address,
            pauser.address,
            state.networkDescriptor,
            state.rowan.address
          )
      ).to.be.revertedWith("!operator");
    });

    it("should be able to change the owner", async function () {
      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);
      await state.bridgeBank.connect(owner).changeOwner(userTwo.address);
      expect(await state.bridgeBank.owner()).to.be.equal(userTwo.address);
    });

    it("should not be able to change the owner if the caller is not the owner", async function () {
      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);

      await expect(
        state.bridgeBank.connect(accounts[7]).changeOwner(userTwo.address)
      ).to.be.revertedWith("!owner");

      expect(await state.bridgeBank.owner()).to.be.equal(owner.address);
    });

    it("should be able to change BridgeBank's operator", async function () {
      expect(await state.bridgeBank.operator()).to.be.equal(operator.address);
      await expect(state.bridgeBank.connect(operator).changeOperator(userTwo.address))
      .not.to.be.reverted;

      expect(await state.bridgeBank.operator()).to.be.equal(userTwo.address);
    });

    it("should not be able to change BridgeBank's operator if the caller is not the operator", async function () {
      expect(await state.bridgeBank.operator()).to.be.equal(operator.address);
      await expect(
        state.bridgeBank.connect(userOne).changeOperator(userTwo.address)
      ).to.be.revertedWith("!operator");

      expect(await state.bridgeBank.operator()).to.be.equal(operator.address);
    });

    it("should not be able to change the operator if the caller is not the operator", async function () {
      expect(await state.cosmosBridge.operator()).to.be.equal(operator.address);
      await expect(
        state.cosmosBridge.connect(userOne).changeOperator(userTwo.address)
      ).to.be.revertedWith("Must be the operator.");

      expect(await state.cosmosBridge.operator()).to.be.equal(operator.address);
    });

    it("should be able to pause the contract", async function () {
      await state.bridgeBank.connect(pauser).pause();
      expect(await state.bridgeBank.paused()).to.be.true;
    });

    it("should not be able to pause the contract if you are not the owner", async function () {
      await expect(state.bridgeBank.connect(userOne).pause()).to.be.revertedWith(
        "PauserRole: caller does not have the Pauser role"
      );

      expect(await state.bridgeBank.paused()).to.be.false;
    });

    it("should be able to add a new pauser if you are a pauser", async function () {
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.false;

      await state.bridgeBank.connect(pauser).addPauser(userOne.address);

      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.true;
    });

    it("should be able to renounce yourself as pauser", async function () {
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.false;

      await state.bridgeBank.connect(pauser).addPauser(userOne.address);
      expect(await state.bridgeBank.pausers(pauser.address)).to.be.true;
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.true;

      await state.bridgeBank.connect(userOne).renouncePauser();
      expect(await state.bridgeBank.pausers(userOne.address)).to.be.false;
    });

    it("should be able to pause and then unpause the contract", async function () {
      // CosmosBank initial values
      await expect(state.bridgeBank.connect(pauser).unpause()).to.be.revertedWith(
        "Pausable: not paused"
      );

      await state.bridgeBank.connect(pauser).pause();
      await expect(state.bridgeBank.connect(pauser).pause()).to.be.revertedWith("Pausable: paused");

      expect(await state.bridgeBank.paused()).to.be.true;
      await state.bridgeBank.connect(pauser).unpause();

      expect(await state.bridgeBank.paused()).to.be.false;
    });

    it("should not be able to lock when contract is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      expect(await state.bridgeBank.paused()).to.be.true;

      await expect(
        state.bridgeBank.connect(userOne).lock(sifRecipient, state.constants.zeroAddress, 100)
      ).to.be.revertedWith("Pausable: paused");
    });

    it("should not be able to burn when contract is paused", async function () {
      await state.bridgeBank.connect(pauser).pause();
      expect(await state.bridgeBank.paused()).to.be.true;

      await expect(
        state.bridgeBank.connect(userOne).burn(sifRecipient, state.rowan.address, 100)
      ).to.be.revertedWith("Pausable: paused");
    });
  });

  // state entire scenario is mimicking the mainnet scenario where there will be
  // cosmos assets on sifchain, and then we hook into an existing ERC20 contract on mainnet
  // that is eRowan. Then we will try to transfer rowan to eRowan to ensure that
  // everything is set up correctly.
  // We will do state by making a new prophecy claim, validating it with the validators
  // Then ensure that the prohpecy claim paid out the person that it was supposed to
  describe("Bridge token burning", function () {
    before(async function () {
      // state test needs to create a new token contract that will
      // effectively be able to be treated as if it was a cosmos native asset
      // even though it was created on top of ethereum

      // Deploy Valset contract
      state = await setup(
        [userOne.address, userTwo.address, userThree.address],
        [33, 33, 33],
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser,
        networkDescriptor,
      ) as TestSecurityState;
    });

    it("should not allow burning of non whitelisted token address", async function () {
      function convertToHex(str: string) {
        let hex = "";
        for (let i = 0; i < str.length; i++) {
          hex += "" + str.charCodeAt(i).toString(16);
        }
        return hex;
      }

      const amount = 100000;
      const sifAddress = "0x" + convertToHex("sif12qfvgsq76eghlagyfcfyt9md2s9nunsn40zu2h");

      // create new fake eRowan token
      const bridgeToken = await BridgeTokenFactory.deploy(
        "rowan",
        "rowan",
        18,
        state.constants.denom.rowan
      );

      // Attempt to burn tokens
      await expect(
        state.bridgeBank.connect(operator).burn(sifAddress, bridgeToken.address, amount)
      ).to.be.revertedWith("Token is not in Cosmos whitelist");
    });
  });

  describe("Consensus Threshold Limits", function () {
    beforeEach(async function () {
      state = await setup(
        [userOne.address, userTwo.address, userThree.address],
        [33, 33, 33],
        operator,
        consensusThreshold,
        owner,
        userOne,
        userThree,
        pauser,
        networkDescriptor,
      ) as TestSecurityState;
    });

    it("should not allow initialization of CosmosBridge with a consensus threshold over 100", async function () {
      const bridge = await CosmosBridgeFactory.deploy();

      await expect(
        bridge
          .connect(operator)
          .initialize(
            operator.address,
            101,
            state.initialValidators,
            state.initialPowers,
            state.networkDescriptor
          )
      ).to.be.revertedWith("Invalid consensus threshold.");
    });

    it("should not allow initialization of oracle with a consensus threshold of 0", async function () {
      const bridge = await CosmosBridgeFactory.deploy();
      await expect(
        bridge
          .connect(operator)
          .initialize(
            operator.address,
            0,
            state.initialValidators,
            state.initialPowers,
            state.networkDescriptor
          )
      ).to.be.revertedWith("Consensus threshold must be positive.");
    });

    it("should not allow a non cosmosbridge account to mint from bridgebank", async function () {
      await expect(
        state.bridgeBank.connect(operator).handleUnpeg(operator.address, state.token1.address, 100)
      ).to.be.revertedWith("!cosmosbridge");
    });
  });

  describe("Network Descriptor Mismatch", function () {
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
      ) as TestSecurityState;
    });

    it("should not allow unlocking tokens upon the processing of a burn prophecy claim with the wrong network descriptor", async function () {
      state.nonce = 1;

      // this is a valid claim in itself (digest, claimData, and signatures all match)
      // but since we set networkDescriptorMismatch=true in beforeEach(), it will be rejected
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

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_NET_DESC");
    });

    it("should not allow unlocking native tokens upon the processing of a burn prophecy claim with the wrong network descriptor", async function () {
      state.nonce = 1;

      // this is a valid claim in itself (digest, claimData, and signatures all match)
      // but since we set networkDescriptorMismatch=true in beforeEach(), it will be rejected
      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.constants.zeroAddress,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.ether,
        [userOne, userTwo, userFour],
      );

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_NET_DESC");
    });
  });

  describe("Troll token tests", function () {
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
      ) as TestSecurityState;

      TrollToken = await deployTrollToken();
      state.troll = TrollToken;
      await state.troll.mint(userOne.address, 100);
    });

    it("should revert when prophecyclaim is submitted out of order", async function () {
      state.nonce = 10;

      // this is a valid claim in itself (digest, claimData, and signatures all match)
      // but it has the wrong nonce (should be 1, not 10)
      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.troll.address,
        state.amount,
        state.name,
        state.symbol,
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.none,
        [userOne, userTwo, userFour],
      );

      await expect(
        state.cosmosBridge
          .connect(userOne)
          .submitProphecyClaimAggregatedSigs(digest, claimData, signatures)
      ).to.be.revertedWith("INV_ORD");
    });

    it("should allow users to unpeg troll token, but then does not receive", async function () {
      // approve and lock tokens
      await state.troll.connect(userOne).approve(state.bridgeBank.address, state.amount);

      // Attempt to lock tokens
      await state.bridgeBank.connect(userOne).lock(state.sender, state.troll.address, state.amount);

      let endingBalance = Number(await state.troll.balanceOf(userOne.address));
      expect(endingBalance).to.be.equal(0);

      state.recipient = userOne;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        state.troll.address,
        state.amount,
        "Troll",
        "TRL",
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.none,
        [userOne, userTwo, userFour],
      );

      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      // user should not receive funds as troll token just burns gas
      endingBalance = Number(await state.troll.balanceOf(userOne.address));
      expect(endingBalance).to.be.equal(0);

      // Last nonce should now be 1
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);
    });

    it("Should not revert on a reentrancy attack, but user should not receive funds either", async function () {
      // Deploy reentrancy attacker token:
      const reentrancyTokenFactory = await ethers.getContractFactory("ReentrancyToken");
      const reentrancyToken = await reentrancyTokenFactory.deploy(
        "Troll Token",
        "TROLL",
        state.cosmosBridge.address,
        userOne.address,
        state.amount
      );
      await reentrancyToken.deployed();

      // approve and lock tokens
      await reentrancyToken.connect(userOne).approve(state.bridgeBank.address, state.amount);

      // Attempt to lock tokens
      await state.bridgeBank
        .connect(userOne)
        .lock(state.sender, reentrancyToken.address, state.amount);

      let endingBalance = Number(await reentrancyToken.balanceOf(userOne.address));
      expect(endingBalance).to.be.equal(0);

      state.recipient = userOne;
      state.nonce = 1;

      const { digest, claimData, signatures } = await getValidClaim(
        state.sender,
        state.senderSequence,
        state.recipient.address,
        reentrancyToken.address,
        state.amount,
        "Troll",
        "TRL",
        state.decimals,
        state.networkDescriptor,
        false,
        state.nonce,
        state.constants.denom.none,
        [userOne, userTwo, userFour],
      );

      // Reentrancy token will try to reenter submitProphecyClaimAggregatedSigs
      await state.cosmosBridge
        .connect(userOne)
        .submitProphecyClaimAggregatedSigs(digest, claimData, signatures);

      // user should not receive funds as a Reentrancy should break the transfer flow
      endingBalance = Number(await reentrancyToken.balanceOf(userOne.address));
      expect(endingBalance).to.be.equal(0);

      // Last nonce should now be 1 (because it does NOT REVERT)
      let lastNonceSubmitted = Number(await state.cosmosBridge.lastNonceSubmitted());
      expect(lastNonceSubmitted).to.be.equal(1);
    });
  });
});
