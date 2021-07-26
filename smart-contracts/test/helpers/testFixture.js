const { ethers, upgrades } = require("hardhat");
const web3 = require("web3")

async function returnContractObjects() {
    let CosmosBridge = await ethers.getContractFactory("CosmosBridge");
    let BridgeBank = await ethers.getContractFactory("BridgeBank");
    let BridgeToken = await ethers.getContractFactory("BridgeToken");

    return {CosmosBridge, BridgeBank, BridgeToken};
}

function getDigestNewProphecyClaim(data) {
  if (!Array.isArray(data)) {
    throw new Error("Input Error: not array");
  }

  const digest = ethers.utils.keccak256(
    ethers.utils.defaultAbiCoder.encode(
      [
        "bytes",
        "uint256",
        "address",
        "address",
        "uint256",
        "bool",
        "uint128",
        "uint256"
      ],
      data
    ),
  );

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

async function multiTokenSetup(
    initialValidators,
    initialPowers,
    operator,
    consensusThreshold,
    owner,
    userOne,
    userThree,
    pauser,
    chainId,
    chainIdMismatch = false
  ) {
    const state = {}

    // Deploy Valset contract
    state.initialValidators = initialValidators;
    state.initialPowers = initialPowers;

    const { CosmosBridge, BridgeBank, BridgeToken } = await returnContractObjects();

    // Chain descriptor
    state.chainId = chainId;

    // Deploy CosmosBridge contract
    state.cosmosBridge = await upgrades.deployProxy(CosmosBridge, [
      operator.address,
      consensusThreshold,
      initialValidators,
      initialPowers,
      chainId
    ]);
    await state.cosmosBridge.deployed();

    // Deploy BridgeBank contract
    state.bridgeBank = await upgrades.deployProxy(BridgeBank, [
      state.cosmosBridge.address,
      owner.address,
      pauser,
      chainIdMismatch ? state.chainId + 1 : chainId
    ]);
    await state.bridgeBank.deployed();

    // Operator sets Bridge Bank
    await state.cosmosBridge.connect(operator).setBridgeBank(state.bridgeBank.address);

    // state is for ERC20 deposits
    state.sender = web3.utils.utf8ToHex(
      "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
    );
    state.cosmosSender = state.sender;
    state.senderSequence = 1;
    state.recipient = userThree;
    state.name = "TEST COIN";
    state.symbol = "TEST";
    state.ethereumToken = "0x0000000000000000000000000000000000000000";
    state.weiAmount = web3.utils.toWei("0.25", "ether");
    state.amount = 100;

    state.rowan = await BridgeToken.deploy("rowan", "rowan", 18);

    await state.rowan.deployed();
    // mint tokens
    await state.rowan.connect(operator).mint(userOne.address, state.amount * 2);
    // add bridgebank as owner of the rowan contract
    await state.rowan.transferOwnership(state.bridgeBank.address);

    await state.rowan.connect(userOne).approve(state.bridgeBank.address, state.amount * 2);

    // Add rowan as an existing bridge token
    await state.bridgeBank.connect(owner).addExistingBridgeToken(state.rowan.address);

    state.token1 = await BridgeToken.deploy(state.name, state.symbol, 18);
    state.token2 = await BridgeToken.deploy(state.name, state.symbol, 18);
    state.token3 = await BridgeToken.deploy(state.name, state.symbol, 18);

    await state.token1.deployed();
    await state.token2.deployed();
    await state.token3.deployed();

    //Load user account with ERC20 tokens for testing
    await state.token1.connect(operator).mint(userOne.address, state.amount * 2);
    await state.token2.connect(operator).mint(userOne.address, state.amount * 2);
    await state.token3.connect(operator).mint(userOne.address, state.amount * 2);

    await state.token1.connect(userOne).approve(state.bridgeBank.address, state.amount * 2);
    await state.token2.connect(userOne).approve(state.bridgeBank.address, state.amount * 2);
    await state.token3.connect(userOne).approve(state.bridgeBank.address, state.amount * 2);

    return state;
}

async function singleSetup(
    initialValidators,
    initialPowers,
    operator,
    consensusThreshold,
    owner,
    userOne,
    userThree,
    pauser,
    chainId,
    chainIdMismatch = false
    ) {
    const state = {};
    // Deploy Valset contract
    state.initialValidators = initialValidators;
    state.initialPowers = initialPowers;

    const { CosmosBridge, BridgeBank, BridgeToken } = await returnContractObjects();

    // Chain descriptor
    state.chainId = chainId;

    // Deploy CosmosBridge contract
    state.cosmosBridge = await upgrades.deployProxy(CosmosBridge, [
      operator.address,
      consensusThreshold,
      initialValidators,
      initialPowers,
      chainId
    ]);
    await state.cosmosBridge.deployed();

    // Deploy BridgeBank contract
    state.bridgeBank = await upgrades.deployProxy(BridgeBank, [
      state.cosmosBridge.address,
      owner.address,
      pauser,
      chainIdMismatch ? state.chainId + 1 : chainId
    ]);
    await state.bridgeBank.deployed();

    // Operator sets Bridge Bank
    await state.cosmosBridge.connect(operator).setBridgeBank(state.bridgeBank.address);

    // state is for ERC20 deposits
    state.sender = web3.utils.utf8ToHex(
      "sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace"
    );
    state.senderSequence = 1;
    state.recipient = userThree;
    state.name = "TEST COIN";
    state.symbol = "TEST";
    state.ethereumToken = "0x0000000000000000000000000000000000000000";
    state.weiAmount = web3.utils.toWei("0.25", "ether");

    state.token = await BridgeToken.deploy(
      state.name,
      state.symbol,
      18
    );

    state.rowan = await BridgeToken.deploy("rowan", "rowan", 18);

    await state.token.deployed();
    state.amount = 100;
    //Load user account with ERC20 tokens for testing
    await state.token.connect(operator).mint(userOne.address, state.amount * 2);

    // Approve tokens to contract
    await state.token.connect(userOne).approve(state.bridgeBank.address, state.amount).should.be.fulfilled;
      
    // Lock tokens on contract
    await state.bridgeBank.connect(userOne).lock(
      state.sender,
      state.token.address,
      state.amount
    ).should.be.fulfilled;

    // Lock tokens on contract
    await state.bridgeBank.connect(userOne).lock(
      state.sender,
      state.ethereumToken,
      state.amount, {
        value: state.amount
      }
    ).should.be.fulfilled;

    return state;
}

async function deployTrollToken() {
  let TrollToken = await ethers.getContractFactory("TrollToken");
  const troll = await TrollToken.deploy("Troll", "TRL");

  return troll;
}

module.exports = {
    multiTokenSetup,
    singleSetup,
    deployTrollToken,
    signHash,
    getDigestNewProphecyClaim
};
