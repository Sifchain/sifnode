import createEthbridgeService from ".";
import { Asset, AssetAmount, Token } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { advanceBlock } from "../../test/utils/advanceBlock";
import { juniper, akasha, ethAccounts } from "../../test/utils/accounts";
import {
  createTestEthService,
  createTestSifService,
} from "../../test/utils/services";
import { getTestingTokens } from "../../test/utils/getTestingToken";
import config from "../../config.localnet.json";
import Web3 from "web3";
import JSBI from "jsbi";
import { sleep } from "../../test/utils/sleep";
import { waitFor } from "../../test/utils/waitFor";

const [ETH, CETH, ATK, CATK, ROWAN, EROWAN] = getTestingTokens([
  "ETH",
  "CETH",
  "ATK",
  "CATK",
  "ROWAN",
  "EROWAN",
]);

describe("EthbridgeService", () => {
  let EthbridgeService: ReturnType<typeof createEthbridgeService>;

  beforeEach(async () => {
    EthbridgeService = createEthbridgeService({
      sifApiUrl: "http://localhost:1317",
      sifWsUrl: "ws://localhost:26657/nosocket",
      sifRpcUrl: "http://localhost:26657",
      sifChainId: "sifchain",
      bridgebankContractAddress: config.bridgebankContractAddress,
      bridgetokenContractAddress: (EROWAN as Token).address,
      getWeb3Provider,
    });
  });

  // We need to only run one test on ebrelayer as we have not got the
  // infrastructure setup to retart it between tests
  // To fix this we would need to deterministically reset the state of both
  // blockchains as well as restart ebrelayer
  test("eth -> ceth -> eth then rowan -> erowan -> rowan ", async () => {
    // Setup services
    const sifService = await createTestSifService(juniper);
    const ethService = await createTestEthService();

    function getEthAddress() {
      return ethAccounts[2].public;
    }

    function getSifAddress() {
      return juniper.address;
    }

    async function getEthBalance() {
      const [bal] = await ethService.getBalance(getEthAddress(), ETH);
      return bal.toBaseUnits();
    }

    async function getCethBalance() {
      const [bal] = await sifService.getBalance(getSifAddress(), CETH);
      return bal.toBaseUnits();
    }

    const web3 = new Web3(await getWeb3Provider());

    ////////////////////////
    // ETH -> cETH
    ////////////////////////

    // Check the balance
    const cethBalance = await getCethBalance();

    const amountToLock = AssetAmount(ETH, "3");

    // Send funds to the smart contract
    await new Promise<void>(async done => {
      EthbridgeService.lockToSifchain(getSifAddress(), amountToLock, 100)
        .onTxHash(() => {
          advanceBlock(100);
        })
        .onComplete(async () => {
          done();
        })
        .onError(err => {
          throw err.payload;
        });
    });

    const expectedCethAmount = JSBI.add(
      cethBalance,
      amountToLock.toBaseUnits()
    );

    await waitFor(
      async () => await getCethBalance(),
      expectedCethAmount,
      "expectedCethAmount"
    );

    ////////////////////////
    // cETH -> ETH
    ////////////////////////

    const recipientBalanceBefore = await getEthBalance();

    const amountToSend = AssetAmount(CETH, "2");
    const feeAmount = AssetAmount(
      Asset.get("ceth"),
      JSBI.BigInt("58560000000000000")
    );

    const message = await EthbridgeService.burnToEthereum({
      fromAddress: getSifAddress(),
      assetAmount: amountToSend,
      feeAmount,
      ethereumRecipient: getEthAddress(),
    });

    // Message has the expected format
    const ethereumChainId = await web3.eth.net.getId();
    expect(message.value.msg).toEqual([
      {
        type: "ethbridge/MsgBurn",
        value: {
          amount: "2000000000000000000",
          ceth_amount: "58560000000000000",
          cosmos_sender: getSifAddress(),
          symbol: "ceth",
          ethereum_chain_id: `${ethereumChainId}`,
          ethereum_receiver: getEthAddress(),
        },
      },
    ]);

    await sifService.signAndBroadcast(message.value.msg);

    const expectedEthAmount = JSBI.add(
      recipientBalanceBefore,
      JSBI.BigInt("2000000000000000000")
    ).toString();

    await waitFor(
      async () => await getEthBalance(),
      expectedEthAmount,
      "expectedEthAmount"
    );

    ////////////////////////
    // Rowan -> eRowan
    ////////////////////////

    async function getERowanBalance() {
      const bals = await ethService.getBalance(getEthAddress(), EROWAN);
      return bals[0].toBaseUnits();
    }

    async function getRowanBalance() {
      const bals = await sifService.getBalance(getSifAddress(), ROWAN);
      return bals[0].toBaseUnits();
    }

    // First get balance in ethereum
    const startingERowanBalance = await getERowanBalance();

    // lock Rowan to eRowan
    const sendRowanAmount = AssetAmount(ROWAN, "100");

    const msg = await EthbridgeService.lockToEthereum({
      fromAddress: getSifAddress(),
      assetAmount: sendRowanAmount,
      ethereumRecipient: getEthAddress(),
      feeAmount: AssetAmount(
        Asset.get("ceth"),
        JSBI.BigInt("54080000000000000")
      ),
    });

    expect(msg.value.msg).toEqual([
      {
        type: "ethbridge/MsgLock",
        value: {
          amount: "100000000000000000000",
          ceth_amount: "54080000000000000",
          cosmos_sender: getSifAddress(),
          ethereum_chain_id: `${ethereumChainId}`,
          ethereum_receiver: getEthAddress(),
          symbol: "rowan",
        },
      },
    ]);

    await sifService.signAndBroadcast(msg.value.msg);

    await sleep(2000);

    const expectedERowanBalance = JSBI.add(
      startingERowanBalance,
      sendRowanAmount.toBaseUnits()
    );

    await waitFor(
      async () => await getERowanBalance(),
      expectedERowanBalance,
      "expectedERowanBalance"
    );

    ////////////////////////
    // eRowan -> Rowan
    ////////////////////////
    const startingRowanBalance = await getRowanBalance();

    const sendERowanAmount = AssetAmount(EROWAN, "10");

    await EthbridgeService.approveBridgeBankSpend(
      getEthAddress(),
      sendERowanAmount
    );

    // Burn eRowan to Rowan
    await new Promise<void>((done, reject) => {
      EthbridgeService.burnToSifchain(
        getSifAddress(),
        sendERowanAmount,
        50,
        getEthAddress()
      )
        .onTxHash(() => {
          advanceBlock(52);
        })
        .onComplete(async () => {
          done();
        })
        .onError(err => {
          reject(err);
        });
    });

    await sleep(2000);

    // wait for the balance to change
    const expectedRowanBalance = JSBI.add(
      startingRowanBalance,
      sendERowanAmount.toBaseUnits()
    );

    await waitFor(
      async () => await getRowanBalance(),
      expectedRowanBalance,
      "expectedRowanBalance"
    );
  });
});
