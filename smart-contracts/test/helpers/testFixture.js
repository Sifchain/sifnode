const { ethers, upgrades } = require("hardhat");
const web3 = require("web3")

async function returnContractObjects() {
    let CosmosBridge = await ethers.getContractFactory("CosmosBridge");
    let BridgeBank = await ethers.getContractFactory("BridgeBank");
    let BridgeToken = await ethers.getContractFactory("BridgeToken");

    return {CosmosBridge, BridgeBank, BridgeToken};
}

async function multiTokenSetup(
    initialValidators,
    initialPowers,
    operator,
    consensusThreshold,
    owner,
    userOne,
    userThree
    ) {
    const state = {}

    // Deploy Valset contract
    state.initialValidators = initialValidators;
    state.initialPowers = initialPowers;

    const { CosmosBridge, BridgeBank, BridgeToken } = await returnContractObjects();

    // Deploy CosmosBridge contract
    state.cosmosBridge = await upgrades.deployProxy(CosmosBridge, [
      operator,
      consensusThreshold,
      initialValidators,
      initialPowers
    ]);
    await state.cosmosBridge.deployed();

    // Deploy BridgeBank contract
    state.bridgeBank = await upgrades.deployProxy(BridgeBank, [
      state.cosmosBridge.address,
      owner,
      pauser
    ]);
    await state.bridgeBank.deployed();

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
    state.amount = 100;

    state.token1 = await BridgeToken.deploy(state.name, state.symbol, 18);

    state.token2 = await BridgeToken.deploy(state.name, state.symbol, 18);

    state.token3 = await BridgeToken.deploy(state.name, state.symbol, 18);

    await state.token1.deployed();
    await state.token2.deployed();
    await state.token3.deployed();

    //Load user account with ERC20 tokens for testing
    await state.token1.mint(userOne.address, state.amount * 2, {
      from: operator
    }).should.be.fulfilled;

    await state.token2.mint(userOne.address, state.amount * 2, {
      from: operator
    }).should.be.fulfilled;

    await state.token3.mint(userOne.address, state.amount * 2, {
      from: operator
    }).should.be.fulfilled;

    return state;
}

async function singleSetup(
    initialValidators,
    initialPowers,
    operator,
    consensusThreshold,
    owner,
    userOne,
    userThree
    ) {
    const state = {};
    // Deploy Valset contract
    state.initialValidators = initialValidators;
    state.initialPowers = initialPowers;

    const { CosmosBridge, BridgeBank, BridgeToken } = await returnContractObjects();

    // Deploy CosmosBridge contract
    state.cosmosBridge = await upgrades.deployProxy(CosmosBridge, [
      operator,
      consensusThreshold,
      initialValidators,
      initialPowers
    ]);
    await state.cosmosBridge.deployed();

    // Deploy BridgeBank contract
    state.bridgeBank = await upgrades.deployProxy(BridgeBank, [
      state.cosmosBridge.address,
      owner,
      pauser
    ]);
    await state.bridgeBank.deployed();

    // Operator sets Bridge Bank
    await state.cosmosBridge.setBridgeBank(state.bridgeBank.address, {
      from: operator
    });

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

    await state.token.deployed();
    state.amount = 100;
    //Load user account with ERC20 tokens for testing
    await state.token.mint(userOne.address, state.amount * 2, {
      from: operator
    }).should.be.fulfilled;

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

module.exports = {
    multiTokenSetup,
    singleSetup
};
