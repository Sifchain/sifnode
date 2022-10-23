import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { BigNumber, BigNumberish, BytesLike, Contract, ContractTransaction } from "ethers";
import { ethers, upgrades } from "hardhat";
import web3 from "web3";
import { Blocklist, Blocklist__factory, BridgeBank, BridgeBank__factory, BridgeToken, BridgeToken__factory, CosmosBridge, CosmosBridge__factory, Erowan, Erowan__factory } from "../../build";

import { ROWAN_DENOM, ETHER_DENOM, DENOM_1, DENOM_2, DENOM_3, DENOM_4, IBC_DENOM } from "./denoms";

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

async function getContractFactories() {
  const CosmosBridge = await ethers.getContractFactory("CosmosBridge");
  const BridgeBank = await ethers.getContractFactory("BridgeBank");
  const BridgeToken = await ethers.getContractFactory("BridgeToken");
  const Blocklist = await ethers.getContractFactory("Blocklist");
  const Rowan = await ethers.getContractFactory("Erowan");

  return { CosmosBridge, BridgeBank, BridgeToken, Blocklist, Rowan };
}

function getDigestNewProphecyClaim(data: unknown[]) {

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

export interface SignedData {
  signer: string;
  _v: number;
  _r: string;
  _s: string;
}

interface TestFixtureStateConstants {
  zeroAddress: typeof ZERO_ADDRESS;
  roles: {
    minter: string;
    admin: "0x0000000000000000000000000000000000000000000000000000000000000000";
  };
  denom: {
    none: "";
    rowan: string;
    ether: string;
    one: string;
    two: string;
    three: string;
    four: string;
    ibc: string;
  };
}

interface TestFixtureAccounts {
  initialValidators: string[];
  initialPowers: number[];
  operator: SignerWithAddress;
  consensusThreshold: number;
  owner: SignerWithAddress;
  user: SignerWithAddress;
  recipient: SignerWithAddress;
  pauser: SignerWithAddress;
  sender: string;
  cosmosSender: string;
  senderSequence: number;
}

interface TestFixtureNetworksAndTokens {
  name: string;
  networkDescriptor: number;
  networkDescriptorMismatch: boolean;
  symbol: string;
  decimals: number;
  weiAmount: string;
  amount: number;
}

interface TestFixtureContracts {
  bridgeBank: BridgeBank;
  cosmosBridge: CosmosBridge;
  blocklist: Blocklist;
  factories: {
    CosmosBridge: CosmosBridge__factory;
    BridgeBank: BridgeBank__factory;
    BridgeToken: BridgeToken__factory;
    Blocklist: Blocklist__factory;
    Rowan: Erowan__factory;
  }
}

interface TestFixtureTokens {
  token: BridgeToken;
  token1: BridgeToken;
  token2: BridgeToken;
  token3: BridgeToken;
  token_ibc: BridgeToken;
  token_noDenom: BridgeToken;
}

export interface TestFixtureState extends
  TestFixtureAccounts,
  TestFixtureNetworksAndTokens,
  TestFixtureContracts,
  TestFixtureTokens {
  constants: TestFixtureStateConstants;
  rowan: Erowan;
  nonce?: number;
}
async function signHash(signers: SignerWithAddress[], hash: BytesLike): Promise<SignedData[]> {
  let sigData: SignedData[] = [];

  for (let i = 0; i < signers.length; i++) {
    const sig = await signers[i].signMessage(ethers.utils.arrayify(hash));

    const splitSig = ethers.utils.splitSignature(sig);
    const signedMessage = {
      signer: signers[i].address,
      _v: splitSig.v,
      _r: splitSig.r,
      _s: splitSig.s,
    };

    sigData.push(signedMessage);
  }

  return sigData;
}

async function setup(
  initialValidators: string[],
  initialPowers: number[],
  operator: SignerWithAddress,
  consensusThreshold: number,
  owner: SignerWithAddress,
  user: SignerWithAddress,
  recipient: SignerWithAddress,
  pauser: SignerWithAddress,
  networkDescriptor: number,
  networkDescriptorMismatch = false,
  lockTokensOnBridgeBank = false,
) {
  const state = await initState(
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
  );

  if (lockTokensOnBridgeBank) {
    // Lock tokens on contract
    await state.bridgeBank.connect(user).lock(state.sender, state.token.address, state.amount)
      .should.not.be.reverted;

    // Lock native tokens on contract
    await state.bridgeBank
      .connect(user)
      .lock(state.sender, state.constants.zeroAddress, state.amount, { value: state.amount }).should
      .not.be.reverted;
  }

  return state;
}

async function initState(
  initialValidators: string[],
  initialPowers: number[],
  operator: SignerWithAddress,
  consensusThreshold: number,
  owner: SignerWithAddress,
  user: SignerWithAddress,
  recipient: SignerWithAddress,
  pauser: SignerWithAddress,
  networkDescriptor: number,
  networkDescriptorMismatch: boolean,
): Promise<TestFixtureState> {
  const sender = web3.utils.utf8ToHex("sif1nx650s8q9w28f2g3t9ztxyg48ugldptuwzpace");
  const testConstants: TestFixtureStateConstants = {
    zeroAddress: ZERO_ADDRESS,
    roles: {
      minter: web3.utils.soliditySha3("MINTER_ROLE") as string,
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
  };

  const testAccounts: TestFixtureAccounts = {
    initialValidators,
    initialPowers,
    operator,
    consensusThreshold,
    owner,
    user,
    recipient,
    pauser,
    sender,
    cosmosSender: sender,
    senderSequence: 1,
  };

  const testNetworkandTokens: TestFixtureNetworksAndTokens = {
    name: "TEST COIN",
    symbol: "TEST",
    networkDescriptor,
    networkDescriptorMismatch,
    decimals: 18,
    weiAmount: web3.utils.toWei("0.25", "ether"),
    amount: 100
  };

// Our upgrades use a delegateCall to the Address lib internally; we'll silence warnings
upgrades.silenceWarnings();

// Add base contracts to state object
const { contracts, tokens } = await deployBaseContracts(testAccounts, testConstants, testNetworkandTokens);
const rowan = await deployRowan(contracts, testNetworkandTokens, testAccounts, testConstants);

return {
  ...contracts,
  ...tokens,
  ...testNetworkandTokens,
  ...testAccounts,
  constants: testConstants,
  rowan,
};
}

interface deployBaseContractsReturn {
  contracts: TestFixtureContracts;
  tokens: TestFixtureTokens;
}

async function deployBaseContracts(accounts: TestFixtureAccounts, constants: TestFixtureStateConstants, tokenInfo: TestFixtureNetworksAndTokens): Promise<deployBaseContractsReturn> {
  const { CosmosBridge, BridgeBank, BridgeToken, Blocklist, Rowan } = await getContractFactories();
  const factories = { CosmosBridge, BridgeBank, BridgeToken, Blocklist, Rowan };

  // Deploy CosmosBridge contract
  const cosmosBridge = await upgrades.deployProxy(
    CosmosBridge,
    [
      accounts.operator.address,
      accounts.consensusThreshold,
      accounts.initialValidators,
      accounts.initialPowers,
      tokenInfo.networkDescriptorMismatch ? tokenInfo.networkDescriptor + 1 : tokenInfo.networkDescriptor,
    ],
    {
      initializer: "initialize(address,uint256,address[],uint256[],int32)",
      unsafeAllow: ["delegatecall"],
    }
  ) as CosmosBridge;
  await cosmosBridge.deployed();

  // Deploy BridgeBank contract
  const bridgeBank = await upgrades.deployProxy(
    BridgeBank,
    [
      accounts.operator.address,
      cosmosBridge.address,
      accounts.owner.address,
      accounts.pauser.address,
      tokenInfo.networkDescriptorMismatch ? tokenInfo.networkDescriptor + 2 : tokenInfo.networkDescriptor,
      constants.zeroAddress
    ],
    {
      initializer: "initialize(address,address,address,address,int32,address)",
      unsafeAllow: ["delegatecall"],
    }
  ) as BridgeBank;
  await bridgeBank.deployed();

  // Operator sets Bridge Bank
  await cosmosBridge.connect(accounts.operator).setBridgeBank(bridgeBank.address);

  // Deploy BridgeTokens
  const token = await BridgeToken.connect(accounts.operator).deploy(
    tokenInfo.name,
    tokenInfo.symbol,
    tokenInfo.decimals,
    constants.denom.one
  );
  const token1 = await BridgeToken.connect(accounts.operator).deploy(
    tokenInfo.name,
    tokenInfo.symbol,
    tokenInfo.decimals,
    constants.denom.two
  );
  const token2 = await BridgeToken.connect(accounts.operator).deploy(
    tokenInfo.name,
    tokenInfo.symbol,
    tokenInfo.decimals,
    constants.denom.three
  );
  const token3 = await BridgeToken.connect(accounts.operator).deploy(
    tokenInfo.name,
    tokenInfo.symbol,
    tokenInfo.decimals,
    constants.denom.four
  );
  const token_noDenom = await BridgeToken.connect(accounts.operator).deploy(
    tokenInfo.name,
    tokenInfo.symbol,
    tokenInfo.decimals,
    constants.denom.none
  );
  const token_ibc = await BridgeToken.connect(accounts.operator).deploy(
    tokenInfo.name,
    tokenInfo.symbol,
    tokenInfo.decimals,
    constants.denom.ibc
  );

  await token.deployed();
  await token1.deployed();
  await token2.deployed();
  await token3.deployed();
  await token_noDenom.deployed();
  await token_ibc.deployed();

  // TODO: Only have IBC and noDenom tokens in this array update the 
  //       unit tests to only use ibc and rowan tokens as bridgetokens
  const tokens = [token_noDenom, token_ibc, token, token1, token2, token3];

  for (const currentToken of tokens) {
    // Grant the MINTER and ADMIN roles to the Bridgebank
    await currentToken.connect(accounts.operator)
      .grantRole(constants.roles.minter, bridgeBank.address)
    await currentToken.connect(accounts.operator)
      .grantRole(constants.roles.admin, bridgeBank.address)
  }

  // Deploy the Blocklist
  const blocklist = await Blocklist.deploy();
  await blocklist.deployed();

  // Register the blocklist on BridgeBank
  await bridgeBank.connect(accounts.operator).setBlocklist(blocklist.address);

  return {
    contracts: {
      bridgeBank,
      cosmosBridge,
      blocklist,
      factories
    },
    tokens: {
      token,
      token1,
      token2,
      token3,
      token_ibc,
      token_noDenom
    }
  }
}

async function deployRowan(contracts: TestFixtureContracts, tokenInfo: TestFixtureNetworksAndTokens, accounts: TestFixtureAccounts, constants: TestFixtureStateConstants) {
  // deploy
  const rowan = await contracts.factories.Rowan.deploy(
    "Erowan",
  );
  await rowan.deployed();

  // mint tokens
  await rowan.connect(accounts.operator).mint(accounts.user.address, tokenInfo.amount * 2);

  // add bridgebank as admin and minter of the rowan contract
  await rowan
    .connect(accounts.operator)
    .addMinter(contracts.bridgeBank.address);
  
  // approve bridgeBank
  await rowan.connect(accounts.user).approve(contracts.bridgeBank.address, tokenInfo.amount * 2);

  // add rowan as an existing bridge token
  await contracts.bridgeBank.connect(accounts.owner).addExistingBridgeToken(rowan.address);

  // Set Rowan as the Rowan special account
  await contracts.bridgeBank.connect(accounts.operator).setRowanTokenAddress(rowan.address);
  
  return rowan;
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
async function deployCommissionToken(devAccount: string, devFee: BigNumberish, userAccount: string, quantity: BigNumberish) {
  const tokenFactory = await ethers.getContractFactory("CommissionToken");
  const token = await tokenFactory.deploy(devAccount, devFee, userAccount, quantity);
  return token;
}

/**
 * Creates a valid claim
 * @returns { digest, signatures, claimData }
 */
async function getValidClaim(
  sender: string,
  senderSequence: number,
  recipientAddress: string,
  tokenAddress: string,
  amount: number,
  tokenName: string,
  tokenSymbol: string,
  tokenDecimals: number,
  networkDescriptor: number,
  bridgeToken: boolean,
  nonce: number,
  cosmosDenom: string,
  validators: SignerWithAddress[],
) {
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

  const sorted_validators = validators.sort((v1, v2) => {
    if (v1.address > v2.address) {
      return 1;
    }
    if (v1.address < v2.address) {
      return -1;
    }
    return 0;
  })

  const signatures = await signHash(sorted_validators, digest);

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

/**
 * This utility function will prefund the given account with amount quantity requested 
 * in tokens by minting tokens. The operator that deployed the tokens or has mint privilege
 * must be provided.
 * @param user The SignerWithAddress that is to be funded
 * @param amount The quantity to fund the account by
 * @param operator The account with mint privlages on all tokens listed
 * @param tokens An array of tokens to mint on
 */
async function prefundAccount(user: SignerWithAddress | Contract, amount: BigNumberish, operator: SignerWithAddress, tokens: BridgeToken[]) {
   let tokenPromises: Promise<ContractTransaction>[] = [];
    for (const token of tokens) {
      tokenPromises.push(token.connect(operator).mint(user.address, amount));
    }
    await Promise.all(tokenPromises);
}

/**
 * This utility function will preapprove the given user from the approvers balance on all tokens listed.
 * @param user The account to be approved to spend approvers funds
 * @param approver The account which is approving the the user passed to be funded from approvers funds
 * @param amount The amount of tokens that the user is approved for
 * @param tokens An array of tokens to approve the account on
 */
async function preApproveAccount(user: SignerWithAddress | Contract, approver: SignerWithAddress, amount: BigNumberish, tokens: BridgeToken[]) {
    let tokenPromises: Promise<ContractTransaction>[] = [];
    for (const token of tokens) {
      tokenPromises.push(token.connect(approver).approve(user.address, amount));
    }
    await Promise.all(tokenPromises);
}

export {
  setup,
  deployTrollToken,
  deployCommissionToken,
  signHash,
  getDigestNewProphecyClaim,
  getValidClaim,
  prefundAccount,
  preApproveAccount
};
