const { ethers, upgrades } = require("hardhat");
const web3 = require("web3");

const { ROWAN_DENOM, ETHER_DENOM, DENOM_1, DENOM_2, DENOM_3, DENOM_4, IBC_DENOM } = require("./denoms");

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

async function getContractFactories() {
  const CosmosBridge = await ethers.getContractFactory("CosmosBridge");
  const BridgeBank = await ethers.getContractFactory("BridgeBank");
  const BridgeToken = await ethers.getContractFactory("BridgeToken");
  const Blocklist = await ethers.getContractFactory("Blocklist");

  return { CosmosBridge, BridgeBank, BridgeToken, Blocklist };
}

function getDigestNewProphecyClaim(data) {
  if (!Array.isArray(data)) {
    throw new Error("Input Error: not array");
  }

  const types = [
    "bytes", // cosmosSender
    "uint256", // cosmosSenderSequence
    "address", // ethereumReceiver
    "address", // tokenAddress
    "uint256", // amount
    "string", // tokenName
    "string", // tokenSymbol
    "uint8", // tokenDecimals
    "int32", // networkDescriptor
    "bool", // bridgetoken
    "uint256", // nonce
    "string", // cosmosDenom
  ];

  if (types.length !== data.length) {
    throw new Error("testFixture::getDigestNewProphecyClaim: invalid data length");
  }

  const digest = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(types, data));

  return digest;
}

async function signHash(signers, hash) {
  let sigData = [];

  for (let i = 0; i < signers.length; i++) {
    let sig = await signers[i].signMessage(ethers.utils.arrayify(hash));

    const splitSig = ethers.utils.splitSignature(sig);
    sig = {
      signer: signers[i].address,
      _v: splitSig.v,
      _r: splitSig.r,
      _s: splitSig.s,
    };

    sigData.push(sig);
  }

  return sigData;
}

async function setup({
  initialValidators,
  initialPowers,
  operator,
  consensusThreshold,
  owner,
  user,
  recipient,
  pauser,
  networkDescriptor,
  networkDescriptorMismatch = false,
  lockTokensOnBridgeBank = false,
}) {
  const state = initState({
    initialValidators,
    initialPowers,
    operator,
    consensusThreshold,
    owner,
    user,
    recipient,
    pauser,
    networkDescriptor,
    networkDescriptorMismatch,
  });

  // Our upgrades use a delegateCall to the Address lib internally; we'll silence warnings
  upgrades.silenceWarnings();

  await deployBaseContracts(state);
  await deployRowan(state);

  if (lockTokensOnBridgeBank) {
    // Lock tokens on contract
    await state.bridgeBank.connect(user).lock(state.sender, state.token.address, state.amount)
      .should.be.fulfilled;

    // Lock native tokens on contract
    await state.bridgeBank
      .connect(user)
      .lock(state.sender, state.constants.zeroAddress, state.amount, { value: state.amount }).should
      .be.fulfilled;
  }

  return state;
}

function initState({
  initialValidators,
  initialPowers,
  operator,
  consensusThreshold,
  owner,
  user,
  recipient,
  pauser,
  networkDescriptor,
  networkDescriptorMismatch,
}) {
  const sender = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace");
  const state = {
    constants: {
      zeroAddress: ZERO_ADDRESS,
      roles: {
        minter: web3.utils.soliditySha3("MINTER_ROLE"),
        admin: "0x0000000000000000000000000000000000000000000000000000000000000000",
      },
      denom: {
        none: "",
        rowan: ROWAN_DENOM,
        ether: ETHER_DENOM,
        one: DENOM_1,
        two: DENOM_2,
        three: DENOM_3,
        four: DENOM_4,
        ibc: IBC_DENOM,
      },
    },
    initialValidators,
    initialPowers,
    operator,
    consensusThreshold,
    owner,
    user,
    recipient,
    pauser,
    networkDescriptor,
    networkDescriptorMismatch,
    sender,
    cosmosSender: sender,
    senderSequence: 1,
    recipient,
    name: "TEST COIN",
    symbol: "TEST",
    decimals: 18,
    weiAmount: web3.utils.toWei("0.25", "ether"),
    amount: 100,
  };

  return state;
}

async function deployBaseContracts(state) {
  const { CosmosBridge, BridgeBank, BridgeToken, Blocklist } = await getContractFactories();
  state.factories = { CosmosBridge, BridgeBank, BridgeToken, Blocklist };

  // Deploy CosmosBridge contract
  state.cosmosBridge = await upgrades.deployProxy(
    CosmosBridge,
    [
      state.operator.address,
      state.consensusThreshold,
      state.initialValidators,
      state.initialPowers,
      state.networkDescriptorMismatch ? state.networkDescriptor + 1 : state.networkDescriptor,
    ],
    {
      initializer: "initialize(address,uint256,address[],uint256[],int32)",
      unsafeAllow: ["delegatecall"],
    }
  );
  await state.cosmosBridge.deployed();

  // Deploy BridgeBank contract
  state.bridgeBank = await upgrades.deployProxy(
    BridgeBank,
    [
      state.operator.address,
      state.cosmosBridge.address,
      state.owner.address,
      state.pauser.address,
      state.networkDescriptorMismatch ? state.networkDescriptor + 2 : state.networkDescriptor,
    ],
    {
      initializer: "initialize(address,address,address,address,int32)",
      unsafeAllow: ["delegatecall"],
    }
  );
  await state.bridgeBank.deployed();

  // Operator sets Bridge Bank
  await state.cosmosBridge.connect(state.operator).setBridgeBank(state.bridgeBank.address);

  // Deploy BridgeTokens
  state.token = await BridgeToken.deploy(
    state.name,
    state.symbol,
    state.decimals,
    state.constants.denom.one
  );
  state.token1 = await BridgeToken.deploy(
    state.name,
    state.symbol,
    state.decimals,
    state.constants.denom.two
  );
  state.token2 = await BridgeToken.deploy(
    state.name,
    state.symbol,
    state.decimals,
    state.constants.denom.three
  );
  state.token3 = await BridgeToken.deploy(
    state.name,
    state.symbol,
    state.decimals,
    state.constants.denom.four
  );
  state.token_noDenom = await BridgeToken.deploy(
    state.name,
    state.symbol,
    state.decimals,
    state.constants.denom.none
  );
  state.token_ibc = await BridgeToken.deploy(
    state.name,
    state.symbol,
    state.decimals,
    state.constants.denom.ibc
  );

  await state.token.deployed();
  await state.token1.deployed();
  await state.token2.deployed();
  await state.token3.deployed();
  await state.token_noDenom.deployed();
  await state.token_ibc.deployed();

  // Grant the MINTER role to the operator:
  await state.token
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.operator.address);
  await state.token1
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.operator.address);
  await state.token2
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.operator.address);
  await state.token3
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.operator.address);
  await state.token_noDenom
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.operator.address);
  await state.token_ibc
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.operator.address);

  // Load user account with ERC20 tokens for testing
  await state.token.connect(state.operator).mint(state.user.address, state.amount * 2);
  await state.token1.connect(state.operator).mint(state.user.address, state.amount * 2);
  await state.token2.connect(state.operator).mint(state.user.address, state.amount * 2);
  await state.token3.connect(state.operator).mint(state.user.address, state.amount * 2);
  await state.token_noDenom.connect(state.operator).mint(state.user.address, state.amount * 2);
  await state.token_ibc.connect(state.operator).mint(state.user.address, state.amount * 2);

  // Approve BridgeBank
  await state.token.connect(state.user).approve(state.bridgeBank.address, state.amount * 2);
  await state.token1.connect(state.user).approve(state.bridgeBank.address, state.amount * 2);
  await state.token2.connect(state.user).approve(state.bridgeBank.address, state.amount * 2);
  await state.token3.connect(state.user).approve(state.bridgeBank.address, state.amount * 2);
  await state.token_noDenom.connect(state.user).approve(state.bridgeBank.address, state.amount * 2);
  await state.token_ibc.connect(state.user).approve(state.bridgeBank.address, state.amount * 2);

  // Deploy the Blocklist
  state.blocklist = await Blocklist.deploy();
  await state.blocklist.deployed();

  // Register the blocklist on BridgeBank
  await state.bridgeBank.connect(state.operator).setBlocklist(state.blocklist.address);
}

async function deployRowan(state) {
  // deploy
  state.rowan = await state.factories.BridgeToken.deploy(
    "rowan",
    "rowan",
    state.decimals,
    state.constants.denom.rowan
  );
  await state.rowan.deployed();

  // mint tokens
  await state.rowan
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.operator.address);
  await state.rowan.connect(state.operator).mint(state.user.address, state.amount * 2);

  // add bridgebank as admin and minter of the rowan contract
  await state.rowan
    .connect(state.operator)
    .grantRole(state.constants.roles.minter, state.bridgeBank.address);
  await state.rowan
    .connect(state.operator)
    .grantRole(state.constants.roles.admin, state.bridgeBank.address);

  // approve bridgeBank
  await state.rowan.connect(state.user).approve(state.bridgeBank.address, state.amount * 2);

  // add rowan as an existing bridge token
  await state.bridgeBank.connect(state.owner).addExistingBridgeToken(state.rowan.address);
}

async function deployTrollToken() {
  let TrollToken = await ethers.getContractFactory("TrollToken");
  const troll = await TrollToken.deploy("Troll", "TRL");

  return troll;
}

/**
 * A Commission Token (CMT) is a token that charges a dev fee commission for every transaction so if I send you
 * 100 CMT with a devFee of 5% you would get 95 CMT and the dev gets 5 CMT.
 * @param {address} devAccount The account which gets the commissions on transfer
 * @param {uint256} devFee The fee to charge per transaction in ten thousandths of a percent
 * @param {address} userAccount The user to mint tokens to
 * @param {uint256} quantity The quantity of tokens to mint
 * @returns An Ethers CommissionTokenContract
 */
async function deployCommissionToken(devAccount, devFee, userAccount, quantity) {
  const tokenFactory = await ethers.getContractFactory("CommissionToken");
  const token = await tokenFactory.deploy(devAccount, devFee, userAccount, quantity);
  return token;
}

/**
 * Creates a valid claim
 * @returns { digest, signatures, claimData }
 */
async function getValidClaim({
  sender,
  senderSequence,
  recipientAddress,
  tokenAddress,
  amount,
  tokenName,
  tokenSymbol,
  tokenDecimals,
  networkDescriptor,
  bridgeToken,
  nonce,
  cosmosDenom,
  validators,
}) {
  const digest = getDigestNewProphecyClaim([
    sender,
    senderSequence,
    recipientAddress,
    tokenAddress,
    amount,
    tokenName,
    tokenSymbol,
    tokenDecimals,
    networkDescriptor,
    bridgeToken,
    nonce,
    cosmosDenom,
  ]);

  const signatures = await signHash(validators, digest);

  const claimData = {
    cosmosSender: sender,
    cosmosSenderSequence: senderSequence,
    ethereumReceiver: recipientAddress,
    tokenAddress,
    amount,
    tokenName,
    tokenSymbol,
    tokenDecimals,
    networkDescriptor,
    bridgeToken,
    nonce,
    cosmosDenom,
  };

  const result = {
    digest,
    signatures,
    claimData,
  };

  return result;
}

module.exports = {
  setup,
  deployTrollToken,
  deployCommissionToken,
  signHash,
  getDigestNewProphecyClaim,
  getValidClaim,
};
